package log

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/facily-tech/go-core/telemetry"
)

var _ Logger = (*Zap)(nil)

// Zap wraps a zap.Logger and implements Logger inteface.
type Zap struct {
	logger *zap.Logger
	tracer telemetry.Tracer
}

type ZapConfig struct {
	Version           string
	DisableStackTrace bool
	Tracer            telemetry.Tracer
}

// NewLoggerZap implements Logger using uber zap structured log package.
func NewLoggerZap(config ZapConfig) (*Zap, error) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	loggerConfig.DisableStacktrace = config.DisableStackTrace
	loggerConfig.InitialFields = map[string]interface{}{
		"version": config.Version,
	}

	logger, err := loggerConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Zap{
		logger: logger,
		tracer: config.Tracer,
	}, nil
}

// fieldsToZap convert Fields ([]Field) to []zap.Field.
// and enbed span trace from context
func fieldsToZap(ctx context.Context, tracer telemetry.Tracer, fs []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fs), len(fs)+1)

	for i := range fs {
		zapFields[i] = zap.Any(fs[i].Key, fs[i].Value)
	}

	if tracer != nil {
		if span, ok := tracer.SpanFromContext(ctx); ok {
			spanCtx := span.Context()
			zapFields = append(zapFields, zap.Any(telemetry.TracerKey, spanCtx.ToMap()))
		}
	}

	return zapFields
}

func (z *Zap) Debug(ctx context.Context, msg string, fields ...Field) {
	z.logger.Debug(msg, fieldsToZap(ctx, z.tracer, fields)...)
}

func (z *Zap) Error(ctx context.Context, msg string, fields ...Field) {
	z.logger.Error(msg, fieldsToZap(ctx, z.tracer, fields)...)
}

func (z *Zap) Fatal(ctx context.Context, msg string, fields ...Field) {
	z.logger.Fatal(msg, fieldsToZap(ctx, z.tracer, fields)...)
}

func (z *Zap) Info(ctx context.Context, msg string, fields ...Field) {
	z.logger.Info(msg, fieldsToZap(ctx, z.tracer, fields)...)
}

func (z *Zap) Panic(ctx context.Context, msg string, fields ...Field) {
	z.logger.Panic(msg, fieldsToZap(ctx, z.tracer, fields)...)
}

func (z *Zap) Warn(ctx context.Context, msg string, fields ...Field) {
	z.logger.Warn(msg, fieldsToZap(ctx, z.tracer, fields)...)
}
