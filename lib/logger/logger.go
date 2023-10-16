// logger 包提供了日志记录功能，基于 zap 实现。支持同时输出到控制台和文件，
// 文件按日期分割，日志格式为 JSON 格式。
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

type DateFileWriter struct {
	sync.Mutex
	logDir      string
	currentFile *os.File
}

// Global Logger
var Log *Logger

type Logger struct {
	logger *zap.Logger
}

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
	logger := zap.New(
		zapcore.NewTee(consoleCore, fileCore),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	// Global Logger
	Log = &Logger{logger: logger}
	return nil
}

// MustInit initializes a single instance of Logger and panic on failure.
// Accepts a path as the log file storage directory.
//
//	package main
//	func main() {logger.MustInit("log/")}
func MustInit(path string) {
	err := NewLogger(path)
	if err != nil {
		panic("Failed to initialize logger" + err.Error())
	}
}

func Debug(msg string, fields ...zap.Field) {
	Log.logger.Debug(msg, fields...)
}

func DebugT(msg string, traceID string, fields ...zap.Field) {
	fields = append(fields, zap.String("traceID", traceID))
	Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Log.logger.Info(msg, fields...)
}

func InfoT(msg string, traceID string, fields ...zap.Field) {
	fields = append(fields, zap.String("traceID", traceID))
	Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.logger.Warn(msg, fields...)
}

func WarnT(msg string, traceID string, fields ...zap.Field) {
	fields = append(fields, zap.String("traceID", traceID))
	Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.logger.Error(msg, fields...)
}

func ErrorT(msg string, traceID string, fields ...zap.Field) {
	fields = append(fields, zap.String("traceID", traceID))
	Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.logger.Fatal(msg, fields...)
}

func FatalT(msg string, traceID string, fields ...zap.Field) {
	fields = append(fields, zap.String("traceID", traceID))
	Fatal(msg, fields...)
}

func Close() error {
	return Log.logger.Sync()
}

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

func (w *DateFileWriter) Sync() error {
	return w.currentFile.Sync()
}

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
