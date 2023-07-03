package fare

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func SetLoggerLevel(level zapcore.Level) {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	l, _ := config.Build()
	logger = l.Sugar()
}
