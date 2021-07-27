package log

import (
	"context"
	"expinc/sunagent/common"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var logger = logrus.New()
var mutex = sync.Mutex{}
var fileOutput *lumberjack.Logger = nil

func init() {
	logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		TimestampFormat: time.RFC3339Nano,
	})
}

func SetLevel(level string) {
	switch {
	case level == "debug":
		logger.SetLevel((logrus.DebugLevel))
	case level == "info":
		logger.SetLevel((logrus.InfoLevel))
	case level == "warn":
		logger.SetLevel((logrus.WarnLevel))
	case level == "error":
		logger.SetLevel((logrus.ErrorLevel))
	case level == "fatal":
		logger.SetLevel((logrus.FatalLevel))
	}
}

func SetRotateFileOutput(fileName string, filelimitmb int) {
	mutex.Lock()
	defer mutex.Unlock()

	if nil != fileOutput {
		fileOutput.Close()
	}

	fileOutput = &lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   filelimitmb,
		LocalTime: true,
	}

	logger.SetOutput(fileOutput)
}

func Debug(obj interface{}) {
	logger.Debug(obj)
}

func Info(obj interface{}) {
	logger.Info(obj)
}

func Warn(obj interface{}) {
	logger.Warn(obj)
}

func Error(obj interface{}) {
	logger.Error(obj)
}

func Fatal(obj interface{}) {
	logger.Fatal(obj)
}

func DebugCtx(ctx context.Context, obj interface{}) {
	logger.WithField(common.TraceIdContextKey, ctx.Value(common.TraceIdContextKey)).Debug(obj)
}

func InfoCtx(ctx context.Context, obj interface{}) {
	logger.WithField(common.TraceIdContextKey, ctx.Value(common.TraceIdContextKey)).Info(obj)
}

func WarnCtx(ctx context.Context, obj interface{}) {
	logger.WithField(common.TraceIdContextKey, ctx.Value(common.TraceIdContextKey)).Warn(obj)
}

func ErrorCtx(ctx context.Context, obj interface{}) {
	logger.WithField(common.TraceIdContextKey, ctx.Value(common.TraceIdContextKey)).Error(obj)
}

func FatalCtx(ctx context.Context, obj interface{}) {
	logger.WithField(common.TraceIdContextKey, ctx.Value(common.TraceIdContextKey)).Fatal(obj)
}
