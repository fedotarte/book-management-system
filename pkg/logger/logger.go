package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

type Logger struct {
	*zap.SugaredLogger
}

var sugarLogger *Logger

// NewLogger создает экземпляр кастомного логгера
func NewLogger() *Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.OutputPaths = []string{"stdout"}      // INFO в stdout
	config.ErrorOutputPaths = []string{"stderr"} // Ошибки и варнинги в stderr

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}

	return &Logger{logger.Sugar()} // Приводим к нашему типу Logger
}

// InitLogger инициализирует глобальный логгер
func InitLogger() {
	sugarLogger = NewLogger()
}

// GetLogger возвращает экз. логгера
func GetLogger() *Logger {
	if sugarLogger == nil {
		InitLogger()
	}
	return sugarLogger
}
