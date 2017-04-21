package logger

import (
	log "github.com/Sirupsen/logrus"
	"github.com/lavrs/docker-monitoring-service/pkg/context"
)

var ctx = context.Get()

func Info(l ...interface{}) {
	if ctx.Debug {
		log.Info(l)
	}
}

func Panic(l ...interface{}) {
	if ctx.Debug {
		log.Panic(l)
	}
}
