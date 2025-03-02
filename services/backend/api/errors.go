package api

import (
	"log/slog"
	"net/http"
)

func logError(r *http.Request, err error) {
	slog.Error(err.Error(), "request_method", r.Method, "request_url", r.URL.String())
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)
	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func ForbiddenResponse(w http.ResponseWriter, r *http.Request, err string) {
	errorResponse(w, r, http.StatusForbidden, err)
}
