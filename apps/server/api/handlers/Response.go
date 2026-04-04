package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
}

func SendResponse(w http.ResponseWriter, statusCode int, success bool, data interface{}, message string, error string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{Success: success, Data: data, Message: message, Error: error})
}
