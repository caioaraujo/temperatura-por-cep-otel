package main

import (
	"net/http"

	"github.com/caioaraujo/temperatura-por-cep-otel/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	weatherHandler := handlers.NewWeatherHandler()
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/temperatura", func(r chi.Router) {
		r.Get("/{cep}", weatherHandler.GetWeather)
	})

	http.ListenAndServe(":8081", r)
}
