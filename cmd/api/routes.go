package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(app.SessionLoad)

	mux.Post("/signup", app.SignupUser)
	mux.Post("/verify-signup", app.VerifyUserSignup)
	mux.Post("/login", app.Login)

	//////////////// private routes ////////////////
	mux.Group(func(r chi.Router) {
		r.Use(app.VerifyJWT)

		r.Route("/pickup", func(r chi.Router) {
			r.Get("/", app.GetPickups)
			r.Post("/", app.CreatePickup)
			r.Put("/", app.UpdatePickup)
			r.Patch("/", app.CancelPickup)
			r.Get("/{id}", app.GetPickup)
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/", app.GetUser)
		})
	})
	////////////////////////////////

	return mux
}
