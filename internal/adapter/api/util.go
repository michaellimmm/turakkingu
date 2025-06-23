package api

import (
	"encoding/json"
	"net/http"
)

func sendJson(w http.ResponseWriter, statusCode int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func sendError(w http.ResponseWriter, statusCode int, err error) error {
	body := ErrorResponse{ErrorMessage: err.Error()}
	return sendJson(w, statusCode, body)
}
