package main

import (
	"11pointer/database"
	"11pointer/logger"
	"11pointer/service"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	database.InitDbDriver()
	logger.LOG.Info("Starting http server on port 80")
	http.HandleFunc("/addUser", service.AddUser)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		logger.LOG.Panic("Panic", zap.Error(err))
	}
}
