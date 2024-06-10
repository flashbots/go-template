package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	promotel "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
)

const (
	UomMicroseconds = "us"
)

var BucketsRequestDuration = []float64{
	0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0, 30.0, 60.0,
}

type MetricsServer struct {
	exporter *promotel.Exporter
	meter    metric.Meter
	provider *sdk.MeterProvider
	registry *prometheus.Registry

	server *http.Server

	float64Histogram   map[string]metric.Float64Histogram
	mxFloat64Histogram sync.RWMutex
}

func New(name, addr string) (metricsServer *MetricsServer, err error) {
	registry := prometheus.NewRegistry()

	exporter, err := promotel.New(promotel.WithRegisterer(registry))
	if err != nil {
		return nil, fmt.Errorf("failed to create prometheus exporter: %w", err)
	}

	provider := sdk.NewMeterProvider(sdk.WithReader(exporter))
	meter := provider.Meter(name)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:    addr,
		Handler: mux,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	metricsServer = &MetricsServer{
		exporter: exporter,
		meter:    meter,
		provider: provider,
		registry: registry,

		server: server,

		float64Histogram: make(map[string]metric.Float64Histogram),
	}

	if err := runtime.Start(
		runtime.WithMeterProvider(metricsServer.provider),
		runtime.WithMinimumReadMemStatsInterval(time.Minute), // querying go-runtime is expensive
	); err != nil {
		return nil, fmt.Errorf("failed to start go-runtime's metrics exporter: %w", err)
	}
	return metricsServer, nil
}

func (ms *MetricsServer) ListenAndServe() error {
	return ms.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections.
func (ms *MetricsServer) Shutdown(ctx context.Context) error {
	return ms.server.Shutdown(ctx)
}

// Float64Histogram returns a float64 histogram with given name and
// parameters.
//
// Once created it's cached and reused further on. All subsequent calls
// to this method that use the same name will retrieve already created
// histogram from the cache.
//
// It is thread-safe.
//
// See also: https://pkg.go.dev/go.opentelemetry.io/otel/metric@v1.21.0#Meter.Float64Histogram
//
//nolint:ireturn,nolintlint
func (ms *MetricsServer) Float64Histogram(
	name string,
	description string,
	uom string,
	bucketBounds ...float64,
) metric.Float64Histogram {
	ms.mxFloat64Histogram.RLock()
	if h, exists := ms.float64Histogram[name]; exists {
		ms.mxFloat64Histogram.RUnlock()
		return h
	}
	ms.mxFloat64Histogram.RUnlock()

	ms.mxFloat64Histogram.Lock()
	defer ms.mxFloat64Histogram.Unlock()

	// avoid race condition between ro-unlock and rw-lock
	if h, exists := ms.float64Histogram[name]; exists {
		return h
	}

	h, err := ms.meter.Float64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithExplicitBucketBoundaries(bucketBounds...),
		metric.WithUnit(uom),
	)
	if err != nil {
		panic(err)
	}

	ms.float64Histogram[name] = h
	return h
}
