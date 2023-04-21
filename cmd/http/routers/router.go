package routers

import (
	"net/http"
	"todolist-api/cmd/http/handlers/activity"
	"todolist-api/cmd/http/handlers/todo"
	"todolist-api/utils"

	"github.com/gorilla/mux"
)

const (
	// GET for requesting data from a specified source or server
	GET = "GET"
	// POS for make requests that usually contain a body and send data to a server
	POS = "POST"
	// PUT for  creates a new resource or replaces a representation of the target resource with the request payload
	PUT = "PUT"
	// DEL for request method deletes the specified resource
	DEL = "DELETE"
)

// InitialRouter for object routers
func InitialRouter(
	activityHandler activity.ActivityHandlerInterface,
	todoHandler todo.TodoHandlerInterface,
) *mux.Router {
	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/activity-groups" {
			res := utils.SetResponseErrURLNotFound(r.URL.Path + " not found")
			res.JSONErrURLNotFound(w)
			return
		}
	})
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/todo-items" {
			res := utils.SetResponseErrURLNotFound(r.URL.Path + " not found")
			res.JSONErrURLNotFound(w)
			return
		}
	})

	// activity
	r.HandleFunc("/activity-groups", activityHandler.CreateActivity).Methods(POS)
	r.HandleFunc("/activity-groups", activityHandler.GetAllActivity).Methods(GET)
	r.HandleFunc("/activity-groups/{id}", activityHandler.GetOneActivity).Methods(GET)
	r.HandleFunc("/activity-groups/{id}", activityHandler.UpdateActivity).Methods(PUT)
	r.HandleFunc("/activity-groups/{id}", activityHandler.DeleteActivity).Methods(DEL)

	// todo
	r.HandleFunc("/todo-items", todoHandler.CreateTodo).Methods(POS)
	r.HandleFunc("/todo-items", todoHandler.GetAllTodo).Methods(GET)
	r.HandleFunc("/todo-items/{id}", todoHandler.GetOneTodo).Methods(GET)
	r.HandleFunc("/todo-items/{id}", todoHandler.UpdateTodo).Methods(PUT)
	r.HandleFunc("/todo-items/{id}", todoHandler.DeleteTodo).Methods(DEL)

	return r
}
