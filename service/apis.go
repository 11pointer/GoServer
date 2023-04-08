package service

import (
	"11pointer/logger"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type addUserRequest struct {
	Email string `json:"email"`
}

func AddUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.LOG.Error("Invalid method for AddUser request", zap.Any("requestMethod", r.Method))
		return
	}

	request := &addUserRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		logger.LOG.Error("Error in unmarshalling AddUser request", zap.Any("requestBody", r.Body))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = addUser(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
