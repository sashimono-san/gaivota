package log

import (
	"log"
	"os"
	"strings"

	"github.com/leoschet/gaivota"
)

func New(prefix string) *Logger {
	return &Logger{
		logger: log.New(os.Stdout, prefix, log.LstdFlags),
	}
}

type Logger struct {
	logger *log.Logger
}

func (l *Logger) Log(level gaivota.LogLevel, format string, v ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}

	if level == gaivota.LogLevelFatal {
		l.logger.Fatalf(format, v)
	}

	go l.logger.Printf(format, v)
}
