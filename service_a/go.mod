module github.com/caioaraujo/temperatura-por-cep-otel/service_a

go 1.22.2

require (
	github.com/caioaraujo/temperatura-por-cep-otel v0.0.0-20240506164915-6fc2c4e47c2f
	github.com/go-chi/chi v1.5.5
	go.opentelemetry.io/otel v1.26.0
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.26.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0
	go.opentelemetry.io/otel/exporters/zipkin v1.26.0
	go.opentelemetry.io/otel/metric v1.26.0
	go.opentelemetry.io/otel/sdk v1.26.0
	go.opentelemetry.io/otel/sdk/metric v1.26.0
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/openzipkin/zipkin-go v0.4.2 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
)
