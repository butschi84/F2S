package logger

import (
	"fmt"
	"log"
	"os"
)

var packageName string

func Initialize(packageName string) *log.Logger {

	packageNameLength := 20
	packageNamePadded := fmt.Sprintf("%-*s", packageNameLength, packageName)

	var logger *log.Logger = log.New(os.Stdout, "["+packageNamePadded+"] ", log.LstdFlags)
	return logger
}
