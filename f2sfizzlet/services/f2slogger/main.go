package f2slogger

import (
	"fmt"
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
	PackageName string
	Logger      *log.Logger
}

func Initialize(PackageName string) F2SLogger {
	packageNameLength := 20
	packageNamePadded := fmt.Sprintf("%-*s", packageNameLength, PackageName+"]")

	return F2SLogger{
		PackageName: packageNamePadded,
		Logger:      log.New(os.Stdout, "["+packageNamePadded+" ", log.LstdFlags),
	}
}

func (l F2SLogger) Info(text ...string) {
	l.Logger.Printf("[INFO] %s", strings.Join(text, " "))
}
func (l F2SLogger) Debug(text ...string) {
	l.Logger.Printf("[DEBUG] %s", strings.Join(text, " "))
}
func (l F2SLogger) Warn(text ...string) {
	l.Logger.Printf("[WARN] %s", strings.Join(text, " "))
}
func (l F2SLogger) Error(err error) {
	if err != nil {
		l.Logger.Printf("[ERROR] %s", err.Error())
	}
}
