package logger

import (
	"os"

	"go.uber.org/zap"
)

// NewLogger returns a new zap logger
func NewLogger() *zap.Logger {
	const (
		logPath = "./logs/library.log"
	)
	os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout", logPath}
	l, err := c.Build()
	if err != nil {
		panic(err)
	}
	return l
}
