package log

import (
	"github.com/senayuki/carrier/pkg/consts"
	"go.uber.org/zap"
)

var loggerBase *zap.Logger

func init() {
	var err error
	loggerBase, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func Logger(component string) *zap.Logger {
	return loggerBase.With(zap.String(consts.Component, component))
}

func Sync() {
	loggerBase.Sync()
}
