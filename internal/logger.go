package logger

import (
	"fmt"
)

type Logger struct {
}

func NewLogger() Logger {

	return Logger{}
}

func (l *Logger) Message(msg string) {

	fmt.Println(msg)
}
