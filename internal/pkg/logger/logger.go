package logger

import (
	"log"
	"os"
	"strings"
)

// LoggingLevel type used in logging level configuration and filtration.
type LoggingLevel string

// Known logging levels (except of "invalid level" defined to catch errors in logging level configuration).
const (
	LevelInvalid LoggingLevel = "invalid level"
	LevelDebug   LoggingLevel = "debug"
	LevelInfo    LoggingLevel = "info"
	LevelWarn    LoggingLevel = "warn"
	LevelError   LoggingLevel = "error"
	LevelFatal   LoggingLevel = "fatal"
)

var loggingLevelToInt = map[LoggingLevel]int{
	LevelDebug: 1,
	LevelInfo:  2,
	LevelWarn:  3,
	LevelError: 4,
	LevelFatal: 4,
}
var loggingLevelArray = []LoggingLevel{
	LevelInvalid, // stub item to validate proper logging level
	LevelDebug,
	LevelInfo,
	LevelWarn,
	LevelError,
	LevelFatal,
}

// Default logging level is set to "Info".
var loggingLevel = loggingLevelToInt[LevelInfo]

// SetLevel configures logger level.
func SetLevel(level LoggingLevel) {
	if loggingLevelToInt[level] == 0 {
		Fatalf("failed to configure logging level: unsupported logging level passed: %q. Supported logging levels: %q", level, validLoggingLevels())
	}
	loggingLevel = loggingLevelToInt[level]

	Debugf("configured logging level to %q", level)
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
	genericLogging(LevelFatal, message, args...)
	os.Exit(1)
}

// Fatalf logs a message at level Fatal + os.exit(1).
func Fatalf(format string, args ...interface{}) {
	genericLogging(LevelFatal, format, args...)
	os.Exit(1)
}

func genericLogging(level LoggingLevel, message string, args ...interface{}) {
	if loggingLevel > loggingLevelToInt[level] {
		return
	}
	log.Printf(strings.ToUpper(string(level))+"\t"+message+"\n", args...)
}

func validLoggingLevels() []string {
	result := make([]string, len(loggingLevelToInt))
	for k, v := range loggingLevelToInt {
		result[v-1] = string(k)
	}

	return result
}
