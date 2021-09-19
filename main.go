package main

import (
	"net/http"

	"github.com/rntrp/go-fitz-rest-example/internal/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const address = ":8080"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", rest.Welcome)
	r.Post("/scale", rest.Scale)
	http.ListenAndServe(address, r)
}
