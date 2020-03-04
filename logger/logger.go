package logger

import (
	"os"
	"runtime"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger *logrus.Logger = logrus.New()

// Instance return singleton instance of logger
func Instance() *logrus.Logger {
	return logger
}

func init() {
	if runtime.GOOS == "windows" {
		logger.Formatter = &logrus.TextFormatter{
			EnvironmentOverrideColors: true,
			ForceColors:               true,
		}
		os.Setenv("CLICOLOR_FORCE", "1")
		logger.Out = colorable.NewColorableStdout()
	} else {
		logger.Formatter = new(prefixed.TextFormatter)
	}
	logger.Level = logrus.DebugLevel
}

// Info output INFO log
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Infoln output INFO log with new line
func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

// Infof output INFO log with format
func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

// Debug output DEBUG log
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Debugln output DEBUG log with new line
func Debugln(args ...interface{}) {
	logger.Debugln(args...)
}

// Debugf output DEBUG log with format
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Error output ERROR log
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Errorln output ERROR log with new line
func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

// Errorf output ERROR log with format
func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

// Warn output ERROR log
func Warn(args ...interface{}) {
	logger.Warn(args...)
}

// Warnln output ERROR log with new line
func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

// Warnf output ERROR log with format
func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

// Fatal output FATAL log
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

// Fatalln output FATAL log with new line
func Fatalln(args ...interface{}) {
	logger.Fatalln(args...)
}

// Fatalf output FATAL log with format
func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
