package infra

import (
	"fmt"
	"os"

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

func Init() error {
	dsn, withSentry := os.LookupEnv("SENTRY_DSN")

	propagators := []propagation.TextMapPropagator{
		propagation.TraceContext{},
		propagation.Baggage{},
	}
	if withSentry {
		propagators = append(propagators, sentryotel.NewSentryPropagator())
	}

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagators...))

	var tracerOptions []trace.TracerProviderOption
	if withSentry {
		tracerOptions = append(tracerOptions, trace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()))
	}

	provider := trace.NewTracerProvider(tracerOptions...)
	otel.SetTracerProvider(provider)

	otel.SetMeterProvider(metric.NewMeterProvider())

	if withSentry {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:           dsn,
			EnableTracing: true,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate: 1.0,
			AttachStacktrace: true,
		}); err != nil {
			return fmt.Errorf("sentry.Init: %w", err)
		}
	}
	return nil
}
