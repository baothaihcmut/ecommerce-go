package logger

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/libs/pkg/constant"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Logger interface to support different log levels
type Logger interface {
	WithCtx(context.Context) Logger
	Debug(args map[string]any, msg string)
	Debugf(args map[string]any, template string, format ...interface{})
	Info(args map[string]any, msg string)
	Infof(args map[string]any, template string, format ...interface{})
	Warn(args map[string]any, msg string)
	Warnf(args map[string]any, template string, format ...interface{})
	Error(args map[string]any, msg string)
	Errorf(args map[string]any, template string, format ...interface{})
	Fatal(args map[string]any, msg string)
	Fatalf(args map[string]any, template string, format ...interface{})
	Panic(args map[string]any, msg string)
	Panicf(args map[string]any, template string, format ...interface{})
}

// logger is the concrete implementation of the Logger interface
type logger struct {
	l *logrus.Logger
	requestId uuid.UUID
}



// NewLogger initializes and returns a logger instance
func NewLogger(l *logrus.Logger) Logger {
	return &logger{
		l: l,
		requestId: uuid.Nil,
	}
}


func (l *logger) WithCtx(ctx context.Context) Logger {
	requestIdTarget := uuid.Nil
	if requestId, ok :=ctx.Value(constant.RequestIdKey).(uuid.UUID); ok {
		requestIdTarget = requestId
	}
	return &logger{
		l: l.l,
		requestId: requestIdTarget,
	}
}

// Debug level logging with context
func (l *logger) Debug(args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Debug(msg)
}

func (l *logger) Debugf(args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Debugf(template, format...)
}

// Info level logging with context
func (l *logger) Info(args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Info(msg)
}

func (l *logger) Infof(args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Infof(template, format...)
}

// Warn level logging with context
func (l *logger) Warn(args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Warn(msg)
}

func (l *logger) Warnf(args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Warnf(template, format...)
}

// Error level logging with context
func (l *logger) Error(args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Error(msg)
}

func (l *logger) Errorf(args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Errorf(template, format...)
}

// Fatal level logging with context
func (l *logger) Fatal( args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Fatal(msg)
}

func (l *logger) Fatalf(args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Fatalf(template, format...)
}

// Panic level logging with context
func (l *logger) Panic(args map[string]any, msg string) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Panic(msg)
}

func (l *logger) Panicf(args map[string]any, template string, format ...interface{}) {
	if args == nil {
		args = make(map[string]any)
	}
	if l.requestId != uuid.Nil {
		args["request_id"] = l.requestId
	}
	l.l.WithFields(logrus.Fields(args)).Panicf(template, format...)
}
