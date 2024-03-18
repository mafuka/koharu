package bot

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger interface defines structured logging capabilities with different severity levels.
type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

// globalLogger is the application-wide Logger instance.
var (
	globalLogger Logger
	once         sync.Once
)

// SetLogger sets the global logger instance.
// It's thread-safe and can be called only once. Subsequent calls will have no effect.
func SetLogger(l Logger) {
	once.Do(func() {
		globalLogger = l
	})
}

// GetLogger returns the globally set Logger instance.
// It's recommended to call SetLogger during application initialization.
func GetLogger() Logger {
	return globalLogger
}

// Log is a shortcut to GetLogger.
func Log() Logger {
	return globalLogger
}

// ZapConfig holds configuration for the ZapLogger.
type ZapConfig struct {
	File     string   `yaml:"file"`     // File specifies the log file path. Use "console" to output to stdout.
	Level    ZapLevel `yaml:"level"`    // Level is the logging level (e.g., debug, info).
	MaxDays  int      `yaml:"max_days"` // MaxDays is the maximum number of days to retain old log files.
	Compress bool     `yaml:"compress"` // Compress determines if the log rotation should compress log files.
}

// ZapLogger wraps zap.SugaredLogger to satisfy the Logger interface.
type ZapLogger struct {
	*zap.SugaredLogger
}

// ZapLevel defines logging levels supported by the ZapLogger.
type ZapLevel string

const (
	DebugLevel ZapLevel = "debug"
	InfoLevel  ZapLevel = "info"
	WarnLevel  ZapLevel = "warn"
	ErrorLevel ZapLevel = "error"
	FatalLevel ZapLevel = "fatal"
)

// zLevel converts a ZapLevel to zap's logging level.
func zLevel(l ZapLevel) zapcore.Level {
	switch l {
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
		return zapcore.InfoLevel // Default logging level is Info.
	}
}

// NewZapLogger creates a new ZapLogger instance based on the provided configuration.
func NewZapLogger(cfg ZapConfig) *ZapLogger {
	level := zLevel(cfg.Level)

	cores := []zapcore.Core{}

	if cfg.File == "console" {
		cores = append(cores, zapcore.NewCore(newZapConsoleEncoder(), newZapConsoleWriter(), level))
	} else {
		cores = append(cores, zapcore.NewCore(newZapJSONEncoder(), newZapFileWriter(cfg), level))
	}

	combinedCore := zapcore.NewTee(cores...)
	zapLogger := zap.New(combinedCore, zap.AddCaller(), zap.AddCallerSkip(0))
	return &ZapLogger{zapLogger.Sugar()}
}

// newZapJSONEncoder creates a new JSON encoder for zap logging.
func newZapJSONEncoder() zapcore.Encoder {
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

// newZapConsoleEncoder creates a new console encoder for zap logging.
func newZapConsoleEncoder() zapcore.Encoder {
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

// newZapFileWriter sets up Lumberjack as the file writer for zap logging.
func newZapFileWriter(cfg ZapConfig) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename: cfg.File,
		MaxAge:   cfg.MaxDays,
		Compress: cfg.Compress,
	})
}

// newZapConsoleWriter returns a console writer for zap logging.
func newZapConsoleWriter() zapcore.WriteSyncer {
	return zapcore.Lock(os.Stdout)
}
