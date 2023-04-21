package activity

import (
	"context"
	"todolist-api/constants"
	"todolist-api/data/models"
	"todolist-api/infra/context/repository"
	"todolist-api/infra/errors"
	"todolist-api/objects/activity"
)

type activityService struct {
	*repository.RepoCtx
}

func (a activityService) CreateActivity(ctx context.Context, req activity.CreateActivity) (activity.Activity, error) {
	tx, err := a.DB.Begin(ctx)
	if err != nil {
		return activity.Activity{}, errors.Wrap(constants.ErrBeginTransaction)
	}

	if req.Title == "" {
		_ = tx.Rollback()
		return activity.Activity{}, errors.Wrap(constants.ErrTitleCannotBeNull)
	}

	activityID, err := a.ActivityRepository.CreateActivity(ctx, tx, models.Activity{
		Title: req.Title,
		Email: req.Email,
	})
	if err != nil {
		_ = tx.Rollback()
		return activity.Activity{}, err
	}

	data, err := a.ActivityRepository.GetOneActivity(ctx, tx, activityID.ActivityID)
	if err != nil {
		_ = tx.Rollback()
		return activity.Activity{}, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return activity.Activity{}, err
	}

	return activity.Activity{
		ID:        data.ActivityID,
		Title:     data.Title,
		Email:     data.Email,
		CreatedAt: data.CreatedAt.UTC().Format(constants.DateTimeFormat),
		UpdatedAt: data.UpdatedAt.UTC().Format(constants.DateTimeFormat),
	}, nil
}

func (a activityService) GetAllActivity(ctx context.Context) ([]activity.Activity, error) {
	tmpActivityData := []activity.Activity{}

	data, err := a.ActivityRepository.GetAllActivity(ctx)
	if err != nil {
		return tmpActivityData, err
	}

	for _, x := range data {
		tmpActivityData = append(tmpActivityData, activity.Activity{
			ID:        x.ActivityID,
			Title:     x.Title,
			Email:     x.Email,
			CreatedAt: x.CreatedAt.UTC().Format(constants.DateTimeFormat),
			UpdatedAt: x.UpdatedAt.UTC().Format(constants.DateTimeFormat),
		})
	}

	return tmpActivityData, nil
}

func (a activityService) GetOneActivity(ctx context.Context, id int) (activity.Activity, error) {
	tx, err := a.DB.Begin(ctx)
	if err != nil {
		return activity.Activity{}, errors.Wrap(constants.ErrBeginTransaction)
	}

	data, err := a.ActivityRepository.GetOneActivity(ctx, tx, id)
	if err != nil {
		_ = tx.Rollback()
		return activity.Activity{}, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return activity.Activity{}, err
	}

	return activity.Activity{
		ID:        data.ActivityID,
		Title:     data.Title,
		Email:     data.Email,
		CreatedAt: data.CreatedAt.UTC().Format(constants.DateTimeFormat),
		UpdatedAt: data.UpdatedAt.UTC().Format(constants.DateTimeFormat),
	}, nil
}

func (a activityService) UpdateActivity(ctx context.Context, id int, req activity.UpdateActivity) (activity.Activity, error) {
	tx, err := a.DB.Begin(ctx)
	if err != nil {
		return activity.Activity{}, errors.Wrap(constants.ErrBeginTransaction)
	}

	if req.Title == "" {
		_ = tx.Rollback()
		return activity.Activity{}, errors.Wrap(constants.ErrTitleCannotBeNull)
	}

	err = a.ActivityRepository.UpdateActivity(ctx, tx, id, models.Activity{
		Title: req.Title,
	})
	if err != nil {
		_ = tx.Rollback()
		return activity.Activity{}, err
	}

	data, err := a.ActivityRepository.GetOneActivity(ctx, tx, id)
	if err != nil {
		_ = tx.Rollback()
		return activity.Activity{}, err
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return activity.Activity{}, err
	}

	return activity.Activity{
		ID:        data.ActivityID,
		Title:     data.Title,
		Email:     data.Email,
		CreatedAt: data.CreatedAt.UTC().Format(constants.DateTimeFormat),
		UpdatedAt: data.UpdatedAt.UTC().Format(constants.DateTimeFormat),
	}, nil
}

func (a activityService) DeleteActivity(ctx context.Context, id int) error {
	tx, err := a.DB.Begin(ctx)
	if err != nil {
		return errors.Wrap(constants.ErrBeginTransaction)
	}

	data, err := a.ActivityRepository.GetOneActivity(ctx, tx, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = a.ActivityRepository.DeleteActivity(ctx, tx, data.ActivityID)
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
