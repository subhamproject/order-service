package otelsvc

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	service     = "order-service"
	environment = "development"
	id          = 1
)

// StdoutTraceExporter returns a console exporter.
func StdoutTraceExporter() (*stdouttrace.Exporter, error) {
	return stdouttrace.New(
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		//stdouttrace.WithoutTimestamps(),
	)
}

// OtelTraceExporter returns a Otel exporter.
func OtelTraceExporter(ctx context.Context, url string) (*otlptrace.Exporter, error) {
	log.Printf("connecting to otel %s", url)
	conn, err := grpc.DialContext(ctx, url, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	reportErr(err, "failed to create gRPC connection to collector")

	// Set up a trace exporter
	traceExporter, err := newExporter(ctx, conn)
	reportErr(err, "failed to create trace exporter")
	return traceExporter, nil
}

func newExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

// // newResource returns a resource describing this application.
// func newResource() *resource.Resource {
// 	r, _ := resource.Merge(
// 		resource.Default(),
// 		resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			semconv.ServiceName("fib"),
// 			semconv.ServiceVersion("v0.1.0"),
// 			attribute.String("environment", "demo"),
// 		),
// 	)
// 	return r
// }

// InitTracerProvider returns an OpenTelemetry InitTracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// InitTracerProvider will also use a Resource configured with all the information
// about the application.
func InitTracerProvider(url string, otelEnable bool) func() {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	// // Create the Jaeger exporter
	// exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	// if err != nil {
	// 	return nil, err
	// }
	var tracerProvider *tracesdk.TracerProvider
	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(service),
		attribute.String("environment", environment),
		attribute.Int64("ID", id),
	)
	stdoutExp, err := StdoutTraceExporter()
	reportErr(err, "failed to create stdout trace exporter")

	if otelEnable {
		// Set up a trace exporter
		traceExporter, err := OtelTraceExporter(ctx, url)
		reportErr(err, "failed to create otel trace exporter")
		tracerProvider = tracesdk.NewTracerProvider(
			// Always be sure to batch in production.
			tracesdk.WithBatcher(traceExporter),
			tracesdk.WithBatcher(stdoutExp),
			// Record information about this application in a Resource.
			tracesdk.WithResource(resources),
		)
	} else {
		tracerProvider = tracesdk.NewTracerProvider(
			// Always be sure to batch in production.
			tracesdk.WithBatcher(stdoutExp),
			// Record information about this application in a Resource.
			tracesdk.WithResource(resources),
		)
	}
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return func() {
		// Shutdown will flush any remaining spans and shut down the exporter.
		reportErr(tracerProvider.Shutdown(ctx), "failed to shutdown TracerProvider")
		cancel()
	}
	//return tp, nil
}

func reportErr(err error, message string) {
	if err != nil {
		log.Printf("%s: %v", message, err)
	}
}
