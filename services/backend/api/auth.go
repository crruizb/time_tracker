package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/crruizb/data"
	"golang.org/x/oauth2"
)

const GithubOauthURLApi = "https://api.github.com/user"

func (rt *Router) oauthLogin(w http.ResponseWriter, r *http.Request) {
	source := r.PathValue("source")
	oauthState := generateStateOauthCookie(w, source)

	u := rt.oauthConfigs[source].AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (rt *Router) oauthCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		slog.Error("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	source := strings.Split(oauthState.Value, ":")[1]

	tokens, err := rt.getTokens(r.FormValue("code"), source)
	if err != nil {
		slog.Error(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	user, _ := GetUserData(tokens)
	// WriteJSON(w, http.StatusOK, tokens, nil)
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		HttpOnly: true, // Prevent JavaScript access
		// SameSite: http.SameSiteStrictMode,
		Secure: false, // Set to true in production with HTTPS
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    user.Username,
		Path:     "/",
		HttpOnly: false, // Prevent JavaScript access
		// SameSite: http.SameSiteStrictMode,
		Secure: false, // Set to true in production with HTTPS
	})

	// WriteJSON(w, http.StatusOK, tokens, nil)
	// Redirect the user to the frontend dashboard after successful login
	http.Redirect(w, r, "http://localhost:5173/", http.StatusFound)
	// http.Redirect(w, r, "/", http.StatusFound)
}

func generateStateOauthCookie(w http.ResponseWriter, source string) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := fmt.Sprintf("%s:%s", base64.URLEncoding.EncodeToString(b), source)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration, HttpOnly: true, Path: "/auth"}
	http.SetCookie(w, &cookie)

	return state
}

func (rt *Router) getTokens(code, source string) (*oauth2.Token, error) {
	token, err := rt.oauthConfigs[source].Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	return token, nil
}

// Validate validates the access token.
func GetUserData(token *oauth2.Token) (*data.User, error) {
	r, _ := http.NewRequest("GET", GithubOauthURLApi, nil)
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
