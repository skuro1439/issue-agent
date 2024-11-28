package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

func NewDefaultLogger() Logger {
	opt := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if lvl := os.Getenv("LOG_LEVEL"); lvl == "debug" {
		opt.Level = slog.LevelDebug
	}

	l := slog.New(slog.NewTextHandler(os.Stdout, opt))
	return DefaultLogger{
		logger: *l,
	}
}

type DefaultLogger struct {
	logger slog.Logger
}

func (l DefaultLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l DefaultLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l DefaultLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

type logLevel int

const (
	Debug logLevel = 30
	Info  logLevel = 10
	Error logLevel = 20
)

type Printer struct {
	level logLevel
}

func NewPrinter() Logger {
	level := Info
	if lvl := os.Getenv("LOG_LEVEL"); lvl == "debug" {
		level = Debug
	}
	return Printer{level}
}

func (l Printer) Info(msg string, args ...any) {
	if l.level >= Info {
		fmt.Printf(msg, args...)
	}
}

func (l Printer) Error(msg string, args ...any) {
	if l.level >= Error {
		fmt.Printf(msg, args...)
	}
}

func (l Printer) Debug(msg string, args ...any) {
	if l.level >= Debug {
		fmt.Printf(msg, args...)
	}
}
