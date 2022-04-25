package fnlog

import (
	"fmt"
	"log"
)

const (
	LogLevelTrace = 6000
	LogLevelDebug = 5000
	LogLevelInfo  = 4000
	LogLevelWarn  = 3000
	LogLevelError = 2000
	LogLevelFatal = 1000
)

var logLevel_ = LogLevelInfo

func SetLogLevel(logLevel int) {
	logLevel_ = logLevel
}

var (
	logLevelMap = map[int]string{
		LogLevelTrace: "TRACE",
		LogLevelDebug: "DEBUG",
		LogLevelInfo:  "INFO",
		LogLevelWarn:  "WARN",
		LogLevelError: "ERROR",
		LogLevelFatal: "FATAL",
	}
)

func Printf(logLevel int, format string, args ...interface{}) {
	if logLevel_ >= logLevel {
		s := fmt.Sprintf(format, args...)
		s2 := fmt.Sprintf("[%v] %v", logLevelMap[logLevel], s)
		log.Default().Output(3, s2)
	}
}

func Tracef(format string, args ...interface{}) {
	if logLevel_ >= LogLevelTrace {
		s := fmt.Sprintf(format, args...)
		s2 := fmt.Sprintf("[TRACE] %v", s)
		log.Default().Output(2, s2)
	}
}

func Debugf(format string, args ...interface{}) {
	if logLevel_ >= LogLevelDebug {
		s := fmt.Sprintf(format, args...)
		s2 := fmt.Sprintf("[DEBUG] %v", s)
		log.Default().Output(2, s2)
	}
}

func Infof(format string, args ...interface{}) {
	if logLevel_ >= LogLevelInfo {
		s := fmt.Sprintf(format, args...)
		s2 := fmt.Sprintf("[INFO] %v", s)
		log.Default().Output(2, s2)
	}
}

func Warnf(format string, args ...interface{}) {
	if logLevel_ >= LogLevelWarn {
		s := fmt.Sprintf(format, args...)
		s2 := fmt.Sprintf("[WARN] %v", s)
		log.Default().Output(2, s2)
	}
}

func Errorf(format string, args ...interface{}) {
	if logLevel_ >= LogLevelError {
		s := fmt.Sprintf(format, args...)
		s2 := fmt.Sprintf("[ERROR] %v", s)
		log.Default().Output(2, s2)
	}
}

func Fatalf(format string, args ...interface{}) {
	if logLevel_ >= LogLevelFatal {
		s := fmt.Sprintf(format, args...)
		s2 := fmt.Sprintf("[FATAL] %v", s)
		log.Default().Output(2, s2)
		panic(s2) // panic
	}
}
