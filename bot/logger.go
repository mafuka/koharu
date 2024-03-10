package bot

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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

// LogConfig holds Logger configurations.
type LogConfig struct {
	File     string   `yaml:"file"`     // Log file path, could be "console".
	Level    LogLevel `yaml:"level"`    // Log Level
	MaxDays  int      `yaml:"max_days"` // Max days to rotate logs
	Compress bool     `yaml:"compress"` // Compress logs using gzip
}

// loggerInst is the application-wide logger instance.
var (
	loggerInst *Logger
	once       sync.Once
)

// InitLogger initializes the global logger.
func InitLogger(cfg LogConfig) {
	once.Do(func() {
		loggerInst = NewLogger(cfg)
	})
}

func Log() *Logger {
	return loggerInst
}

// Logger wraps zap.SugaredLogger to provide formatted logging capabilities.
type Logger struct {
	*zap.SugaredLogger
}

// NewLogger creates a new Logger instance.
func NewLogger(cfg LogConfig) *Logger {
	level := zapLevel(cfg.Level)

	cores := []zapcore.Core{}

	if cfg.File == "console" {
		cores = append(cores, zapcore.NewCore(newConsoleEncoder(), newConsoleWriter(), level))
	} else {
		cores = append(cores, zapcore.NewCore(newJSONEncoder(), newFileWriter(cfg), level))
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
func newFileWriter(cfg LogConfig) zapcore.WriteSyncer {
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
