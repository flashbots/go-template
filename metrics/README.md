# metrics

Introduction:
- https://opentelemetry.io/docs/languages/go/instrumentation/

Clients:
- https://github.com/open-telemetry/opentelemetry-go
- https://github.com/prometheus/client_golang
- https://github.com/VictoriaMetrics/metrics

Onboarding:
- https://www.timescale.com/blog/four-types-prometheus-metrics-to-collect/
- https://blog.pvincent.io/2017/12/prometheus-blog-series-part-1-metrics-and-labels/
- https://pierrevincent.github.io/2017/12/prometheus-blog-series-part-2-metric-types/
- https://pierrevincent.github.io/2017/12/prometheus-blog-series-part-3-exposing-and-collecting-metrics/
- https://pierrevincent.github.io/2017/12/prometheus-blog-series-part-4-instrumenting-code-in-go-and-java/
- https://developers.soundcloud.com/blog/prometheus-monitoring-at-soundcloud

Best practices:
- https://prometheus.io/docs/practices/naming/

See also:
- https://pkg.go.dev/github.com/prometheus/client_golang@v1.17.0/prometheus#NewRegistry
- https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus@v0.44.0#New
- https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric@v1.21.0#NewMeterProvider
- https://pkg.go.dev/go.opentelemetry.io/otel/sdk/metric@v1.21.0#MeterProvider.Meter
- https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/runtime@v0.46.1#Start