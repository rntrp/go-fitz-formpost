package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rntrp/go-fitz-formpost/internal/config"
	"github.com/rntrp/go-fitz-formpost/internal/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	config.Load()
}

func main() {
	if err := start(); err != nil {
		log.Fatalln(err)
	}
	log.Println("Bye.")
}

func start() error {
	srvout := make(chan error, 1)
	signal := make(chan os.Signal, 1)
	srv := server(signal)
	go shutdownMonitor(signal, srvout, srv)
	log.Println("Starting server at " + srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return <-srvout
}

func server(sig chan os.Signal) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", rest.Index)
	r.Get("/index.html", rest.Index)
	r.Get("/live", rest.Live)
	r.Post("/pages", rest.NumPage)
	r.Post("/scale", rest.Convert)
	if config.IsEnablePrometheus() {
		r.Handle("/metrics", promhttp.Handler())
	}
	if config.IsEnableShutdown() {
		r.Post("/shutdown", shutdownFn(sig))
	}
	return &http.Server{Addr: config.GetTCPAddress(), Handler: r}
}

func shutdownFn(sig chan os.Signal) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		go func() { sig <- os.Interrupt }()
	}
}

func shutdownMonitor(sig chan os.Signal, out chan error, srv *http.Server) {
	timeout := config.GetShutdownTimeout()
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sigName := (<-sig).String()
	log.Println("Signal received: " + sigName)
	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	out <- srv.Shutdown(ctx)
}
