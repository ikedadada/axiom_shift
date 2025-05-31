package domain

import "testing"

func TestNewEnemy(t *testing.T) {
	tests := []struct {
		name       string
		matrix     *Matrix
		growthRate float64
		wantName   string
	}{
		{"normal", NewMatrix(2, 2), 0.5, "test"},
		{"zero matrix", NewMatrix(0, 0), 1.0, "zero"},
		{"nil matrix", nil, 1.0, "nil"},
		{"negative growth", NewMatrix(2, 2), -1.0, "neg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEnemy(tt.wantName, tt.matrix, tt.growthRate)
			if e.Name != tt.wantName || e.Matrix != tt.matrix || e.GrowthRate != tt.growthRate {
				t.Errorf("NewEnemy fields not set correctly: got %+v", e)
			}
		})
	}
}

func TestEnemyGrow(t *testing.T) {
	tests := []struct {
		name   string
		input  float64
		rows   int
		cols   int
		growth float64
	}{
		{"input 0.0", 0.0, 2, 2, 1.0},
		{"input 1.0", 1.0, 2, 2, 1.0},
		{"input 0.5", 0.5, 2, 2, 2.0},
		{"input negative", -0.5, 2, 2, 1.0},
		{"input >1", 2.0, 2, 2, 1.0},
		{"zero matrix", 0.0, 0, 0, 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.rows, tt.cols)
			e := NewEnemy("test", m, tt.growth)
			rule := NewMatrix(tt.rows, tt.cols)
			// set some values for rule
			if tt.rows > 0 && tt.cols > 0 {
				rule.Data[0][0] = 1.0
				rule.Data[tt.rows-1][tt.cols-1] = 2.0
			}
			e.Grow(tt.input, rule)
			if m != nil && m.Rows > 0 && m.Cols > 0 {
				// At least one element should be incremented
				found := false
				for i := 0; i < m.Rows; i++ {
					for j := 0; j < m.Cols; j++ {
						if m.Data[i][j] != 0 {
							found = true
						}
					}
				}
				if !found {
					t.Errorf("Grow did not update any element: got %v", m.Data)
				}
			}
		})
	}
}

func TestEnemyGetMatrix(t *testing.T) {
	m := NewMatrix(2, 2)
	e := NewEnemy("test", m, 1.0)
	if e.GetMatrix() != m {
		t.Error("GetMatrix did not return correct matrix")
	}
}
