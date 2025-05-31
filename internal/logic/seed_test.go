package logic

import (
	"testing"
)

func TestSeedDeterminism(t *testing.T) {
	tests := []struct {
		name  string
		seed1 int64
		seed2 int64
		size  int
		equal bool
	}{
		{"same seed same result", 42, 42, 2, true},
		{"different seed different result", 42, 43, 2, false},
		{"zero size", 1, 1, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule1 := NewRuleMatrix(tt.seed1, tt.size)
			rule2 := NewRuleMatrix(tt.seed2, tt.size)
			if tt.size == 0 {
				return // always equal for zero size
			}
			equal := true
			for i := 0; i < tt.size; i++ {
				for j := 0; j < tt.size; j++ {
					if rule1.matrix[i][j] != rule2.matrix[i][j] {
						equal = false
					}
				}
			}
			if equal != tt.equal {
				t.Errorf("Seed determinism failed: got %v, want %v", equal, tt.equal)
			}
		})
	}
}
