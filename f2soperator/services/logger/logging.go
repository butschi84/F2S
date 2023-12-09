package logger

import (
	"os"
	"strings"

	"golang.org/x/exp/slog"
)

type F2SLogger struct {
	ComponentName string
	Logger        *slog.Logger
}

func Initialize(ComponentName string) *F2SLogger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)).With("component", ComponentName)
	return &F2SLogger{
		ComponentName: ComponentName,
		Logger:        logger,
	}
}

// standard logging functions
func (l F2SLogger) Info(text ...string) {
	l.Logger.Info(strings.Join(text, " "), "type", "log")
}
func (l F2SLogger) Debug(text ...string) {
	l.Logger.Debug(strings.Join(text, " "), "type", "log")
}
func (l F2SLogger) Warn(text ...string) {
	l.Logger.Warn(strings.Join(text, " "), "type", "log")
}
func (l F2SLogger) Error(text ...string) {
	l.Logger.Error(strings.Join(text, " "), "type", "log")
}

// log an event
func (l F2SLogger) Event(text ...string) {
	l.Logger.Info(strings.Join(text, " "), "type", "event")
}
