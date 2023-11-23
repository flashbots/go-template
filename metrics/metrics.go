package metrics

import (
	"context"
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

var (
	metrics *Metrics
	mx      sync.Mutex
)

type Metrics struct {
	exporter *promotel.Exporter
	meter    metric.Meter
	provider *sdk.MeterProvider
	registry *prometheus.Registry

	server *http.Server

	float64Histogram   map[string]metric.Float64Histogram
	mxFloat64Histogram sync.RWMutex
}

func Init(name, address string) {
	mx.Lock()
	defer mx.Unlock()

	if metrics != nil {
		panic("double initialisation of metrics")
	}

	registry := prometheus.NewRegistry()

	exporter, err := promotel.New(promotel.WithRegisterer(registry))
	if err != nil {
		panic(err) // we don't do anything fancy above => there should be no errors
	}

	provider := sdk.NewMeterProvider(sdk.WithReader(exporter))
	meter := provider.Meter(name)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	server := &http.Server{
		Addr:    address,
		Handler: mux,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	metrics = &Metrics{
		exporter: exporter,
		meter:    meter,
		provider: provider,
		registry: registry,

		server: server,

		float64Histogram: make(map[string]metric.Float64Histogram),
	}
}

func ListenAndServe() error {
	if err := runtime.Start(
		runtime.WithMeterProvider(metrics.provider),
		runtime.WithMinimumReadMemStatsInterval(time.Minute), // querying go-runtime is expensive
	); err != nil {
		return err
	}
	return metrics.server.ListenAndServe()
}

func Shutdown(ctx context.Context) error {
	return metrics.server.Shutdown(ctx)
}

//nolint:ireturn,nolintlint
func Float64Histogram(
	name string,
	description string,
	uom string,
	bucketBounds ...float64,
) metric.Float64Histogram {
	metrics.mxFloat64Histogram.RLock()
	if h, exists := metrics.float64Histogram[name]; exists {
		metrics.mxFloat64Histogram.RUnlock()
		return h
	}
	metrics.mxFloat64Histogram.RUnlock()

	metrics.mxFloat64Histogram.Lock()
	defer metrics.mxFloat64Histogram.Unlock()

	// avoid race condition between ro-unlock and rw-lock
	if h, exists := metrics.float64Histogram[name]; exists {
		return h
	}

	h, err := metrics.meter.Float64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithExplicitBucketBoundaries(bucketBounds...),
		metric.WithUnit(uom),
	)
	if err != nil {
		panic(err)
	}

	metrics.float64Histogram[name] = h
	return h
}
