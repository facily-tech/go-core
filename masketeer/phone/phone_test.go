package phone

import "testing"

func TestMask(t *testing.T) {
	type args struct {
		phone string
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
				phone: "+12 123 123456789",
				opt:   nil,
			},
			want: "...6789",
		},
		{
			name: "not a phone number",
			args: args{
				phone: "test@example.com",
				opt:   nil,
			},
			want: "",
		},
		{
			name: "a highter number of visible chars than there is on phone",
			args: args{
				phone: "1234",
				opt: &Option{
					NumberOfVisibleCharsOnSufix: 5,
				},
			},
			want: "...1234",
		},
		{
			name: "accepting more chars",
			args: args{
				phone: "+12 124 123.46-98",
				opt: &Option{
					NumberOfVisibleCharsOnSufix: 7,
					UseAsVisibleChars:           DefaultUseAsVisibleChars + ".-",
				},
			},
			want: "...3.46-98",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mask(tt.args.phone, tt.args.opt); got != tt.want {
				t.Errorf("Mask() = \"%v\", want \"%v\"", got, tt.want)
			}
		})
	}
}
