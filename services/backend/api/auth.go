package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func (rt *Router) oauthLogin(w http.ResponseWriter, r *http.Request) {
	source := r.PathValue("source")
	println(source)
	oauthState := generateStateOauthCookie(w, source)

	u := rt.oauthConfigs[source].AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (rt *Router) oauthCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, err := r.Cookie("oauthstate")
	println(err)

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
	WriteJSON(w, http.StatusOK, tokens, nil)
}

func generateStateOauthCookie(w http.ResponseWriter, source string) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := fmt.Sprintf("%s:%s", base64.URLEncoding.EncodeToString(b), source)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
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
