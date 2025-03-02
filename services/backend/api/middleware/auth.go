package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/crruizb/api"
	"github.com/crruizb/data"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

const (
	authorizationHeader = "authorization"
	authorizationPrefix = "bearer"
	githubOauthURLApi   = "https://api.github.com/user"
)

var (
	errInvalidToken = "invalid token"
	errExpiredToken = "expired token"
)

type UserStore interface {
	GetUser(username, source string) (*data.User, error)
	InsertUser(username, source string) (*data.User, error)
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

			fields := strings.Fields(r.Header.Get(authorizationHeader))
			if len(fields) < 2 {
				api.ForbiddenResponse(w, r, errInvalidToken)
				return
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
			oauthUser, err := GetUserData(token)
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

// Validate validates the access token.
func GetUserData(token *oauth2.Token) (*data.User, error) {
	r, _ := http.NewRequest("GET", githubOauthURLApi, nil)
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed getting user info")
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	user := &data.User{
		Source: "Github",
	}
	json.Unmarshal(contents, user)

	return user, nil
}
