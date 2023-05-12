package app

import (
	"context"
	"errors"
	"github.com/Parsa-Sedigh/rebottle/internal/appjwt"
	"github.com/Parsa-Sedigh/rebottle/pkg/jsonutil"
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
				jsonutil.ErrorJSON(w, app.logger, errors.New("no authorization header received"), http.StatusUnauthorized)
				return
			}

			JWTData, err := appjwt.ExtractClaims(headerParts[1])
			if err != nil {
				jsonutil.ErrorJSON(w, app.logger, errors.New("you're unauthorized"), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "JWTData", JWTData)

			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		jsonutil.ErrorJSON(w, app.logger, errors.New("you're unauthorized"), http.StatusUnauthorized)
	})
}
