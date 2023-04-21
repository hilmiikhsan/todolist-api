package activity

import (
	"context"
	"todolist-api/data/models"
	"todolist-api/infra/db"

	"github.com/jmoiron/sqlx"
)

type ActivityRepositoryInterface interface {
	CreateActivity(ctx context.Context, tx *sqlx.Tx, data models.Activity) (models.Activity, error)
	GetAllActivity(ctx context.Context) ([]models.Activity, error)
	GetOneActivity(ctx context.Context, tx *sqlx.Tx, id int) (models.Activity, error)
	UpdateActivity(ctx context.Context, tx *sqlx.Tx, id int, data models.Activity) error
	DeleteActivity(ctx context.Context, tx *sqlx.Tx, id int) error
}

func NewActivityRepository(db *db.DB) ActivityRepositoryInterface {
	return &activityRepository{
		db,
	}
}
