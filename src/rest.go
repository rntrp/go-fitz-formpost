package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const address = ":8080"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", Welcome)
	r.Post("/scale", Scale)
	http.ListenAndServe(address, r)
}
