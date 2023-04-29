package main

import (
	"context"
	"errors"
	"github.com/Parsa-Sedigh/rebottle/internal/appjwt"
	"net/http"
	"strings"
)

func (app *application) SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func (app *application) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			headerParts := strings.Split(r.Header.Get("Authorization"), " ")

			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				app.errorJSON(w, errors.New("no authorization header received"), http.StatusUnauthorized)
				return
			}

			JWTData, err := appjwt.ExtractClaims(headerParts[1])
			if err != nil {
				app.errorJSON(w, errors.New("you're unauthorized"), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "JWTData", JWTData)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		app.errorJSON(w, errors.New("you're unauthorized"), http.StatusUnauthorized)
	})
}
