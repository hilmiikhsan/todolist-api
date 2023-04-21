package todo

import (
	"context"
	"todolist-api/infra/context/repository"
	"todolist-api/objects/todo"
)

type TodoServiceInterface interface {
	CreateTodo(ctx context.Context, req todo.CreateTodo) (todo.Todo, error)
	GetAllTodo(ctx context.Context) ([]todo.Todo, error)
	GetOneTodo(ctx context.Context, id int) (todo.Todo, error)
	UpdateTodo(ctx context.Context, id int, req todo.UpdateTodo) (todo.Todo, error)
	DeleteTodo(ctx context.Context, id int) error
}

func NewTodoService(ctx *repository.RepoCtx) TodoServiceInterface {
	return &todoService{
		ctx,
	}
}
