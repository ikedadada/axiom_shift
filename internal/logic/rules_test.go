package logic

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
			if len(rule.matrix) != tt.wantN {
				t.Errorf("matrix row size: got %d, want %d", len(rule.matrix), tt.wantN)
			}
			if tt.wantN > 0 && len(rule.matrix[0]) != tt.wantN {
				t.Errorf("matrix col size: got %d, want %d", len(rule.matrix[0]), tt.wantN)
			}
		})
	}
}

func TestRuleMatrixGetMatrix(t *testing.T) {
	tests := []struct {
		name string
		seed int64
		size int
	}{
		{"normal", 42, 2},
		{"empty", 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewRuleMatrix(tt.seed, tt.size)
			mat := rule.GetMatrix()
			if len(mat) != tt.size {
				t.Errorf("GetMatrix row size: got %d, want %d", len(mat), tt.size)
			}
			if tt.size > 0 && len(mat[0]) != tt.size {
				t.Errorf("GetMatrix col size: got %d, want %d", len(mat[0]), tt.size)
			}
		})
	}
}
