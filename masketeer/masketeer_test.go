package masketeer

import "testing"

func TestMasketeer_Email(t *testing.T) {
	type args struct {
		eml string
	}
	tests := []struct {
		name string
		opt  *Option
		args args
		want string
	}{
		{
			name: "email mask",
			opt:  nil,
			args: args{
				eml: "test@example.com",
			},
			want: "tes...@example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.opt)
			if got := m.Email(tt.args.eml); got != tt.want {
				t.Errorf("Masketeer.Email() = %v, want %v", got, tt.want)
			}
		})
	}
}
