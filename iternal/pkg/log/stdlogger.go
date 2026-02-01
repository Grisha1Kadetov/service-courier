package log

import (
	"fmt"
	"log"
	"strings"
)

type StdLogger struct{}

func (l *StdLogger) Debug(msg string, fields ...Field) {
	log.Println("[Debug]", msg, formatFields(fields))
}

func (l *StdLogger) Info(msg string, fields ...Field) {
	log.Println("[Info]", msg, formatFields(fields))
}

func (l *StdLogger) Warn(msg string, fields ...Field) {
	log.Println("[Warn]", msg, formatFields(fields))
}

func (l *StdLogger) Error(msg string, fields ...Field) {
	log.Println("[Error]", msg, formatFields(fields))
}

func (l *StdLogger) Fatal(msg string, fields ...Field) {
	log.Fatal("[Fatal]", msg, formatFields(fields))
}

func NewStdLogger() Logger {
	return &StdLogger{}
}

func formatFields(fields []Field) string {
	var b strings.Builder
	for i, f := range fields {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%s=%v", f.Key, f.Value)
	}
	return b.String()
}
