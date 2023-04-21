package todo

import (
	"context"
	"todolist-api/constants"
	"todolist-api/data/models"
	"todolist-api/infra/context/repository"
	"todolist-api/infra/errors"
	"todolist-api/objects/todo"
)

type todoService struct {
	*repository.RepoCtx
}

func (t todoService) CreateTodo(ctx context.Context, req todo.CreateTodo) (todo.Todo, error) {
	tx, err := t.DB.Begin(ctx)
	if err != nil {
		return todo.Todo{}, errors.Wrap(constants.ErrBeginTransaction)
	}

	if req.Title == "" {
		_ = tx.Rollback()
		return todo.Todo{}, errors.Wrap(constants.ErrTitleCannotBeNull)
	}

	if req.Priority == "" {
		req.Priority = constants.Priority
	}

	todoID, err := t.TodoRepository.CreateTodo(ctx, tx, models.Todo{
		Title:           req.Title,
		ActivityGroupID: req.ActivityGroupID,
		IsActive:        req.IsActive,
		Priority:        req.Priority,
	})
	if err != nil {
		_ = tx.Rollback()
		return todo.Todo{}, err
	}

	data, err := t.TodoRepository.GetOneTodo(ctx, tx, todoID.TodoID)
	if err != nil {
		_ = tx.Rollback()
		return todo.Todo{}, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return todo.Todo{}, err
	}

	return todo.Todo{
		ID:              data.TodoID,
		Title:           data.Title,
		ActivityGroupID: data.ActivityGroupID,
		IsActive:        data.IsActive,
		Priority:        data.Priority,
		CreatedAt:       data.CreatedAt.UTC().Format(constants.DateTimeFormat),
		UpdatedAt:       data.UpdatedAt.UTC().Format(constants.DateTimeFormat),
	}, nil
}

func (t todoService) GetAllTodo(ctx context.Context) ([]todo.Todo, error) {
	tmpTodoData := []todo.Todo{}

	data, err := t.TodoRepository.GetAllTodo(ctx)
	if err != nil {
		return tmpTodoData, err
	}

	for _, x := range data {
		tmpTodoData = append(tmpTodoData, todo.Todo{
			ID:              x.TodoID,
			Title:           x.Title,
			ActivityGroupID: x.ActivityGroupID,
			IsActive:        x.IsActive,
			Priority:        x.Priority,
			UpdatedAt:       x.UpdatedAt.UTC().Format(constants.DateTimeFormat),
			CreatedAt:       x.CreatedAt.UTC().Format(constants.DateTimeFormat),
		})
	}

	return tmpTodoData, nil
}

func (t todoService) GetOneTodo(ctx context.Context, id int) (todo.Todo, error) {
	tx, err := t.DB.Begin(ctx)
	if err != nil {
		return todo.Todo{}, errors.Wrap(constants.ErrBeginTransaction)
	}

	data, err := t.TodoRepository.GetOneTodo(ctx, tx, id)
	if err != nil {
		_ = tx.Rollback()
		return todo.Todo{}, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return todo.Todo{}, err
	}

	return todo.Todo{
		ID:              data.TodoID,
		Title:           data.Title,
		ActivityGroupID: data.ActivityGroupID,
		IsActive:        data.IsActive,
		Priority:        data.Priority,
		UpdatedAt:       data.UpdatedAt.UTC().Format(constants.DateTimeFormat),
		CreatedAt:       data.CreatedAt.UTC().Format(constants.DateTimeFormat),
	}, nil
}

func (t todoService) UpdateTodo(ctx context.Context, id int, req todo.UpdateTodo) (todo.Todo, error) {
	tx, err := t.DB.Begin(ctx)
	if err != nil {
		return todo.Todo{}, errors.Wrap(constants.ErrBeginTransaction)
	}

	if req.Title == "" {
		_ = tx.Rollback()
		return todo.Todo{}, errors.Wrap(constants.ErrTitleCannotBeNull)
	}

	if req.Priority == "" {
		req.Priority = constants.Priority
	}

	err = t.TodoRepository.UpdateTodo(ctx, tx, id, models.Todo{
		Title:    req.Title,
		IsActive: req.IsActive,
		Priority: req.Priority,
	})
	if err != nil {
		_ = tx.Rollback()
		return todo.Todo{}, err
	}

	data, err := t.TodoRepository.GetOneTodo(ctx, tx, id)
	if err != nil {
		_ = tx.Rollback()
		return todo.Todo{}, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return todo.Todo{}, err
	}

	return todo.Todo{
		ID:              data.TodoID,
		Title:           data.Title,
		ActivityGroupID: data.ActivityGroupID,
		IsActive:        data.IsActive,
		Priority:        data.Priority,
		UpdatedAt:       data.UpdatedAt.UTC().Format(constants.DateTimeFormat),
		CreatedAt:       data.CreatedAt.UTC().Format(constants.DateTimeFormat),
	}, nil
}

func (t todoService) DeleteTodo(ctx context.Context, id int) error {
	tx, err := t.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(constants.ErrBeginTransaction)
	}

	data, err := t.TodoRepository.GetOneTodo(ctx, tx, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = t.TodoRepository.DeleteTodo(ctx, tx, data.TodoID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
