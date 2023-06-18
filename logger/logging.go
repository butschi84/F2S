package logger

import (
	"log"
	"os"
)

var packageName string

func Initialize(packageName string) *log.Logger {

	var logger *log.Logger = log.New(os.Stdout, "["+packageName+"] ", log.LstdFlags)
	return logger
}
