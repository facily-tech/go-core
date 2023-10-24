package telemetry

import (
	"context"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"testing"
)

func TestDataDog_StartSpanForCommand(t *testing.T) {
	type args struct {
		action string
		params []CommandResultParam
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should create a span for a command",
			args: args{
				action: "test",
				params: []CommandResultParam{
					{
						Key:   "name",
						Value: "John",
					},
					{
						Key:   "lastname",
						Value: "Doe",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt := mocktracer.Start()
			defer mt.Stop()

			ctx := context.TODO()

			d := DataDog{}
			d.StartSpanForCommand(
				ctx,
				tt.args.action,
				func(ctx context.Context) []CommandResultParam {
					return tt.args.params
				},
			)

			spans := mt.FinishedSpans()
			totalSpans := len(spans)

			if totalSpans != 1 {
				t.Errorf("expected 1 span, got %d", totalSpans)
			}
		})
	}
}
