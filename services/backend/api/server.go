package api

import (
	"net/http"

	"github.com/crruizb/data"
)

type Server struct {
	addr string
	ps ProjectsStore
}

type ProjectsStore interface {
	CreateProject(name, description string) (*data.Project, error)
}

func NewServer(addr string, ps ProjectsStore) *Server {
	return &Server{
		addr: addr,
		ps: ps,
	}
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr: s.addr,
		Handler: s.routes(),
	}

	return srv.ListenAndServe()
}


func (s *Server) createProject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Description string `json:"description"`
	}

	err := s.readJSON(w, r, &input)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	project, err := s.ps.CreateProject(input.Name, input.Description)
	if err != nil {
		s.serverErrorResponse(w, r, err)
		return
	}

	s.writeJSON(w, http.StatusAccepted, project, nil)
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) startTask(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) stopTask(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) projectReport(w http.ResponseWriter, r *http.Request) {
	
}