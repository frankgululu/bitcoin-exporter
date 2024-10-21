package logger

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true, //会输出log原文
	})

	Logger = slog.New(handler)
}
