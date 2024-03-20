package slogx

import (
	"github.com/cocktail828/go-kits/pkg"
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func levelC(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

func NewLoggerWithLumberjack(logf pkg.Log) Logger {
	var lvl slog.LevelVar
	lvl.Set(levelC(logf.Level))

	return NewLoggerWithSlog(slog.New(slog.NewJSONHandler(&lumberjack.Logger{
		Filename:   logf.Filename,
		MaxSize:    logf.MaxSize,
		MaxBackups: logf.MaxBackups,
		MaxAge:     logf.MaxAge,
		Compress:   logf.Compress,
	}, &slog.HandlerOptions{Level: &lvl})))
}
