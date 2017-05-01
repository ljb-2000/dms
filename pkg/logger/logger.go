package logger

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lavrs/docker-monitoring-service/pkg/context"
)

// Info print info log
func Info(l ...interface{}) {
	if context.Get().Debug {
		log.Info(l)
	}
}

// Panic print panic log
func Panic(l ...interface{}) {
	log.Panic(l)
}
