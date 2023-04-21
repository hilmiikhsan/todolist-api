package activity

import (
	"context"
	"time"
	"todolist-api/data/models"
	"todolist-api/infra/db"
	"todolist-api/infra/errors"
	"todolist-api/utils"

	"github.com/jmoiron/sqlx"
)

type activityRepository struct {
	db *db.DB
}

func (a activityRepository) CreateActivity(ctx context.Context, tx *sqlx.Tx, data models.Activity) (models.Activity, error) {
	result, err := tx.ExecContext(
		ctx,
		queryCreateActivity,
		data.Title,
		data.Email,
		time.Now(),
	)
	if err != nil {
		return models.Activity{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Activity{}, err
	}

	data.ActivityID = int(id)

	return data, nil
}

func (a activityRepository) GetAllActivity(ctx context.Context) ([]models.Activity, error) {
	results := []models.Activity{}
	err := a.db.Slave().SelectContext(
		ctx,
		&results,
		queryGetAllActivity,
	)
	if err != nil {
		return results, err
	}

	return results, nil
}

func (a activityRepository) GetOneActivity(ctx context.Context, tx *sqlx.Tx, id int) (models.Activity, error) {
	results := []models.Activity{}
	err := tx.SelectContext(
		ctx,
		&results,
		queryGetOneActivity,
		id,
	)
	if err != nil {
		return models.Activity{}, err
	}

	if len(results) == 0 {
		return models.Activity{}, errors.Wrap(utils.ErrDataNotFound(id))
	}

	return results[0], nil
}

func (a activityRepository) UpdateActivity(ctx context.Context, tx *sqlx.Tx, id int, data models.Activity) error {
	result, err := tx.ExecContext(
		ctx,
		queryUpdateActivity,
		data.Title,
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}

	activityID, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(utils.ErrDataNotFound(activityID))
	}

	return nil
}

func (a activityRepository) DeleteActivity(ctx context.Context, tx *sqlx.Tx, id int) error {
	_, err := tx.ExecContext(
		ctx,
		queryDeleteActivity,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}
