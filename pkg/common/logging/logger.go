package logging

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a wrapper around zap logger
type Logger struct {
	*zap.SugaredLogger
}

// Config represents the configuration for the logger
type Config struct {
	Level      string
	ServiceName string
	Environment string
}

// NewLogger creates a new logger with the specified configuration
func NewLogger(cfg *Config) *Logger {
	// Parse log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel // Default to info level
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Add common fields
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).
		With(
			zap.String("service", cfg.ServiceName),
			zap.String("env", cfg.Environment),
		)

	return &Logger{
		SugaredLogger: logger.Sugar(),
	}
}

// WithRequestID adds a request ID field to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		SugaredLogger: l.With("request_id", requestID),
	}
}

// WithUserID adds a user ID field to the logger
func (l *Logger) WithUserID(userID string) *Logger {
	return &Logger{
		SugaredLogger: l.With("user_id", userID),
	}
}

// WithTimestamp logs with the current timestamp
func (l *Logger) WithTimestamp() *Logger {
	return &Logger{
		SugaredLogger: l.With("timestamp", time.Now().Format(time.RFC3339)),
	}
} 