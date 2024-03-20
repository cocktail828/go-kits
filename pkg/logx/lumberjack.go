package logx

import (
	"io"

	"github.com/cocktail828/go-kits/pkg"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Lumberjack struct {
	logger    *lumberjack.Logger
	formatter logrus.Formatter
}

func (hook *Lumberjack) Fire(entry *logrus.Entry) error {
	msg, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.logger.Write([]byte(msg))
	return err
}

func (hook *Lumberjack) Levels() []logrus.Level {
	return logrus.AllLevels
}

func levelC(level string) logrus.Level {
	switch level {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.ErrorLevel
	}
}

func NewLoggerWithLumberjack(logf pkg.Log) Logger {
	return NewLoggerWithLogrus(func() *logrus.Logger {
		l := logrus.New()
		l.SetOutput(io.Discard)
		l.SetLevel(levelC(logf.Level))
		l.AddHook(&Lumberjack{
			logger: &lumberjack.Logger{
				Filename:   logf.Filename,
				MaxSize:    logf.MaxSize,
				MaxBackups: logf.MaxBackups,
				MaxAge:     logf.MaxAge,
				Compress:   logf.Compress,
			},
			formatter: &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"},
		})
		return l
	}())
}
