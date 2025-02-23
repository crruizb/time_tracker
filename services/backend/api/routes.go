package api

import "net/http"

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/project", s.createProject)
	mux.HandleFunc("POST /api/project/{projectId}/task", s.createTask)
	mux.HandleFunc("GET /api/task/{taskId}/start", s.startTask)
	mux.HandleFunc("GET /api/task/{taskId}/stop", s.stopTask)
	mux.HandleFunc("GET /api/project/{projectId}", s.projectReport)


	return mux
}