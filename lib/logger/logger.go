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

type Logger struct {
	logger *zap.Logger
}

func NewLogger(outputPath string) (*Logger, error) {
	writer, err := NewDateFileWriter(filepath.Dir(outputPath))
	if err != nil {
		return nil, err
	}

	encoderConfig := zapcore.EncoderConfig{
		LevelKey:    "level",
		TimeKey:     "time",
		CallerKey:   "caller",
		MessageKey:  "msg",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(fmt.Sprintf("%s:%d:%s", filepath.Base(caller.File), caller.Line, caller.Function))
		},
	}

	// consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleEncoder := zapcore.NewJSONEncoder(encoderConfig)
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

	consoleLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})

	fileLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.InfoLevel
	})

	consoleOutput := zapcore.Lock(os.Stdout)
	fileOutput := zapcore.AddSync(writer)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleOutput, consoleLevel),
		zapcore.NewCore(fileEncoder, fileOutput, fileLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return &Logger{logger: logger}, nil
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Close() error {
	return l.logger.Sync()
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
	// fmt.Println("Opening new file:", logFilePath)

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
