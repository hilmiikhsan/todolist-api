package activity

import (
	"context"
	"todolist-api/infra/context/repository"
	"todolist-api/objects/activity"
)

type ActivityServiceInterface interface {
	CreateActivity(ctx context.Context, req activity.CreateActivity) (activity.Activity, error)
	GetAllActivity(ctx context.Context) ([]activity.Activity, error)
	GetOneActivity(ctx context.Context, id int) (activity.Activity, error)
	UpdateActivity(ctx context.Context, id int, req activity.UpdateActivity) (activity.Activity, error)
	DeleteActivity(ctx context.Context, id int) error
}

func NewActivityService(ctx *repository.RepoCtx) ActivityServiceInterface {
	return &activityService{
		ctx,
	}
}
