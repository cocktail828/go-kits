package slogx_test

import (
	"testing"

	"github.com/cocktail828/go-kits/pkg"
	"github.com/cocktail828/go-kits/pkg/slogx"
)

func TestSlog(t *testing.T) {
	l := slogx.NewLoggerWithLumberjack(pkg.Log{
		Level:      "error",
		Filename:   "/log/server/error.log",
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	})
	l = l.With("a1", "b1")
	l.Info("finished", "key", "value")
	l.Info("finishedxxx", "key", "value")
}
