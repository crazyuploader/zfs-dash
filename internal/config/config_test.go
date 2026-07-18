package config

import (
	"testing"
	"time"
)

func TestParseFlexDuration(t *testing.T) {
	def := 300 * time.Second
	tests := []struct {
		name string
		in   any
		want time.Duration
	}{
		{"nil", nil, def},
		{"int seconds", 120, 120 * time.Second},
		{"int64 seconds", int64(60), 60 * time.Second},
		{"float seconds", 90.0, 90 * time.Second},
		{"numeric string", "120", 120 * time.Second},
		{"duration string", "5m", 5 * time.Minute},
		{"compound duration", "1h30m", 90 * time.Minute},
		{"garbage", "not-a-duration", def},
		{"zero", 0, def},
		{"negative", -5, def},
		{"negative duration string", "-5m", def},
		{"unsupported type", []string{"5m"}, def},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFlexDuration(tt.in, def); got != tt.want {
				t.Errorf("parseFlexDuration(%v) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
