package tracer

import (
	"context"
	"log/slog"

	"github.com/walnuts1018/mpeg-dash-encoder/consts"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var Tracer = otel.Tracer("github.com/walnuts1018/mpeg-dash-encoder")

func NewTracerProvider(ctx context.Context) (func(), error) {
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(consts.ApplicationName),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	closeFunc := func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown TracerProvider", slog.Any(
				"error", err,
			))
		}
	}
	return closeFunc, nil
}
