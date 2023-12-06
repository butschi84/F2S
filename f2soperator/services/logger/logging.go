package logger

import (
	"log/slog"
	"os"
)

type IF2SLogger interface {
	Debug()
	Info()
	Warn()
	Error()
}

type F2SLogger struct {
	ComponentName string
	Logger        *slog.Logger
}

func Initialize(ComponentName string) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil)).With("component", ComponentName)
}
