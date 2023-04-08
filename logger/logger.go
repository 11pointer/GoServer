package logger

import "go.uber.org/zap"

var LOG = getLogger()

func getLogger() *zap.Logger {

	res, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return res
}
