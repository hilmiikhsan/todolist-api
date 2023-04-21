package activity

import (
	"net/http"
	"todolist-api/infra/context/service"
)

type ActivityHandlerInterface interface {
	CreateActivity(w http.ResponseWriter, r *http.Request)
	GetAllActivity(w http.ResponseWriter, r *http.Request)
	GetOneActivity(w http.ResponseWriter, r *http.Request)
	UpdateActivity(w http.ResponseWriter, r *http.Request)
	DeleteActivity(w http.ResponseWriter, r *http.Request)
}

func NewActivityHandler(serviceCtx *service.Ctx) ActivityHandlerInterface {
	return &activityHandler{
		serviceCtx,
	}
}
