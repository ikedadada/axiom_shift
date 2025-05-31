package domain

import "testing"

func TestNewEnemy(t *testing.T) {
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
			e := NewEnemy("test", m, tt.growth)
			if e.GrowthRate != tt.growth {
				t.Errorf("GrowthRate: got %v, want %v", e.GrowthRate, tt.growth)
			}
			if e.MatrixState.Rows != m.Rows || e.MatrixState.Cols != m.Cols {
				t.Errorf("MatrixState: got %v, want %v", e.MatrixState, m)
			}
			for i := 0; i < m.Rows; i++ {
				for j := 0; j < m.Cols; j++ {
					if e.MatrixState.Data[i][j] != m.Data[i][j] {
						t.Errorf("MatrixState.Data[%d][%d]: got %v, want %v", i, j, e.MatrixState.Data[i][j], m.Data[i][j])
					}
				}
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
			if e.MatrixState != nil && e.MatrixState.Rows > 0 && e.MatrixState.Cols > 0 {
				// At least one element should be incremented
				found := false
				for i := 0; i < e.MatrixState.Rows; i++ {
					for j := 0; j < e.MatrixState.Cols; j++ {
						if e.MatrixState.Data[i][j] != 0 {
							found = true
						}
					}
				}
				if !found {
					t.Errorf("Grow did not update any element: got %v", e.MatrixState.Data)
				}
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
		{"zero matrix and zero rule", NewEnemy("e2", NewMatrix(0, 0), 1.0), NewMatrix(0, 0), 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic or update anything
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Grow panicked: %v", r)
				}
			}()
			before := tt.enemy.MatrixState
			if before != nil {
				beforeCopy := before.Copy()
				tt.enemy.Grow(tt.input, tt.rule)
				// MatrixState should not change for zero matrix
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
	m := NewMatrix(2, 2)
	m.Data[0][0], m.Data[0][1] = 1, 2
	m.Data[1][0], m.Data[1][1] = 3, 4
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
		{"non-nil MatrixState", NewEnemy("e", NewMatrix(2, 2), 1.0), NewMatrix(2, 2)},
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

// テーブル駆動型TestEnemy_Resetのみ残す
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
			m := NewMatrix(2, 2)
			m.Data[0][0] = tt.initVal
			e := NewEnemy("e", m, 1.0)
			e.MatrixState.Data[0][0] = tt.modVal
			e.Reset()
			if e.MatrixState.Data[0][0] != tt.initVal {
				t.Errorf("Enemy.Reset: got %v, want %v", e.MatrixState.Data[0][0], tt.initVal)
			}
		})
	}
}
