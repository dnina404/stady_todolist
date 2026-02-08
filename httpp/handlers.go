package httpp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"restapi/todo"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHeandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHeandlers {
	return &HTTPHeandlers{
		todoList: todoList,
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	dto := ErrorDTO{
		Message: err.Error(),
		Time:    time.Now(),
	}
	http.Error(w, dto.ToString(), status)
}

/*
pattern: /tasks
method: POST
info:   JSON in HTTP request body

succeed:
  - status code: 201 Created
  - responce body: JSON represent created task

failed:

  - status code: 400 Bad Request, 409 Conflict Already Exists, 500 error on service

  - responce body:  JSON with error + time

    Контракт удовлетворяет требования парадигмы РестАпи
*/
func (h *HTTPHeandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var TaskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&TaskDTO); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := TaskDTO.ValidateForCreate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	todoTask := todo.NewTask(TaskDTO.Title, TaskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {
		if errors.Is(err, todo.ErrTaskAlreadyExists) {
			writeError(w, http.StatusConflict, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}
		return
	}

	b, err := json.MarshalIndent(todoTask, "", "    ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: GET
info: pattern

succeed:
  - status code: 200 Ok
  - responce body: JSON represent found task

failed:
  - status code: 400, 404 Not Found, 500
  - responce body: JSON with error + time
*/
func (h *HTTPHeandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title, _ := mux.Vars(r)["title"]
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("BadRequest")
		return
	}
	task, err := h.todoList.GetTask(title)
	if err != nil {
		if errors.Is(err, todo.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, err)
		} else {
			writeError(w, http.StatusInternalServerError, err)
		}

		return
	}
	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks
method: GET
info:	-

succeed:
  - status code: 200 Ok
  - responce body: JSON represent found tasks

failed:
  - status code: 400, 500
  - responce body: JSON with error + time
*/
func (h *HTTPHeandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListTasks()
	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks?completed=true
method: GET
info:   query params

succeed:
  - status code: 200 Ok
  - responce body: JSON represent found tasks

failed:
  - status code: 400, 500
  - responde body: JSON with error + time
*/
func (h *HTTPHeandlers) HandleGetUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListNotCompletedTasks()
	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: PATCH
info:   pattern + JSON in request body

succeed:
  - status code: 200 Ok
  - responce body: JSON represent changed task

failed:
  - status code: 400, 404, 409, 500
  - responde body: JSON with error + time
*/
func (h *HTTPHeandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	title, _ := mux.Vars(r)["title"]
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("BadRequest")
		return
	}
	if completeDTO.Complete {
		if err := h.todoList.CompleteTask(title); err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
	}

	if !completeDTO.Complete {
		if err := h.todoList.UnCompleteTask(title); err != nil {
			writeError(w, http.StatusNotFound, err)
			return
		}
	}
	task, err := h.todoList.GetTask(title)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response", err)
		return
	}
}

/*
pattern: /tasks/{title}
method: DELETE
info: pattern

succeed:
  - status code: 204 No Content
  - responce body: -

failed:
  - status code: 400, 404, 500
  - responce body: JSON with error + time
*/
func (h *HTTPHeandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title, _ := mux.Vars(r)["title"]
	if title == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("BadRequest")
		return
	}

	if err := h.todoList.DeleteTask(title); err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
