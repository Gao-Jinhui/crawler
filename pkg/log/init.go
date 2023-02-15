package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger() *zap.Logger {
	// log
	plugin := NewStdoutPlugin(zapcore.InfoLevel)
	logger := NewLogger(plugin)
	logger.Info("log init end")
	return logger
}
