package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type application struct {
	config config
	router *chi.Mux
}

type config struct {
	addr string
}

func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthcheckHandler)
	})

	app.router = router
	return router
}

func (app *application) run() error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      app.router,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Server has started on %s", app.config.addr)
	return srv.ListenAndServe()
}
