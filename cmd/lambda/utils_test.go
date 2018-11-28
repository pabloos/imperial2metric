package main

import "testing"

func Test_isAZip(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "zip filename testcase",
			args: args{
				filename: "filename.zip",
			},
			want: true,
		},
		{
			name: "other filename testcase",
			args: args{
				filename: "filename.rar",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAZip(tt.args.filename); got != tt.want {
				t.Errorf("isAZip() = %v, want %v", got, tt.want)
			}
		})
	}
}
