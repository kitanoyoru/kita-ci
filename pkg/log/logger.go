package log

import (
	"log"
	"strings"

	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(level string) *Logger {
	var logger *zap.Logger

	switch strings.ToLower(level) {
	case "debug":
		debugLogger, err := zap.NewDevelopment()
		if err != nil {
			log.Fatal("Can't initialize logger")
		}
		logger = debugLogger
	case "info":
		infoLogger, err := zap.NewProduction()
		if err != nil {
			log.Fatal("Can't initialize logger")
		}
		logger = infoLogger
	default:
		log.Fatalf("Please specify correct log level")
	}

	return &Logger{
		logger,
	}
}

func (l *Logger) Info(msg string) {
	defer l.logger.Sync()
	l.logger.Info(msg)
}

func (l *Logger) Error(msg string, err error) {
	defer l.logger.Sync()
	l.logger.Error(msg, zap.Error(err))
}

func (l *Logger) Fatal(msg string, err error) {
	defer l.logger.Sync()
	l.logger.Fatal(msg, zap.Error(err))
}
