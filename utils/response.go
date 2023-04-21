package utils

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	contentType              = "Content-Type"
	contentTypeValue         = "application/json; charset=utf-8"
	xContentTypeOptions      = "X-Content-Type-Options"
	xContentTypeOptionsValue = "nosniff"

	STATUS_INTERNAL_ERR = "STATUS_INTERNAL_ERROR"
	STATUS_BAD_REQUEST  = "STATUS_BAD_REQUEST"
	STATUS_UNAUTHORIZED = "STATUS_UNAUTHORIZED"
	STATUS_FORBIDDEN    = "STATUS_FORBIDDEN"
	STATUS_NOT_FOUND    = "STATUS_NOT_FOUND"

	MESSAGE_SUCCESS             = "Success"
	MESSAGE_BAD_REQUEST         = "Bad Request"
	MESSAGE_INTERNAL_SERVER_ERR = "Internal Server Error"
	MESSAGE_NOT_FOUND           = "Not Found"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseErr struct {
	Status  interface{} `json:"status"`
	Message string      `json:"message"`
}

type ResponseErrNotFound struct {
	Message string `json:"message"`
}

func SetResponseJSON(status, message string, data interface{}) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func SetResponseErrJSON(status interface{}, message string) *ResponseErr {
	return &ResponseErr{
		Status:  status,
		Message: message,
	}
}

func SetResponseErrURLNotFound(message string) *ResponseErrNotFound {
	return &ResponseErrNotFound{
		Message: message,
	}
}

func SetResponseErrNotFound(status interface{}, message string) *ResponseErr {
	return &ResponseErr{
		Status:  status,
		Message: message,
	}
}

func (r *Response) JSONResponse(w http.ResponseWriter) {
	w.Header().Set(contentType, contentTypeValue)
	w.Header().Set(xContentTypeOptions, xContentTypeOptionsValue)
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Error(err)
	}
}

func (r *Response) JSONSuccessResponse(w http.ResponseWriter) {
	w.Header().Set(contentType, contentTypeValue)
	w.Header().Set(xContentTypeOptions, xContentTypeOptionsValue)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Error(err)
	}
}

func (r *ResponseErr) JSONErrResponse(w http.ResponseWriter) {
	w.Header().Set(contentType, contentTypeValue)
	w.Header().Set(xContentTypeOptions, xContentTypeOptionsValue)
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Error(err)
	}
}

func (r *ResponseErr) JSONErrInternalServerResponse(w http.ResponseWriter) {
	w.Header().Set(contentType, contentTypeValue)
	w.Header().Set(xContentTypeOptions, xContentTypeOptionsValue)
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Error(err)
	}
}

func (r *ResponseErrNotFound) JSONErrURLNotFound(w http.ResponseWriter) {
	w.Header().Set(contentType, contentTypeValue)
	w.Header().Set(xContentTypeOptions, xContentTypeOptionsValue)
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Error(err)
	}
}

func (r *ResponseErr) JSONErrNotFound(w http.ResponseWriter) {
	w.Header().Set(contentType, contentTypeValue)
	w.Header().Set(xContentTypeOptions, xContentTypeOptionsValue)
	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Error(err)
	}
}
