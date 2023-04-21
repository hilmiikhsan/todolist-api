package activity

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"todolist-api/constants"
	"todolist-api/infra/context/service"
	"todolist-api/objects/activity"

	"todolist-api/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

type activityHandler struct {
	*service.Ctx
}

func (a activityHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var req activity.CreateActivity
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error(err)
		res := utils.SetResponseErrJSON(http.StatusBadRequest, err.Error())
		res.JSONErrResponse(w)
		return
	}

	if err = validator.Validate(req); err != nil {
		log.Error(err)
		res := utils.SetResponseErrJSON(http.StatusBadRequest, err.Error())
		res.JSONErrResponse(w)
		return
	}

	data, err := a.ActivityService.CreateActivity(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrTitleCannotBeNull.Error()) {
			log.Error(err)
			res := utils.SetResponseErrJSON(utils.MESSAGE_BAD_REQUEST, err.Error())
			res.JSONErrResponse(w)
			return
		}
		res := utils.SetResponseErrJSON(utils.MESSAGE_INTERNAL_SERVER_ERR, err.Error())
		res.JSONErrInternalServerResponse(w)
		return
	}

	res := utils.SetResponseJSON(utils.MESSAGE_SUCCESS, "Success", data)
	res.JSONResponse(w)
}

func (a activityHandler) GetAllActivity(w http.ResponseWriter, r *http.Request) {
	data, err := a.ActivityService.GetAllActivity(r.Context())
	if err != nil {
		res := utils.SetResponseErrJSON(utils.MESSAGE_INTERNAL_SERVER_ERR, err.Error())
		res.JSONErrInternalServerResponse(w)
		return
	}

	res := utils.SetResponseJSON(utils.MESSAGE_SUCCESS, "Success", data)
	res.JSONSuccessResponse(w)
}

func (a activityHandler) GetOneActivity(w http.ResponseWriter, r *http.Request) {
	queryParamID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(queryParamID)
	if err != nil {
		if strings.Contains(err.Error(), ":id") {
			log.Error(err)
			res := utils.SetResponseErrNotFound(utils.MESSAGE_NOT_FOUND, utils.ErrDataNotFound(":id").Error())
			res.JSONErrNotFound(w)
			return
		}
		log.Error(err)
		res := utils.SetResponseErrJSON(utils.MESSAGE_BAD_REQUEST, err.Error())
		res.JSONErrResponse(w)
		return
	}

	data, err := a.ActivityService.GetOneActivity(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			log.Error(err)
			res := utils.SetResponseErrNotFound(utils.MESSAGE_NOT_FOUND, err.Error())
			res.JSONErrNotFound(w)
			return
		}
		res := utils.SetResponseErrJSON(utils.MESSAGE_INTERNAL_SERVER_ERR, err.Error())
		res.JSONErrInternalServerResponse(w)
		return
	}

	res := utils.SetResponseJSON(utils.MESSAGE_SUCCESS, "Success", data)
	res.JSONSuccessResponse(w)
}

func (a activityHandler) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	queryParamID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(queryParamID)
	if err != nil {
		if strings.Contains(err.Error(), ":id") {
			log.Error(err)
			res := utils.SetResponseErrNotFound(utils.MESSAGE_NOT_FOUND, utils.ErrDataNotFound(":id").Error())
			res.JSONErrNotFound(w)
			return
		}
		log.Error(err)
		res := utils.SetResponseErrJSON(utils.MESSAGE_BAD_REQUEST, err.Error())
		res.JSONErrResponse(w)
		return
	}

	var req activity.UpdateActivity
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error(err)
		res := utils.SetResponseErrJSON(http.StatusBadRequest, err.Error())
		res.JSONErrResponse(w)
		return
	}

	if err = validator.Validate(req); err != nil {
		log.Error(err)
		res := utils.SetResponseErrJSON(http.StatusBadRequest, err.Error())
		res.JSONErrResponse(w)
		return
	}

	data, err := a.ActivityService.UpdateActivity(r.Context(), id, req)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrTitleCannotBeNull.Error()) {
			log.Error(err)
			res := utils.SetResponseErrJSON(utils.MESSAGE_BAD_REQUEST, err.Error())
			res.JSONErrResponse(w)
			return
		}

		if strings.Contains(err.Error(), "Not Found") {
			log.Error(err)
			res := utils.SetResponseErrNotFound(utils.MESSAGE_NOT_FOUND, err.Error())
			res.JSONErrNotFound(w)
			return
		}
		res := utils.SetResponseErrJSON(utils.MESSAGE_INTERNAL_SERVER_ERR, err.Error())
		res.JSONErrInternalServerResponse(w)
		return
	}

	res := utils.SetResponseJSON(utils.MESSAGE_SUCCESS, "Success", data)
	res.JSONSuccessResponse(w)
}

func (a activityHandler) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	queryParamID := mux.Vars(r)["id"]
	id, err := strconv.Atoi(queryParamID)
	if err != nil {
		if strings.Contains(err.Error(), ":id") {
			log.Error(err)
			res := utils.SetResponseErrNotFound(utils.MESSAGE_NOT_FOUND, utils.ErrDataNotFound(":id").Error())
			res.JSONErrNotFound(w)
			return
		}
		log.Error(err)
		res := utils.SetResponseErrJSON(utils.MESSAGE_BAD_REQUEST, err.Error())
		res.JSONErrResponse(w)
		return
	}

	err = a.ActivityService.DeleteActivity(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			log.Error(err)
			res := utils.SetResponseErrNotFound(utils.MESSAGE_NOT_FOUND, err.Error())
			res.JSONErrNotFound(w)
			return
		}
		res := utils.SetResponseErrJSON(utils.MESSAGE_INTERNAL_SERVER_ERR, err.Error())
		res.JSONErrInternalServerResponse(w)
		return
	}

	data := make(map[string]interface{})

	res := utils.SetResponseJSON(utils.MESSAGE_SUCCESS, "Success", data)
	res.JSONSuccessResponse(w)
}
