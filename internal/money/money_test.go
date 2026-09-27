package money

import "testing"

func TestYenString(t *testing.T) {
	tests := []struct {
		name string
		in   Yen
		want string
	}{
		{"0円", 0, "¥0"},
		{"3桁", 980, "¥980"},
		{"4桁は区切る", 1980, "¥1,980"},
		{"7桁は2か所で区切る", 1234567, "¥1,234,567"},
		{"マイナス", -1500, "-¥1,500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.String(); got != tt.want {
				t.Errorf("Yen(%d).String() = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
