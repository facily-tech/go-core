package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMake(t *testing.T) {
	type args struct {
		basePath        string
		requestPath     string
		queryParameters interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				basePath:    "https://example.com/",
				requestPath: "/resource",
				queryParameters: struct {
					Name string `url:"name"`
				}{Name: "facily"},
			},
			want: "https://example.com/resource?name=facily",
		},
		{
			name: "success, if empty queryParameters",
			args: args{
				basePath:        "https://example.com/",
				requestPath:     "/resource",
				queryParameters: nil,
			},
			want: "https://example.com/resource",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Make(tt.args.basePath, tt.args.requestPath, tt.args.queryParameters)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
