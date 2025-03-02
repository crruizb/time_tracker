package api

import (
	"net/http"

	"golang.org/x/oauth2"
)

type Router struct {
	*http.ServeMux
	oauthConfigs map[string]*oauth2.Config
	ps           ProjectsStore
	ts           TasksStore
}

func NewRouter(oauthConfigs map[string]*oauth2.Config, ps ProjectsStore, ts TasksStore) *Router {
	rt := &Router{
		ServeMux:     http.NewServeMux(),
		oauthConfigs: oauthConfigs,
		ps:           ps,
		ts:           ts,
	}

	rt.HandleFunc("POST /api/project", rt.createProject)
	rt.HandleFunc("POST /api/project/{projectId}/task", rt.createTask)
	rt.HandleFunc("POST /api/task/{taskId}/start", rt.startTask)
	rt.HandleFunc("POST /api/task/{taskId}/stop", rt.stopTask)
	rt.HandleFunc("GET /api/project/{projectId}", rt.projectReport)

	rt.HandleFunc("GET /auth/login/{source}", rt.oauthLogin)
	rt.HandleFunc("GET /auth/callback", rt.oauthCallback)

	return rt
}
