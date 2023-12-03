package logger

import (
	"log"
	"os"
	"strings"
)

var packageName string

type IF2SLogger interface {
	Debug()
	Info()
	Warn()
	Error()
}

type F2SLogger struct {
	ComponentName string
	Logger        *log.Logger
}

func Initialize(ComponentName string) F2SLogger {
	return F2SLogger{
		ComponentName: ComponentName,
		Logger:        log.New(os.Stdout, "component=\""+ComponentName+"\" ", log.LstdFlags),
	}
}

func (l F2SLogger) Info(text ...string) {
	l.Logger.Printf("loglevel=info msg=\"%s\"", strings.Join(text, " "))
}
func (l F2SLogger) Debug(text ...string) {
	l.Logger.Printf("loglevel=debug msg=\"%s\"", strings.Join(text, " "))
}
func (l F2SLogger) Warn(text ...string) {
	l.Logger.Printf("loglevel=warn msg=\"%s\"", strings.Join(text, " "))
}
func (l F2SLogger) Error(err error) {
	if err != nil {
		l.Logger.Printf("loglevel=error msg=\"%s\"", err.Error())
	}
}
