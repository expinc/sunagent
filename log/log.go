package log

import (
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
