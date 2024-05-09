package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/caioaraujo/temperatura-por-cep-otel/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var logger = log.New(os.Stderr, "zipkin-example", log.Ldate|log.Ltime|log.Llongfile)

func main() {
	url := flag.String("zipkin", "http://localhost:9411/api/v2/spans", "zipkin url")
	flag.Parse()

	consoleMetricExporter, err := newMetricExporter()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	exporter, err := zipkin.New(
		*url,
		zipkin.WithLogger(logger),
	)
	if err != nil {
		panic(err)
	}

	batcher := trace.NewBatchSpanProcessor(exporter)

	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(batcher),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("zipkin-test"),
		)),
	)
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(ctx)

	meterProvider := newMeterProvider(consoleMetricExporter)
	defer meterProvider.Shutdown(ctx)
	otel.SetMeterProvider(meterProvider)

	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	cepHandler := handlers.NewCepHandler()
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/cep", func(r chi.Router) {
		r.Post("/", cepHandler.PostCep)
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

func newTraceExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func newMetricExporter() (metric.Exporter, error) {
	return stdoutmetric.New()
}

func newMeterProvider(meterExporter metric.Exporter) *metric.MeterProvider {
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(meterExporter)),
	)
	return meterProvider
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	)
}
