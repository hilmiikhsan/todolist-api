package service

import (
	"todolist-api/cmd/services/activity"
	"todolist-api/cmd/services/todo"
)

// Ctx service context
type Ctx struct {
	ActivityService activity.ActivityServiceInterface
	TodoService     todo.TodoServiceInterface
}
