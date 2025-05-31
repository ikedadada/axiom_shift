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

func TestSeedManager_Basic(t *testing.T) {
	sm := NewSeedManager()
	seed := sm.GetSeed()
	sm.SetSeed(seed + 1)
	if sm.GetSeed() != seed+1 {
		t.Error("SetSeed or GetSeed failed")
	}
}

func TestSeedManagerWithFixedValue(t *testing.T) {
	sm := NewSeedManagerWithFixedValue(123)
	if sm.GetSeed() != 123 {
		t.Error("NewSeedManagerWithFixedValue failed")
	}
}

func TestSeedManager_Randoms(t *testing.T) {
	sm := NewSeedManagerWithFixedValue(42)
	f := sm.RandomFloat64()
	if f < 0 || f > 1 {
		t.Error("RandomFloat64 out of range")
	}
	_ = sm.RandomInt(0, 10)
}
