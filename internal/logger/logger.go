package logger

import (
	"go.uber.org/zap"
	"log/slog"
	"os"
)

type Logger struct {
	Sl *slog.Logger
}

func NewLogger() *Logger {
	sl := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	sugar.Infow("hello from logger!")

	return &Logger{Sl: sl}
}
