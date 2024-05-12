package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/caioaraujo/temperatura-por-cep-otel/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serviceName = "temperatura-cep"
	serviceUrl  = "collector:4317"
)

func main() {
	log.Println("Starting service-b server.")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutdown, err := initProvider(ctx, serviceName, serviceUrl)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	weatherHandler := handlers.NewWeatherHandler()
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Route("/temperatura", func(r chi.Router) {
		r.Get("/{cep}", weatherHandler.GetWeather)
	})

	http.ListenAndServe(":8081", r)
}

func initProvider(ctx context.Context, serviceName, serviceURL string) (shutdown func(context.Context) error, err error) {
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(serviceName)))
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		serviceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	tracerExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	bsp := trace.NewBatchSpanProcessor(tracerExporter)
	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(traceProvider)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceProvider.Shutdown, nil
}
