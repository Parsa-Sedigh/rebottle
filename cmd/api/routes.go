package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(app.SessionLoad)

	mux.Route("/pickup", func(r chi.Router) {
		r.Use(app.VerifyJWT)

		r.Get("/", app.GetPickups)
		r.Post("/", app.CreatePickup)
		r.Put("/", app.UpdatePickup)
		r.Patch("/", app.CancelPickup)
		r.Get("/{id}", app.GetPickup)
	})

	mux.Route("/user", func(r chi.Router) {
		r.Use(app.VerifyJWT)

		r.Get("/", app.GetUser)
	})

	mux.Post("/signup", app.SignupUser)
	mux.Post("/verify-signup", app.VerifyUserSignup)
	mux.Post("/login", app.Login)

	return mux
}
