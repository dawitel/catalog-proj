package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

var defaultLogger *slog.Logger

func init() {
	defaultLogger = slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level: slog.LevelInfo,
	}))
}

func SetDefault(l *slog.Logger) {
	if l != nil {
		defaultLogger = l
	}
}

func Default() *slog.Logger {
	return defaultLogger
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}
