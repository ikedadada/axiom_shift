package domain

import "testing"

func TestNewEnemy(t *testing.T) {
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
			e := NewEnemy("test", m, tt.growth)
			if e.GrowthRate != tt.growth {
				t.Errorf("GrowthRate: got %v, want %v", e.GrowthRate, tt.growth)
			}
			if m != nil && (e.MatrixState.Rows != m.Rows || e.MatrixState.Cols != m.Cols) {
				t.Errorf("MatrixState: got %v, want %v", e.MatrixState, m)
			}
			if m != nil {
				for i := 0; i < m.Rows; i++ {
					for j := 0; j < m.Cols; j++ {
						if e.MatrixState.Data[i][j] != m.Data[i][j] {
							t.Errorf("MatrixState.Data[%d][%d]: got %v, want %v", i, j, e.MatrixState.Data[i][j], m.Data[i][j])
						}
					}
				}
			}
		})
	}
}

func TestEnemyGrow(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		mdata [][]float64
		rdata [][]float64
	}{
		{"normal", 0.5, [][]float64{{1, 2}, {3, 4}}, [][]float64{{1, 0}, {0, 1}}},
		{"zero matrix and zero rule", 0.5, [][]float64{}, [][]float64{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.mdata)
			rule := NewMatrix(tt.rdata)
			e := NewEnemy("e", m, 1.0)
			e.Grow(tt.input, rule)
			// No panic and matrix state should remain valid
			if e.MatrixState != nil && (e.MatrixState.Rows != len(tt.mdata) || (len(tt.mdata) > 0 && e.MatrixState.Cols != len(tt.mdata[0]))) {
				t.Errorf("MatrixState: got %dx%d, want %dx%d", e.MatrixState.Rows, e.MatrixState.Cols, len(tt.mdata), len(tt.mdata[0]))
			}
		})
	}
}

func TestEnemyGrow_GuardClauses(t *testing.T) {
	tests := []struct {
		name  string
		enemy *Enemy
		rule  *Matrix
		input float64
	}{
		{"nil enemy matrix and nil rule", NewEnemy("e", nil, 1.0), nil, 0.5},
		{"zero matrix and zero rule", NewEnemy("e2", NewMatrix([][]float64{}), 1.0), NewMatrix([][]float64{}), 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Grow panicked: %v", r)
				}
			}()
			before := tt.enemy.MatrixState
			if before != nil {
				beforeCopy := before.Copy()
				tt.enemy.Grow(tt.input, tt.rule)
				if before.Rows == 0 && before.Cols == 0 && !equal(tt.enemy.MatrixState, beforeCopy) {
					t.Error("MatrixState should not change for zero matrix")
				}
			} else {
				tt.enemy.Grow(tt.input, tt.rule)
				if tt.enemy.MatrixState != nil {
					t.Error("MatrixState should remain nil for nil input")
				}
			}
		})
	}
}

func TestEnemyGetMatrix(t *testing.T) {
	m := NewMatrix([][]float64{{1, 2}, {3, 4}})
	e := NewEnemy("test", m, 1.0)
	got := e.GetMatrix()
	if !equal(got, m) {
		t.Error("GetMatrix did not return correct matrix content")
	}
}

func TestEnemyGetMatrix_Nil(t *testing.T) {
	tests := []struct {
		name  string
		enemy *Enemy
		want  *Matrix
	}{
		{"nil MatrixState", NewEnemy("e", nil, 1.0), nil},
		{"non-nil MatrixState", NewEnemy("e", NewMatrix([][]float64{{0, 0}, {0, 0}}), 1.0), NewMatrix([][]float64{{0, 0}, {0, 0}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.enemy.GetMatrix()
			if tt.want == nil && got != nil {
				t.Error("GetMatrix should return nil for nil MatrixState")
			}
			if tt.want != nil && got == nil {
				t.Error("GetMatrix should not return nil for non-nil MatrixState")
			}
		})
	}
}

func TestEnemy_Reset(t *testing.T) {
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
			m := NewMatrix([][]float64{{tt.initVal, 0}, {0, 0}})
			e := NewEnemy("e", m, 1.0)
			e.MatrixState.Data[0][0] = tt.modVal
			e.Reset()
			if e.MatrixState.Data[0][0] != tt.initVal {
				t.Errorf("Enemy.Reset: got %v, want %v", e.MatrixState.Data[0][0], tt.initVal)
			}
		})
	}
}
