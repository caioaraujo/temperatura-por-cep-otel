package main

import (
	"net/http"

	"github.com/caioaraujo/temperatura-por-cep-otel/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	// ctx := context.Background()
	// traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	// if err != nil {
	// 	panic(err)
	// }

	cepHandler := handlers.NewCepHandler()
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/cep", func(r chi.Router) {
		r.Post("/", cepHandler.PostCep)
	})

	http.ListenAndServe(":8080", r)
}