package domain

import "testing"

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name       string
		matrix     *Matrix
		growthRate float64
		wantNil    bool
	}{
		{"normal", NewMatrix(2, 2), 0.5, false},
		{"zero matrix", NewMatrix(0, 0), 1.0, false},
		{"nil matrix", nil, 1.0, false},
		{"negative growth", NewMatrix(2, 2), -1.0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(tt.matrix, tt.growthRate)
			if p.MatrixState != tt.matrix || p.GrowthRate != tt.growthRate {
				t.Errorf("NewPlayer fields not set correctly: got %+v", p)
			}
		})
	}
}

func TestPlayerUpdateMatrix(t *testing.T) {
	tests := []struct {
		name       string
		input      float64
		rows, cols int
		growth     float64
		wantIndex  [2]int
	}{
		{"input 0.0", 0.0, 2, 2, 1.0, [2]int{0, 0}},
		{"input 1.0", 1.0, 2, 2, 1.0, [2]int{1, 1}},
		{"input 0.5", 0.5, 2, 2, 2.0, [2]int{1, 0}},
		{"input negative", -0.5, 2, 2, 1.0, [2]int{0, 0}},
		{"input >1", 2.0, 2, 2, 1.0, [2]int{1, 1}},
		{"zero matrix", 0.0, 0, 0, 1.0, [2]int{0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.rows, tt.cols)
			p := NewPlayer(m, tt.growth)
			p.UpdateMatrix(tt.input)
			if m.Rows == 0 || m.Cols == 0 {
				return // nothing to check
			}
			// Check that the target index is incremented by growth
			if m.Data[tt.wantIndex[0]][tt.wantIndex[1]] < tt.growth {
				t.Errorf("Target index not incremented as expected: got %v", m.Data)
			}
		})
	}
}

func TestPlayerGetMatrix(t *testing.T) {
	m := NewMatrix(2, 2)
	p := NewPlayer(m, 1.0)
	if p.GetMatrix() != m {
		t.Error("GetMatrix did not return correct matrix")
	}
}
