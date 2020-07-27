package logger

import (
	"log"
	"strings"
)

// LoggingLevel type used in logging level configuration and filtration.
type LoggingLevel string

const (
	LevelDebug LoggingLevel = "debug"
	LevelInfo  LoggingLevel = "info"
	LevelWarn  LoggingLevel = "warn"
	LevelError LoggingLevel = "fatal"
)

var loggingLevelToInt = map[LoggingLevel]int{
	LevelDebug: 0,
	LevelInfo:  1,
	LevelWarn:  2,
	LevelError: 3,
}
var loggingLevelArray = []LoggingLevel{
	LevelDebug,
	LevelInfo,
	LevelWarn,
	LevelError,
}

var loggingLevel = 1

// SetLevel configures logger level.
func SetLevel(level LoggingLevel) {
	loggingLevel = loggingLevelToInt[level]
}

// GetLevel returns logger level.
func GetLevel() LoggingLevel {
	return loggingLevelArray[loggingLevel]
}

// Info logs a message at level Info.
func Info(message string, args ...interface{}) {
	genericLogging(LevelInfo, message, args...)
}

// Infof logs a message at level Info.
func Infof(format string, args ...interface{}) {
	genericLogging(LevelInfo, format, args...)
}

// Debug logs a message at level Debug.
func Debug(message string, args ...interface{}) {
	genericLogging(LevelDebug, message, args...)
}

// Debugf logs a message at level Info.
func Debugf(format string, args ...interface{}) {
	genericLogging(LevelDebug, format, args...)
}

// Error logs a message at level Error.
func Error(message string, args ...interface{}) {
	genericLogging(LevelError, message, args...)
}

// Errorf logs a message at level Error.
func Errorf(format string, args ...interface{}) {
	genericLogging(LevelError, format, args...)
}

// Warn logs a message at level Info.
func Warn(message string, args ...interface{}) {
	genericLogging(LevelWarn, message, args...)
}

// Warnf logs a message at level Warn.
func Warnf(format string, args ...interface{}) {
	genericLogging(LevelWarn, format, args...)
}

// Fatal logs a message at level Fatal + os.exit(1).
func Fatal(message string, args ...interface{}) {
	log.Fatal(append([]interface{}{message}, args...)...)
}

// Fatalf logs a message at level Fatal + os.exit(1).
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func genericLogging(level LoggingLevel, message string, args ...interface{}) {
	if loggingLevel > loggingLevelToInt[level] {
		return
	}
	log.Printf(strings.ToUpper(string(level))+"\t"+message+"\n", args...)
}
