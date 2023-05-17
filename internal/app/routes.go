package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"path/filepath"
	"runtime"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	//mux.Use(zapLogger(app.logger))
	mux.Use(app.SessionLoad)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Route("/v1", func(r chi.Router) {
		r.Post("/signup", app.SignupUser)
		r.Post("/signup-driver", app.SignupDriver)

		r.Post("/verify-signup", app.VerifyUserSignup)
		r.Post("/verify-signup-driver", app.VerifyDriverSignup)
		r.Post("/login", app.Login)
		r.Post("/token", app.NewAuthTokens)

		//////////////// private routes ////////////////
		r.Group(func(r chi.Router) {
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
				r.Put("/{id}", app.UpdateDriver)
			})
		})
	})
	////////////////////////////////

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Join(b + "../../.."))
	fileserver := http.FileServer(http.Dir(basepath + "/static"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileserver))

	return mux
}

//func zapLogger(l *zap.Logger) func(next http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
//			t1 := time.Now()
//			defer func() {
//				l.Info("Served",
//					zap.String("proto", r.Proto),
//					zap.String("path", r.URL.Path),
//					zap.Duration("lat", time.Since(t1)),
//					zap.Int("status", ww.Status()),
//					zap.Int("size", ww.BytesWritten()),
//					zap.String("reqId", middleware.GetReqID(r.Context())))
//			}()
//
//			next.ServeHTTP(ww, r)
//		})
//	}
//}
