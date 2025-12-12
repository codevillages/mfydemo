package logger

import (
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New builds a zap.Logger from config.
func New(cfg Config) (*zap.Logger, error) {
	encoding := cfg.Encoding
	if encoding == "" {
		encoding = "json"
	}
	level := zap.InfoLevel
	if cfg.Level != "" {
		if parsed, err := zapcore.ParseLevel(strings.ToLower(cfg.Level)); err == nil {
			level = parsed
		}
	}

	zapCfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: false,
		Encoding:    encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorOutputPaths,
	}

	if len(zapCfg.OutputPaths) == 0 {
		zapCfg.OutputPaths = []string{"stdout"}
	}
	if len(zapCfg.ErrorOutputPaths) == 0 {
		zapCfg.ErrorOutputPaths = []string{"stderr"}
	}

	return zapCfg.Build()
}

// SetupGoZero bridges go-zero logx to zap.
func SetupGoZero(logger *zap.Logger) {
	if logger == nil {
		return
	}
	logx.DisableStat()
	logx.SetWriter(&zapWriter{log: logger})
}

type zapWriter struct {
	log *zap.Logger
}

func (w *zapWriter) Alert(v any) {
	w.log.Error(asMessage(v))
}

func (w *zapWriter) Close() error {
	return nil
}

func (w *zapWriter) Debug(v any, fields ...logx.LogField) {
	w.log.Debug(asMessage(v), zapFields(fields)...)
}

func (w *zapWriter) Info(v any, fields ...logx.LogField) {
	w.log.Info(asMessage(v), zapFields(fields)...)
}

func (w *zapWriter) Error(v any, fields ...logx.LogField) {
	w.log.Error(asMessage(v), zapFields(fields)...)
}

func (w *zapWriter) Severe(v any) {
	w.log.Error(asMessage(v))
}

func (w *zapWriter) Slow(v any, fields ...logx.LogField) {
	w.log.Warn(asMessage(v), zapFields(fields)...)
}

func (w *zapWriter) Stack(v any) {
	w.log.Error(asMessage(v))
}

func (w *zapWriter) Stat(v any, fields ...logx.LogField) {
	w.log.Info(asMessage(v), zapFields(fields)...)
}

func asMessage(v interface{}) string {
	if v == nil {
		return ""
	}
	if msg, ok := v.(string); ok {
		return msg
	}
	return ""
}

func zapFields(fields []logx.LogField) []zap.Field {
	zFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zFields = append(zFields, zap.Any(f.Key, f.Value))
	}
	return zFields
}
