package generator

import (
	"fmt"
	"io"
)

type Logger struct {
	w io.Writer
}

func NewLogger(w io.Writer) *Logger {
	return &Logger{
		w: w,
	}
}

func (l *Logger) Write(format string, args ...interface{}) {
	io.WriteString(l.w, fmt.Sprintf("[protoc-gen-graphql] "+format+"\n", args...)) // nolint: errcheck
}
