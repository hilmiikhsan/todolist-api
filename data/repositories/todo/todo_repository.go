package todo

import (
	"context"
	"time"
	"todolist-api/data/models"
	"todolist-api/infra/db"
	"todolist-api/infra/errors"
	"todolist-api/utils"

	"github.com/jmoiron/sqlx"
)

type todoRepository struct {
	db *db.DB
}

func (t todoRepository) CreateTodo(ctx context.Context, tx *sqlx.Tx, data models.Todo) (models.Todo, error) {
	result, err := tx.ExecContext(
		ctx,
		queryCreateTodo,
		data.Title,
		data.ActivityGroupID,
		data.IsActive,
		data.Priority,
		time.Now(),
	)
	if err != nil {
		return models.Todo{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Todo{}, err
	}

	data.TodoID = int(id)

	return data, nil
}

func (t todoRepository) GetAllTodo(ctx context.Context) ([]models.Todo, error) {
	results := []models.Todo{}
	err := t.db.Slave().SelectContext(
		ctx,
		&results,
		queryGetAllTodo,
	)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (t todoRepository) GetOneTodo(ctx context.Context, tx *sqlx.Tx, id int) (models.Todo, error) {
	results := []models.Todo{}
	err := tx.SelectContext(
		ctx,
		&results,
		queryGetOneTodo,
		id,
	)
	if err != nil {
		return models.Todo{}, err
	}

	if len(results) == 0 {
		return models.Todo{}, errors.Wrap(utils.ErrDataNotFound(id))
	}

	return results[0], nil
}

func (t todoRepository) UpdateTodo(ctx context.Context, tx *sqlx.Tx, id int, data models.Todo) error {
	result, err := tx.ExecContext(
		ctx,
		queryUpdateTodo,
		data.Title,
		data.IsActive,
		data.Priority,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}

	todoID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(utils.ErrDataNotFound(todoID))
	}

	return nil
}

func (t todoRepository) DeleteTodo(ctx context.Context, tx *sqlx.Tx, id int) error {
	_, err := tx.ExecContext(
		ctx,
		queryDeleteTodo,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
