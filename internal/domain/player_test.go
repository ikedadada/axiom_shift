package domain

import "testing"

func TestNewPlayer(t *testing.T) {
	tests := []struct {
		name   string
		data   [][]float64
		growth float64
	}{
		{"normal", [][]float64{{1, 2}, {3, 4}}, 0.5},
		{"zero_matrix", [][]float64{}, 1},
		{"negative_growth", [][]float64{{1, 2}, {3, 4}}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.data)
			p := NewPlayer(m, tt.growth)
			if p.GrowthRate != tt.growth {
				t.Errorf("GrowthRate: got %v, want %v", p.GrowthRate, tt.growth)
			}
			if m != nil && (p.MatrixState.Rows != m.Rows || p.MatrixState.Cols != m.Cols) {
				t.Errorf("MatrixState: got %v, want %v", p.MatrixState, m)
			}
			if m != nil {
				for i := 0; i < m.Rows; i++ {
					for j := 0; j < m.Cols; j++ {
						if p.MatrixState.Data[i][j] != m.Data[i][j] {
							t.Errorf("MatrixState.Data[%d][%d]: got %v, want %v", i, j, p.MatrixState.Data[i][j], m.Data[i][j])
						}
					}
				}
			}
		})
	}
}

func TestPlayerUpdateMatrix(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		data  [][]float64
	}{
		{"normal", 0.5, [][]float64{{1, 2}, {3, 4}}},
		{"zero size", 0.5, [][]float64{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.data)
			p := NewPlayer(m, 1.0)
			p.UpdateMatrix(tt.input)
			// No panic and matrix state should remain valid
			if p.MatrixState != nil && (p.MatrixState.Rows != len(tt.data) || (len(tt.data) > 0 && p.MatrixState.Cols != len(tt.data[0]))) {
				t.Errorf("MatrixState: got %dx%d, want %dx%d", p.MatrixState.Rows, p.MatrixState.Cols, len(tt.data), len(tt.data[0]))
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
	m := NewMatrix([][]float64{{1, 2}, {3, 4}})
	p := NewPlayer(m, 1.0)
	got := p.GetMatrix()
	if !equal(got, m) {
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
		{"non-nil MatrixState", NewPlayer(NewMatrix([][]float64{{0, 0}, {0, 0}}), 1.0), NewMatrix([][]float64{{0, 0}, {0, 0}})},
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
	data := [][]float64{{1, 2}, {3, 4}}
	m := NewMatrix(data)
	p := NewPlayer(m, 1.0)
	p.MatrixState.Data[0][0] = 99
	p.Reset()
	if p.MatrixState.Data[0][0] != 1 {
		t.Error("Reset did not restore initial state")
	}
}

func TestPlayer_NewPlayer(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]float64
		wantRows int
		wantCols int
	}{
		{"normal", [][]float64{{1, 2}, {3, 4}}, 2, 2},
		{"empty", [][]float64{}, 0, 0},
		{"empty inner", [][]float64{{}}, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.data)
			p := NewPlayer(m, 1.0)
			if p.MatrixState == nil {
				if tt.wantRows != 0 || tt.wantCols != 0 {
					t.Errorf("MatrixState is nil, want %dx%d", tt.wantRows, tt.wantCols)
				}
				return
			}
			if p.MatrixState.Rows != tt.wantRows || p.MatrixState.Cols != tt.wantCols {
				t.Errorf("NewPlayer: got %dx%d, want %dx%d", p.MatrixState.Rows, p.MatrixState.Cols, tt.wantRows, tt.wantCols)
			}
		})
	}
}
