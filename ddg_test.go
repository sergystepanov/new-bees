package main

import "testing"

func Test_token(t *testing.T) {
	tests := []struct {
		name string
		len  int
		want int
	}{
		{name: "a", len: 16, want: 16},
		{name: "b", len: 20, want: 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := token(tt.len)
			t.Logf("token: %v", got)
			if len(got) != tt.want {
				t.Errorf("token() = %v, want %v", got, tt.want)
			}
		})
	}
}
