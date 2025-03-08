package middleware

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/crruizb/api"
	"github.com/crruizb/data"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

const (
	authorizationHeader     = "authorization"
	cookieAuthorizationName = "access_token"
	authorizationPrefix     = "bearer"
)

var (
	errInvalidToken = "invalid token"
	errExpiredToken = "expired token"
)

type UserStore interface {
	GetUser(username, source string) (*data.User, error)
	InsertUser(username, source string) (*data.User, error)
}

func CorsMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Allow specific origins (localhost and 127.0.0.1)
			if origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Maneja las solicitudes OPTIONS (preflight)
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Auth is a middleware that checks for OAuth2 authentication.
func Auth(us UserStore, excludedPaths []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, path := range excludedPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			fields := []string{}
			c := r.CookiesNamed(cookieAuthorizationName)
			if len(c) == 0 {
				fields = strings.Fields(r.Header.Get(authorizationHeader))
				if len(fields) < 2 {
					api.ForbiddenResponse(w, r, errInvalidToken)
					return
				}
			} else {
				fields = append(fields, "bearer")
				fields = append(fields, c[0].Value)
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != authorizationPrefix {
				api.ForbiddenResponse(w, r, errInvalidToken)
				return
			}

			token := &oauth2.Token{
				AccessToken: fields[1],
			}

			// validate token
			oauthUser, err := api.GetUserData(token)
			if err != nil {
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					slog.Info("token is expired")
					api.ForbiddenResponse(w, r, errExpiredToken)
					return
				default:
					api.ForbiddenResponse(w, r, errExpiredToken)
					return
				}
			}

			user, err := us.GetUser(oauthUser.Username, oauthUser.Source)
			if err != nil {
				switch {
				case errors.Is(err, sql.ErrNoRows):
					user, err = us.InsertUser(oauthUser.Username, oauthUser.Source)
					if err != nil {
						api.ServerErrorResponse(w, r, err)
						return
					}
				default:
					api.ForbiddenResponse(w, r, errExpiredToken)
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), api.ContextUser, user)))
		})
	}
}
