package main

import (
	"context"
	"net/http"
	"time"

	"github.com/caioaraujo/temperatura-por-cep-otel/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	consoleTraceExporter, err := newTraceExporter()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	tracerProvider := newTraceProvider(consoleTraceExporter)
	defer tracerProvider.Shutdown(ctx)
	otel.SetTracerProvider(tracerProvider)

	cepHandler := handlers.NewCepHandler()
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/cep", func(r chi.Router) {
		r.Post("/", cepHandler.PostCep)
	})

	http.ListenAndServe(":8080", r)
}

func newTraceExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func newTraceProvider(traceExporter trace.SpanExporter) *trace.TracerProvider {
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second)),
	)
	return traceProvider
}
