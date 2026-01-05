package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type Logger struct {
	level     LogLevel
	prefix    string
	debugLog  *log.Logger
	infoLog   *log.Logger
	warnLog   *log.Logger
	errorLog  *log.Logger
	fatalLog  *log.Logger
}

func NewLogger(prefix string, level LogLevel) *Logger {
	flags := log.LstdFlags | log.Lshortfile

	return &Logger{
		level:    level,
		prefix:   prefix,
		debugLog: log.New(os.Stdout, fmt.Sprintf("[DEBUG] %s: ", prefix), flags),
		infoLog:  log.New(os.Stdout, fmt.Sprintf("[INFO] %s: ", prefix), flags),
		warnLog:  log.New(os.Stderr, fmt.Sprintf("[WARN] %s: ", prefix), flags),
		errorLog: log.New(os.Stderr, fmt.Sprintf("[ERROR] %s: ", prefix), flags),
		fatalLog: log.New(os.Stderr, fmt.Sprintf("[FATAL] %s: ", prefix), flags),
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.debugLog.Printf(msg, args...)
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.infoLog.Printf(msg, args...)
	}
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.warnLog.Printf(msg, args...)
	}
}

func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.errorLog.Printf(msg, args...)
	}
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.fatalLog.Fatalf(msg, args...)
}

func (l *Logger) WithContext(context string) *ContextLogger {
	return &ContextLogger{
		logger:  l,
		context: context,
	}
}

type ContextLogger struct {
	logger  *Logger
	context string
}

func (cl *ContextLogger) Debug(msg string, args ...interface{}) {
	fullMsg := fmt.Sprintf("[%s] %s", cl.context, msg)
	cl.logger.Debug(fullMsg, args...)
}

func (cl *ContextLogger) Info(msg string, args ...interface{}) {
	fullMsg := fmt.Sprintf("[%s] %s", cl.context, msg)
	cl.logger.Info(fullMsg, args...)
}

func (cl *ContextLogger) Warn(msg string, args ...interface{}) {
	fullMsg := fmt.Sprintf("[%s] %s", cl.context, msg)
	cl.logger.Warn(fullMsg, args...)
}

func (cl *ContextLogger) Error(msg string, args ...interface{}) {
	fullMsg := fmt.Sprintf("[%s] %s", cl.context, msg)
	cl.logger.Error(fullMsg, args...)
}

type TimingLogger struct {
	logger *Logger
	name   string
	start  time.Time
}

func (l *Logger) StartTiming(name string) *TimingLogger {
	return &TimingLogger{
		logger: l,
		name:   name,
		start:  time.Now(),
	}
}

func (tl *TimingLogger) End() {
	duration := time.Since(tl.start)
	tl.logger.Info("%s completed in %v", tl.name, duration)
}

func (tl *TimingLogger) EndWithValue(value interface{}) {
	duration := time.Since(tl.start)
	tl.logger.Info("%s completed in %v with result: %v", tl.name, duration, value)
}
