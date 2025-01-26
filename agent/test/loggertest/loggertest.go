package loggertest

import "github.com/clover0/issue-agent/logger"

type testLogger struct{}

func NewTestLogger() logger.Logger {
	return &testLogger{}
}

func (l *testLogger) Info(msg string, args ...any)  {}
func (l *testLogger) Error(msg string, args ...any) {}
func (l *testLogger) Debug(msg string, args ...any) {}
