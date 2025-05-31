package domain

import "testing"

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name   string
		matrix *Matrix
		growth float64
	}{
		{"normal", NewMatrix(2, 2), 0.5},
		{"zero_matrix", NewMatrix(0, 0), 1},
		{"negative_growth", NewMatrix(2, 2), -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.matrix
			m.Data = [][]float64{{1, 2}, {3, 4}}
			p := NewPlayer(m, tt.growth)
			if p.GrowthRate != tt.growth {
				t.Errorf("GrowthRate: got %v, want %v", p.GrowthRate, tt.growth)
			}
			if p.MatrixState.Rows != m.Rows || p.MatrixState.Cols != m.Cols {
				t.Errorf("MatrixState: got %v, want %v", p.MatrixState, m)
			}
			for i := 0; i < m.Rows; i++ {
				for j := 0; j < m.Cols; j++ {
					if p.MatrixState.Data[i][j] != m.Data[i][j] {
						t.Errorf("MatrixState.Data[%d][%d]: got %v, want %v", i, j, p.MatrixState.Data[i][j], m.Data[i][j])
					}
				}
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
			if p.MatrixState == nil || p.MatrixState.Rows == 0 || p.MatrixState.Cols == 0 {
				return // nothing to check
			}
			// Check that the target index is incremented by growth
			if p.MatrixState.Data[tt.wantIndex[0]][tt.wantIndex[1]] < tt.growth {
				t.Errorf("Target index not incremented as expected: got %v", p.MatrixState.Data)
			}
		})
	}
}

func TestPlayerUpdateMatrix_GuardClauses(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("UpdateMatrix panicked with nil MatrixState: %v", r)
		}
	}()
	p := NewPlayer(nil, 1.0)
	p.UpdateMatrix(0.5)
	if p.MatrixState != nil {
		t.Error("MatrixState should remain nil when initialized with nil")
	}
}

func TestPlayerGetMatrix(t *testing.T) {
	m := NewMatrix(2, 2)
	m.Data[0][0], m.Data[0][1] = 1, 2
	m.Data[1][0], m.Data[1][1] = 3, 4
	p := NewPlayer(m, 1.0)
	got := p.GetMatrix()
	if !matricesEqual(got, m) {
		t.Error("GetMatrix did not return correct matrix content")
	}
}

func TestPlayerGetMatrix_Nil(t *testing.T) {
	tests := []struct {
		name   string
		player *Player
		want   *Matrix
	}{
		{"nil MatrixState", NewPlayer(nil, 1.0), nil},
		{"non-nil MatrixState", NewPlayer(NewMatrix(2, 2), 1.0), NewMatrix(2, 2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.player.GetMatrix()
			if tt.want == nil && got != nil {
				t.Error("GetMatrix should return nil for nil MatrixState")
			}
			if tt.want != nil && got == nil {
				t.Error("GetMatrix should not return nil for non-nil MatrixState")
			}
		})
	}
}

func TestPlayer_Reset(t *testing.T) {
	tests := []struct {
		name    string
		initVal float64
		modVal  float64
		want    float64
	}{
		{"reset to 1", 1, 99, 1},
		{"reset to 0", 0, 42, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(2, 2)
			m.Data[0][0] = tt.initVal
			p := NewPlayer(m, 1.0)
			p.MatrixState.Data[0][0] = tt.modVal
			p.Reset()
			if p.MatrixState.Data[0][0] != tt.initVal {
				t.Errorf("Player.Reset: got %v, want %v", p.MatrixState.Data[0][0], tt.initVal)
			}
		})
	}
}
