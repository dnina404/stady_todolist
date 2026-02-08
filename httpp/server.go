package httpp

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	HTTPHeandlers *HTTPHeandlers
}

func NewHTTPServer(httpHandler *HTTPHeandlers) *HTTPServer {
	return &HTTPServer{
		HTTPHeandlers: httpHandler,
	}
}

func (s *HTTPServer) StartServer() error {
	router := mux.NewRouter()
	router.Path("/tasks").Methods("POST").HandlerFunc(s.HTTPHeandlers.HandleCreateTask)
	router.Path("/tasks/{title}").Methods("GET").HandlerFunc(s.HTTPHeandlers.HandleGetTask)
	router.Path("/tasks").Methods("GET").HandlerFunc(s.HTTPHeandlers.HandleGetAllTasks)
	router.Path("/tasks").Methods("GET").Queries("completed", "true").HandlerFunc(s.HTTPHeandlers.HandleGetUncompletedTasks)
	router.Path("/tasks/{title}").Methods("PATCH").HandlerFunc(s.HTTPHeandlers.HandleCompleteTask)
	router.Path("/tasks/{title}").Methods("DELETE").HandlerFunc(s.HTTPHeandlers.HandleDeleteTask)

	return http.ListenAndServe(":9091", router)
}
