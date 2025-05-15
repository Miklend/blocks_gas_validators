package logging

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type writerHook struct {
	Writer    []io.Writer
	Loglevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		w.Write([]byte(line))
	}
	return nil
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.Loglevels
}

var (
	e    *logrus.Entry
	once sync.Once
)

type Logger struct {
	*logrus.Entry
}

func InitLogger(levelStr string) *Logger {
	once.Do(func() {
		l := logrus.New()

		l.SetReportCaller(true)
		l.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				filename := path.Base(frame.File)
				return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
			},
			DisableColors: false,
			FullTimestamp: true,
		})

		err := os.MkdirAll("logs", 0755)
		if err != nil {
			panic(fmt.Errorf("failed to create logs dir: %w", err))
		}

		allFile, _ := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		errorFile, _ := os.OpenFile("logs/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

		l.SetOutput(io.Discard)

		l.AddHook(&writerHook{
			Writer:    []io.Writer{allFile, os.Stdout},
			Loglevels: logrus.AllLevels,
		})

		l.AddHook(&writerHook{
			Writer:    []io.Writer{errorFile},
			Loglevels: []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
		})

		level, err := logrus.ParseLevel(strings.ToLower(levelStr))
		if err != nil {
			level = logrus.InfoLevel
		}
		l.SetLevel(level)

		e = logrus.NewEntry(l)
	})
	return &Logger{e}
}

func GetLogger() *Logger {
	if e == nil {
		return InitLogger("info")
	}
	return &Logger{e}
}

func (l *Logger) GetLoggerWithField(k string, v interface{}) *Logger {
	return &Logger{l.WithField(k, v)}
}
