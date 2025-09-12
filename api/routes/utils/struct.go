package utils

import (
	"encoding/json"
	"net/http"
)

// JSON RESPONSES

func WriteJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func WriteSuccessResponse(w http.ResponseWriter, s string) {
	WriteJSONResponse(w, NewSuccessMessage(s), http.StatusOK)
}

func WriteConflictResponse(w http.ResponseWriter, s string) {
	WriteJSONResponse(w, NewConflictMessage(s), http.StatusConflict)
}

func WriteInternalServerErrorResponse(w http.ResponseWriter, s string) {
	WriteJSONResponse(w, NewInternalServerErrorMessage(s), http.StatusInternalServerError)
}

// MESSAGE RESPONSES:

type Response struct {
	Message string `json:"message"`
}

func NewSuccessMessage(s string) *Response {
	return &Response{Message: s}
}

func NewConflictMessage(s string) *Response {
	return &Response{Message: s}
}

func NewInternalServerErrorMessage(s string) *Response {
	return &Response{Message: s}
}
