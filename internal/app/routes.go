package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(app.SessionLoad)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Post("/signup", app.SignupUser)
	mux.Post("/signup-driver", app.SignupDriver)

	mux.Post("/verify-signup", app.VerifyUserSignup)
	mux.Post("/verify-signup-driver", app.VerifyDriverSignup)
	mux.Post("/login", app.Login)
	mux.Post("/token", app.NewAuthTokens)

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

		r.Route("/driver", func(r chi.Router) {
			r.Put("/", app.UpdateDriver)
		})
	})
	////////////////////////////////

	fileserver := http.FileServer(http.Dir("../../static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))

	return mux
}
