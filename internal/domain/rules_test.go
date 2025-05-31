package domain

import (
	"testing"
)

func TestNewRuleMatrix(t *testing.T) {
	tests := []struct {
		name  string
		seed  int64
		size  int
		wantN int
	}{
		{"normal 2x2", 42, 2, 2},
		{"size 0", 42, 0, 0},
		{"size 1", 1, 1, 1},
		{"large", 99, 10, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewRuleMatrix(tt.seed, tt.size)
			if rule == nil {
				t.Fatal("NewRuleMatrix returned nil")
			}
			if len(rule.Matrix.Data) != tt.wantN {
				t.Errorf("matrix row size: got %d, want %d", len(rule.Matrix.Data), tt.wantN)
			}
			if tt.wantN > 0 && len(rule.Matrix.Data[0]) != tt.wantN {
				t.Errorf("matrix col size: got %d, want %d", len(rule.Matrix.Data[0]), tt.wantN)
			}
		})
	}
}
