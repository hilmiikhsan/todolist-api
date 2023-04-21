package todo

import (
	"net/http"
	"todolist-api/infra/context/service"
)

type TodoHandlerInterface interface {
	CreateTodo(w http.ResponseWriter, r *http.Request)
	GetAllTodo(w http.ResponseWriter, r *http.Request)
	GetOneTodo(w http.ResponseWriter, r *http.Request)
	UpdateTodo(w http.ResponseWriter, r *http.Request)
	DeleteTodo(w http.ResponseWriter, r *http.Request)
}

func NewTodoHandler(serviceCtx *service.Ctx) TodoHandlerInterface {
	return &todoHandler{
		serviceCtx,
	}
}
