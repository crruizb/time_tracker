package api

import (
	"errors"
	"net/http"

	"github.com/crruizb/data"
)

type ProjectsStore interface {
	CreateProject(name, description, userId string) (*data.Project, error)
}

type TasksStore interface {
	CreateTask(projectId, name, description string) (*data.Task, error)
	StartTask(taskId, userId string) error
	StopTask(taskId, userId string) error
}

type contextUserKey string

const ContextUser = contextUserKey("ctxUser")

func (s *Router) createProject(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ContextUser).(*data.User)
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := ReadJSON(w, r, &input)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	project, err := s.ps.CreateProject(input.Name, input.Description, user.Id)
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}

	WriteJSON(w, http.StatusAccepted, project, nil)
}

func (s *Router) createTask(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("projectId")
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := ReadJSON(w, r, &input)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	task, err := s.ts.CreateTask(projectId, input.Name, input.Description)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrProjectNotFound):
			NotFoundResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}

		return
	}

	WriteJSON(w, http.StatusAccepted, task, nil)
}

func (s *Router) startTask(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ContextUser).(*data.User)
	taskId := r.PathValue("taskId")
	err := s.ts.StartTask(taskId, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrTaskAlreadyStarted):
			BadRequestResponse(w, r, err)
		default:
			ServerErrorResponse(w, r, err)
		}
		return

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Router) stopTask(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(ContextUser).(*data.User)
	taskId := r.PathValue("taskId")
	err := s.ts.StopTask(taskId, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrTaskNotStarted):
			BadRequestResponse(w, r, err)
		case errors.Is(err, data.ErrTaskAlreadyFinished):
			BadRequestResponse(w, r, err)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (s *Router) projectReport(w http.ResponseWriter, r *http.Request) {

}
