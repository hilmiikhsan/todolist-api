package todo

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"todolist-api/constants"
	"todolist-api/infra/context/service"
	"todolist-api/objects/todo"
	"todolist-api/utils"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

type todoHandler struct {
	*service.Ctx
}

func (t todoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req todo.CreateTodo
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

	data, err := t.TodoService.CreateTodo(r.Context(), req)
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

func (t todoHandler) GetAllTodo(w http.ResponseWriter, r *http.Request) {
	data, err := t.TodoService.GetAllTodo(r.Context())
	if err != nil {
		res := utils.SetResponseErrJSON(utils.MESSAGE_INTERNAL_SERVER_ERR, err.Error())
		res.JSONErrInternalServerResponse(w)
		return
	}

	res := utils.SetResponseJSON(utils.MESSAGE_SUCCESS, "Success", data)
	res.JSONSuccessResponse(w)
}

func (t todoHandler) GetOneTodo(w http.ResponseWriter, r *http.Request) {
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

	data, err := t.TodoService.GetOneTodo(r.Context(), id)
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

func (t todoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
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

	var req todo.UpdateTodo
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

	data, err := t.TodoService.UpdateTodo(r.Context(), id, req)
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

func (t todoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
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

	err = t.TodoService.DeleteTodo(r.Context(), id)
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
