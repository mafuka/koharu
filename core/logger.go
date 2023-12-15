package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

// LogLevel defines the level of logging.
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// LoggerConfig holds Logger configurations.
type LoggerConfig struct {
	File     string   `yaml:"file"`     // Log file path
	Level    LogLevel `yaml:"level"`    // Log Level
	MaxDays  int      `yaml:"max_days"` // Max days to rotate logs
	Compress bool     `yaml:"compress"` // Compress logs using gzip
}

// DefaultLoggerConfig provides a basic default LogCfg.
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		File:     "console",
		Level:    DebugLevel,
		MaxDays:  3,
		Compress: false,
	}
}

// loggerInst is the application-wide logger instance.
var (
	loggerInst *Logger
	once       sync.Once
)

// RegisterLogger initializes the global logger.
func RegisterLogger(cfg LoggerConfig) {
	once.Do(func() {
		loggerInst = buildLogger(cfg)
	})
}

func Log() *Logger {
	return loggerInst
}

// Logger wraps zap.SugaredLogger to provide formatted logging capabilities.
type Logger struct {
	*zap.SugaredLogger
}

// buildLogger creates a new Logger instance.
func buildLogger(cfg LoggerConfig) *Logger {
	level := zapLevel(cfg.Level)

	cores := []zapcore.Core{
		zapcore.NewCore(newConsoleEncoder(), newConsoleWriter(), level),
	}

	if cfg.File != "console" {
		fileWriter := newFileWriter(cfg)
		cores = append(cores, zapcore.NewCore(newJSONEncoder(), fileWriter, level))
	}

	combinedCore := zapcore.NewTee(cores...)
	zapLogger := zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(1))
	return &Logger{zapLogger.Sugar()}
}

// zapLevel converts LogLevel to zap's logging level.
func zapLevel(level LogLevel) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// newJSONEncoder prepares the JSON encoder for the logger.
func newJSONEncoder() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(cfg)
}

// newConsoleEncoder prepares the console encoder for the logger.
func newConsoleEncoder() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.CapitalColorLevelEncoder,
		EncodeTime:   zapcore.TimeEncoderOfLayout("15:04:05.000"),
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(cfg)
}

// newFileWriter sets up Lumberjack as the file writer.
func newFileWriter(cfg LoggerConfig) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename: cfg.File,
		MaxAge:   cfg.MaxDays,
		Compress: cfg.Compress,
	})
}

// newConsoleWriter returns a console writer.
func newConsoleWriter() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
}

// Debug logs a message at the debug level.
func (l *Logger) Debug(format string, a ...interface{}) {
	l.SugaredLogger.Debugf(format, a...)
}

// Info logs a message at the info level.
func (l *Logger) Info(format string, a ...interface{}) {
	l.SugaredLogger.Infof(format, a...)
}

// Warn logs a message at the warn level.
func (l *Logger) Warn(format string, a ...interface{}) {
	l.SugaredLogger.Warnf(format, a...)
}

// Error logs a message at the error level.
func (l *Logger) Error(format string, a ...interface{}) {
	l.SugaredLogger.Errorf(format, a...)
}

// Fatal logs a message at the fatal level and typically causes the program to terminate.
func (l *Logger) Fatal(format string, a ...interface{}) {
	l.SugaredLogger.Fatalf(format, a...)
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}

// Update re-initializes the global logger with new configuration.
func (l *Logger) Update(cfg LoggerConfig) {
	loggerInst = buildLogger(cfg)
}
