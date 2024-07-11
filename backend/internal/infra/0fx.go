package infra

import (
	"github.com/getsentry/sentry-go"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"platform/pkg/db/tx"
)

var Module = fx.Module("infra",
	fx.Provide(NewLogger),
	fx.Provide(tx.NewManager),
)

func init() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
		sentryotel.NewSentryPropagator(),
	))

	provider := trace.NewTracerProvider(
		trace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()),
	)
	otel.SetTracerProvider(provider)

	otel.SetMeterProvider(metric.NewMeterProvider())

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:           "https://7c1d438cc260efab4d91c20023a6656f@o4507585172537344.ingest.de.sentry.io/4507585191673936",
		EnableTracing: true,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	}); err != nil {
		panic(err)
	}

}
