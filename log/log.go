package log

import (
	"fmt"
	"log"
	"os"

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
	msg := format
	if len(v) > 0 {
		msg = fmt.Sprintf(format, v)
	}

	if level == gaivota.LogLevelFatal {
		l.logger.Fatalln(msg)
	}

	go l.logger.Println(msg)
}
