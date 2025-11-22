package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Otel struct {
	OTELTracer      trace.Tracer
	RequestNameOtel string
}

func New(trace trace.Tracer, serviceName string) *Otel {
	return &Otel{
		OTELTracer:      trace,
		RequestNameOtel: serviceName,
	}
}

func InitOtel(serverName, collectorUrl string) error {

	ctx := context.Background()

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serverName)))

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, collectorUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(tracerProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	defer func() {
		_ = tracerProvider.Shutdown(context.Background())
		//if err := shutdown(ctx); err != nil {
		//log.Fatal("failed to shutdown TracerProvider: %w", err)}
	}()

	return nil
}
