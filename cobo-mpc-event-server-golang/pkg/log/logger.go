package log

import (
	"fmt"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	DefaultLogger *logrus.Entry // the exposed logger (actually an logrus.Entry)

	defaultInnerLogger = logrus.New() // logger under the hood
	once               sync.Once
)

//nolint:nolintlint,gochecknoinits
func init() {
	once.Do(func() {
		defaultInnerLogger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		})
		defaultInnerLogger.SetLevel(logrus.DebugLevel)
		DefaultLogger = logrus.NewEntry(defaultInnerLogger)
	})
}

func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

func Debugln(args ...interface{}) {
	DefaultLogger.Debugln(args...)
}

func Debugf(format string, args ...interface{}) {
	DefaultLogger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Infoln(args ...interface{}) {
	DefaultLogger.Infoln(args...)
}

func Infof(format string, args ...interface{}) {
	DefaultLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

func Warnln(args ...interface{}) {
	DefaultLogger.Warnln(args...)
}

func Warnf(format string, args ...interface{}) {
	DefaultLogger.Warnf(format, args...)
}

// Find the first error in an arg list, and return a logger WithError'ed.
func getLoggerWithError(args []interface{}) logrus.FieldLogger {
	var tmpLogger logrus.FieldLogger = DefaultLogger
	for _, arg := range args {
		if err, ok := arg.(error); ok {
			tmpLogger = DefaultLogger.WithError(err)
			break
		}
	}
	return tmpLogger
}

func Error(args ...interface{}) {
	getLoggerWithError(args).Error(args...)
}

func Errorln(args ...interface{}) {
	getLoggerWithError(args).Errorln(args...)
}

func Errorf(format string, args ...interface{}) {
	getLoggerWithError(args).Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	getLoggerWithError(args).Fatal(args...)
}

func Fatalln(args ...interface{}) {
	getLoggerWithError(args).Fatalln(args...)
}

func Fatalf(format string, args ...interface{}) {
	getLoggerWithError(args).Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	getLoggerWithError(args).Panic(args...)
}

func Panicln(args ...interface{}) {
	getLoggerWithError(args).Panicln(args...)
}

func Panicf(format string, args ...interface{}) {
	getLoggerWithError(args).Panicf(format, args...)
}

func ErrorStack(args ...any) {
	getLoggerWithError(args).Error("panic %s\n%s", fmt.Sprint(args...), string(debug.Stack()))
}

func ErrorStackf(format string, args ...any) {
	getLoggerWithError(args).Errorf("%s\n%s", fmt.Sprintf(format, args...), string(debug.Stack()))
}

func WithError(err error) *logrus.Entry {
	return DefaultLogger.WithError(err)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return DefaultLogger.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return DefaultLogger.WithFields(fields)
}
