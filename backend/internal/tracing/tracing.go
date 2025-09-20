package tracing

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	// DefaultTracerProvider is the default tracer provider
	DefaultTracerProvider *sdktrace.TracerProvider
)

// InitializeDefaultTracerProvider initializes the default tracer provider with the given configuration
func InitializeDefaultTracerProvider(config Config) error {
	if !config.Enabled {
		log.Println("Tracing is disabled")
		return nil
	}

	// Create a new resource with service information
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.DeploymentEnvironmentKey.String(config.Environment),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Create a new trace exporter based on the configuration
	var exporter sdktrace.SpanExporter

	if config.JaegerURL != "" {
		// Create a Jaeger exporter
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerURL)))
		if err != nil {
			return fmt.Errorf("failed to create Jaeger exporter: %w", err)
		}
	} else if config.ZipkinURL != "" {
		// Create a Zipkin exporter
		exporter, err = zipkin.New(config.ZipkinURL)
		if err != nil {
			return fmt.Errorf("failed to create Zipkin exporter: %w", err)
		}
	} else {
		// If no exporter URL is provided, use a console exporter for development
		exporter = &ConsoleExporter{}
	}

	// Create a new trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(config.SamplingRate)),
	)

	// Set the global tracer provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Set the default tracer provider
	DefaultTracerProvider = tp

	log.Println("Tracing initialized successfully")
	return nil
}

// ShutdownTracerProvider shuts down the default tracer provider
func ShutdownTracerProvider(ctx context.Context) error {
	if DefaultTracerProvider == nil {
		return nil
	}

	return DefaultTracerProvider.Shutdown(ctx)
}

// Tracer returns a tracer with the given name
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer(name).Start(ctx, name)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attributes ...trace.Attribute) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.AddEvent(name, trace.WithAttributes(attributes...))
	}
}

// SetSpanError sets an error on the current span
func SetSpanError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span != nil && err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}
}

// GetTraceID returns the trace ID from the context
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sc := span.SpanContext()
	if !sc.IsValid() {
		return ""
	}

	return sc.TraceID().String()
}

// GetSpanID returns the span ID from the context
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sc := span.SpanContext()
	if !sc.IsValid() {
		return ""
	}

	return sc.SpanID().String()
}

// ConsoleExporter is a simple console exporter for development
type ConsoleExporter struct{}

// ExportSpans exports spans to the console
func (e *ConsoleExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, span := range spans {
		fmt.Printf("TraceID: %s, SpanID: %s, ParentSpanID: %s, Name: %s, StartTime: %s, EndTime: %s\n",
			span.SpanContext().TraceID(),
			span.SpanContext().SpanID(),
			span.Parent().SpanID(),
			span.Name(),
			span.StartTime(),
			span.EndTime(),
		)
	}
	return nil
}

// Shutdown shuts down the console exporter
func (e *ConsoleExporter) Shutdown(ctx context.Context) error {
	return nil
}