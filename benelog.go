package benelog

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
)

type correlationIdType int

const (
	requestIdKey correlationIdType = iota
	sessionIdKey
)

var logger zap.Logger

func NewLogger(options ...zap.Option) (*zap.Logger, error) {
	logConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "ts",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     logConfig,
		OutputPaths:       nil,
		ErrorOutputPaths:  nil,
		InitialFields:     nil,
	}
	fields := zap.Fields(zap.Int("pid", os.Getpid()),
		zap.String("exe", path.Base(os.Args[0])))
	options = append(options, fields)

	return cfg.Build(options...)
}

func init() {
	tmpLogger, err := NewLogger()
	print("Hello friend")
	if err != nil {
		panic("error creating simple logger")
	}
	logger = *tmpLogger
}

// WithRqId returns a context which knows its request ID
func WithRqId(ctx context.Context, rqId string) context.Context {
	print("WithRqId\n")
	return context.WithValue(ctx, requestIdKey, rqId)
}

// WithSessionId returns a context which knows its session ID
func WithSessionId(ctx context.Context, sessionId string) context.Context {
	print("WithSessionId\n")
	return context.WithValue(ctx, sessionIdKey, sessionId)
}

// Logger returns a zap logger with as much context as possible
func Logger(ctx context.Context) zap.Logger {
	newLogger := logger
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {
			print("trying rqid\n")
			newLogger = *newLogger.With(zap.String("rqId", ctxRqId))
		}
		if ctxSessionId, ok := ctx.Value(sessionIdKey).(string); ok {
			print("trying sessionid\n")
			newLogger = *newLogger.With(zap.String("sessionId", ctxSessionId))
		}
	}
	fmt.Printf("%v\n", newLogger)
	return newLogger
}
