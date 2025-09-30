package o11y

type Telemetry interface {
	Tracer() Tracer
	Metrics() Metrics
	Logger() Logger
}
