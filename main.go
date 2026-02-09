package main

import (
	"fmt"
	"net/http"
	"restapi/httpp"
	"restapi/todo"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
}

func main() {
	todoList := todo.NewList()
	httpHandlers := httpp.NewHTTPHandlers(todoList)
	httpServer := httpp.NewHTTPServer(httpHandlers)

	if err := httpServer.StartServer(); err != nil {
		fmt.Println("failed to start http server:", err)
	}
}
