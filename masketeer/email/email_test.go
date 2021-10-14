package email

import "testing"

func TestMask(t *testing.T) {
	type args struct {
		email string
		opt   *Option
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple test",
			args: args{
				email: "test@example.com",
				opt:   nil,
			},
			want: "tes...@example.com",
		},
		{
			name: "single word on prefix",
			args: args{
				email: "a@example.com",
				opt:   nil,
			},
			want: "a...@example.com",
		},
		{
			name: "no prefix",
			args: args{
				email: "@example.com",
				opt:   nil,
			},
			want: "",
		},
		{
			name: "no domain",
			args: args{
				email: "test",
				opt:   nil,
			},
			want: "",
		},
		{
			name: "no domain with @",
			args: args{
				email: "test@",
				opt:   nil,
			},
			want: "tes...@",
		},
		{
			name: "custom number of visible char on prefix",
			args: args{
				email: "test@example.com",
				opt: &Option{
					NumberOfVisibleCharsOnPrefix: 2,
				},
			},
			want: "te...@example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mask(tt.args.email, tt.args.opt); got != tt.want {
				t.Errorf("Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}
