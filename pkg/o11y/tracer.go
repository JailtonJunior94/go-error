package o11y

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	StartSpan(ctx context.Context, name string, attrs ...any) (context.Context, EndSpanFunc)
	WithSpan(ctx context.Context, name string, fn func(ctx context.Context) error) error
}

type EndSpanFunc func(err error)

type tracer struct {
	tracer trace.Tracer
}

func NewTracer(ctx context.Context, endpoint, serviceName string, resource *resource.Resource) (Tracer, func(context.Context) error, error) {
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize trace exporter grpc: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSyncer(traceExporter),
		sdktrace.WithResource(resource),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	shutdown := func(ctx context.Context) error {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	}

	return &tracer{tracer: tracerProvider.Tracer(serviceName)}, shutdown, nil
}

func (o *tracer) StartSpan(ctx context.Context, name string, attrs ...any) (context.Context, EndSpanFunc) {
	var kvs []attribute.KeyValue
	for i := 0; i+1 < len(attrs); i += 2 {
		k, ok1 := attrs[i].(string)
		v := attrs[i+1]
		if ok1 {
			switch tv := v.(type) {
			case string:
				kvs = append(kvs, attribute.String(k, tv))
			case int64:
				kvs = append(kvs, attribute.Int64(k, tv))
			case int:
				kvs = append(kvs, attribute.Int(k, tv))
			case bool:
				kvs = append(kvs, attribute.Bool(k, tv))
			default:
				kvs = append(kvs, attribute.String(k, fmt.Sprintf("%v", tv)))
			}
		}
	}

	ctx2, span := o.tracer.Start(ctx, name)
	if len(kvs) > 0 {
		span.SetAttributes(kvs...)
	}

	ended := false
	endFn := func(err error) {
		if ended {
			return
		}

		ended = true
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			return
		}

		span.SetStatus(codes.Ok, "OK")
		span.End()
	}
	return ctx2, endFn
}

func (o *tracer) WithSpan(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	ctx2, end := o.StartSpan(ctx, name)
	start := time.Now()
	err := fn(ctx2)
	end(err)
	_ = start
	return err
}
