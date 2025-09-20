package tracing

// Config represents the configuration for the OpenTelemetry tracer
type Config struct {
	// Enabled indicates whether tracing is enabled
	Enabled bool
	
	// ServiceName is the name of the service
	ServiceName string
	
	// Environment is the environment the service is running in
	Environment string
	
	// JaegerURL is the URL of the Jaeger collector
	JaegerURL string
	
	// ZipkinURL is the URL of the Zipkin collector
	ZipkinURL string
	
	// SamplingRate is the sampling rate for traces (0.0 to 1.0)
	SamplingRate float64
}