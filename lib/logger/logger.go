// Package logger offers structured logging with file rotation based on date using zap.
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DateFileWriter handles daily log file rotation.
type DateFileWriter struct {
	sync.Mutex
	logDir      string
	currentFile *os.File
}

// Log is a singleton instance for application-wide logging.
var Log *Logger

// Logger wraps zap's SugaredLogger for convenience.
type Logger struct {
	sugaredLogger *zap.SugaredLogger
}

// NewLogger configures and instantiates a Logger with file and console outputs.
func NewLogger(outputPath string) error {
	writer, err := NewDateFileWriter(filepath.Dir(outputPath))
	if err != nil {
		return err
	}

	encoderConfig := zapcore.EncoderConfig{
		LevelKey:   "level",
		TimeKey:    "time",
		CallerKey:  "caller",
		MessageKey: "msg",
		// FunctionKey:   "func",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime:    zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line))
		},
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	consoleLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.ErrorLevel
	})

	fileLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.ErrorLevel
	})

	consoleOutput := zapcore.Lock(os.Stdout)
	fileOutput := zapcore.AddSync(writer)

	consoleCore := zapcore.NewCore(consoleEncoder, consoleOutput, consoleLevel)
	fileCore := zapcore.NewCore(fileEncoder, fileOutput, fileLevel)

	// Create the logger with the async core
	baseLogger := zap.New(
		zapcore.NewTee(consoleCore, fileCore),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	// Create a sugared logger
	sugaredLogger := baseLogger.Sugar()

	// Global Logger
	Log = &Logger{sugaredLogger: sugaredLogger}
	return nil
}

// NewDateFileWriter prepares a directory for log files and initializes a writer.
func NewDateFileWriter(logDir string) (*DateFileWriter, error) {
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		return nil, err
	}

	writer := &DateFileWriter{
		logDir: logDir,
	}

	if err := writer.openNewFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

// Write outputs log entries to the current file, rotating if necessary.
func (w *DateFileWriter) Write(p []byte) (n int, err error) {
	w.Lock()
	defer w.Unlock()

	if err := w.checkDateChange(); err != nil {
		return 0, err
	}

	n, err = w.currentFile.Write(p)
	if err != nil {
		return n, err
	}

	err = w.Sync()
	if err != nil {
		return n, err
	}

	return n, nil
}

// Sync flushes the current log file to disk.
func (w *DateFileWriter) Sync() error {
	return w.currentFile.Sync()
}

// checkDateChange handles the rotation of the log file if the date has changed.
func (w *DateFileWriter) checkDateChange() error {
	now := time.Now()
	date := now.Format("2006-01-02")

	currentFileDate := filepath.Base(w.currentFile.Name())
	if currentFileDate == date {
		return nil
	}

	if err := w.openNewFile(); err != nil {
		return err
	}

	return nil
}

// openNewFile opens a new log file based on the current date.
func (w *DateFileWriter) openNewFile() error {
	now := time.Now()
	date := now.Format("2006-01-02")
	logFilePath := filepath.Join(w.logDir, fmt.Sprintf("%s.log", date))

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if w.currentFile != nil {
		w.currentFile.Close()
	}

	w.currentFile = file

	return nil
}

// MustInit is a convenience function to initialize logging and panic on error.
func MustInit(path string) {
	err := NewLogger(path)
	if err != nil {
		panic("Failed to initialize logger" + err.Error())
	}
}

// Debug logs a formatted debug message.
func Debug(template string, args ...interface{}) {
	Log.sugaredLogger.Debugf(template, args...)
}

// Info logs a formatted informational message.
func Info(template string, args ...interface{}) {
	Log.sugaredLogger.Infof(template, args...)
}

// Warn logs a formatted warning message.
func Warn(template string, args ...interface{}) {
	Log.sugaredLogger.Warnf(template, args...)
}

// Error logs a formatted error message.
func Error(template string, args ...interface{}) {
	Log.sugaredLogger.Errorf(template, args...)
}

// Fatal logs a formatted fatal message and terminates the application.
func Fatal(template string, args ...interface{}) {
	Log.sugaredLogger.Fatalf(template, args...)
}

// Close flushes any buffered log entries.
func Close() error {
	return Log.sugaredLogger.Sync()
}
