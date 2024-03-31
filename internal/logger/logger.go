package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	Sl *slog.Logger
}

func NewLogger() *Logger {
	sl := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	return &Logger{Sl: sl}
}
