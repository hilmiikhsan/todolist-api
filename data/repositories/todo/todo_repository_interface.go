package todo

import (
	"context"
	"todolist-api/data/models"
	"todolist-api/infra/db"

	"github.com/jmoiron/sqlx"
)

type TodoRepositoryInterface interface {
	CreateTodo(ctx context.Context, tx *sqlx.Tx, data models.Todo) (models.Todo, error)
	GetAllTodo(ctx context.Context) ([]models.Todo, error)
	GetOneTodo(ctx context.Context, tx *sqlx.Tx, id int) (models.Todo, error)
	UpdateTodo(ctx context.Context, tx *sqlx.Tx, id int, data models.Todo) error
	DeleteTodo(ctx context.Context, tx *sqlx.Tx, id int) error
}

func NewTodoRepository(db *db.DB) TodoRepositoryInterface {
	return &todoRepository{
		db,
	}
}
