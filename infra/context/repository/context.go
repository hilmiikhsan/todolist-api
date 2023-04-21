package repository

import (
	"todolist-api/data/repositories/activity"
	"todolist-api/data/repositories/todo"
	"todolist-api/infra/db"
)

// RepoCtx struct for repository context
type RepoCtx struct {
	DB                 *db.DB
	ActivityRepository activity.ActivityRepositoryInterface
	TodoRepository     todo.TodoRepositoryInterface
}
