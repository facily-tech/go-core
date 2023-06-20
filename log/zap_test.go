package log

import (
	"context"
	"reflect"
	"testing"

	"github.com/facily-tech/go-core/telemetry"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func Test_fieldsToZap(t *testing.T) {
	type args struct {
		ctx    context.Context
		tracer telemetry.Tracer
		fs     []Field
	}
	tests := []struct {
		name string
		args args
		want []zap.Field
	}{
		{
			name: "simple",
			args: args{
				ctx:    context.Background(),
				tracer: nil,
				fs: []Field{
					{
						Key:   "foo",
						Value: "boo",
					},
				},
			},
			want: func() []zap.Field {
				zapFields := make([]zap.Field, 1)
				zapFields[0] = zap.Any("foo", "boo")

				return zapFields
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fieldsToZap(tt.args.ctx, tt.args.tracer, tt.args.fs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fieldsToZap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZap_With(t *testing.T) {
	buff := &zaptest.Buffer{}
	type fields struct {
		logger *zap.Logger
		tracer telemetry.Tracer
	}
	type args struct {
		ctx    context.Context
		fields []Field
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Logger
	}{
		{
			name: "included field must be present?",
			fields: fields{
				logger: zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewDevelopmentConfig().EncoderConfig), buff, zap.DebugLevel)),
				tracer: nil,
			},
			args: args{
				ctx:    context.Background(),
				fields: []Field{Any("versionApp", "4.39.2")},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Zap{
				logger: tt.fields.logger,
				tracer: tt.fields.tracer,
			}
			z.With(tt.args.ctx, tt.args.fields...).Info(tt.args.ctx, "ola")
			assert.Contains(t, buff.String(), "versionApp")
			assert.Contains(t, buff.String(), "4.39.2")

		})
	}
}
