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

func NewDefaultLogger(level string) Logger {
	opt := &slog.HandlerOptions{
		Level: slogLevel(level),
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

type Printer struct {
	level slog.Level
}

func NewPrinter(levelStr string) Logger {
	level := slogLevel(levelStr)

	return Printer{level}
}

func (l Printer) Debug(msg string, args ...any) {
	if l.level <= slog.LevelDebug {
		fmt.Printf(msg, args...)
	}
}
func (l Printer) Info(msg string, args ...any) {
	if l.level <= slog.LevelInfo {
		fmt.Printf(msg, args...)
	}
}

func (l Printer) Error(msg string, args ...any) {
	if l.level <= slog.LevelError {
		fmt.Printf(msg, args...)
	}
}

func slogLevel(l string) slog.Level {
	switch l {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "error":
		return slog.LevelError
	case "":
		return slog.LevelInfo
	default:
		panic(fmt.Sprintf("unknown log level: %s", l))
	}
}
