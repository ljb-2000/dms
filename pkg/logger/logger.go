package logger

import (
	"github.com/lavrs/dms/pkg/context"
	"go.uber.org/zap"
)

var sugar *zap.SugaredLogger

// init sugar logger
func init() {
	config := zap.NewDevelopmentConfig()
	config.DisableCaller = true

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	sugar = logger.Sugar()
}

// Info print info log
func Info(data ...interface{}) {
	if context.Get().Debug {
		sugar.Info(data)
	}
}

// Panic print panic log
func Panic(err ...interface{}) {
	sugar.Panic(err)
}
