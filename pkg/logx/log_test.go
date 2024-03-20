package logx_test

import (
	"testing"

	"github.com/cocktail828/go-kits/pkg"
	"github.com/cocktail828/go-kits/pkg/logx"
)

func TestLogger(t *testing.T) {
	l := logx.NewLoggerWithLumberjack(pkg.Log{
		Level:      "debug",
		Filename:   "/log/server/error.log",
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     1,
		Compress:   false,
	})
	l.Debug("Debug msg") // not be written
	l.Info("Info msg")   // not be written
	l.Warn("Warn msg")   // written in general.log
	l.Error("Error msg") // written in error.log
}
