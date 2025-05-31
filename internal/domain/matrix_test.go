package domain

import "testing"

func TestNewMatrix(t *testing.T) {
	tests := []struct {
		name     string
		data     [][]float64
		wantRows int
		wantCols int
	}{
		{"normal", [][]float64{{1, 2}, {3, 4}}, 2, 2},
		{"empty outer", [][]float64{}, 0, 0},
		{"empty inner", [][]float64{{}}, 1, 0},
		{"jagged (should use first row)", [][]float64{{1, 2}, {3}}, 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.data)
			if m.Rows != tt.wantRows || m.Cols != tt.wantCols {
				t.Errorf("NewMatrix: got %dx%d, want %dx%d", m.Rows, m.Cols, tt.wantRows, tt.wantCols)
			}
			if m.Rows > 0 && m.Cols > 0 && len(m.Data) != m.Rows {
				t.Errorf("NewMatrix: Data row count mismatch")
			}
			if m.Rows > 0 && m.Cols > 0 && len(m.Data[0]) != m.Cols {
				t.Errorf("NewMatrix: Data col count mismatch")
			}
		})
	}
}

func TestMatrixMultiply(t *testing.T) {
	tests := []struct {
		name    string
		m1, m2  [][]float64
		wantNil bool
		wantVal float64 // expected value at [0][0] if not nil
	}{
		{"normal", [][]float64{{1, 2}, {3, 4}}, [][]float64{{5, 6}, {7, 8}}, false, 19},
		{"mismatch", [][]float64{{1, 2, 3}, {4, 5, 6}}, [][]float64{{1, 2}, {3, 4}}, true, 0},
		{"zero size", [][]float64{}, [][]float64{}, false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1 := NewMatrix(tt.m1)
			m2 := NewMatrix(tt.m2)
			res := m1.Multiply(m2)
			if tt.wantNil && res != nil {
				t.Error("Multiply should return nil for mismatched sizes")
			}
			if !tt.wantNil && len(tt.m1) > 0 && len(tt.m2) > 0 && res != nil && res.Data[0][0] != tt.wantVal {
				t.Errorf("Matrix Multiply failed: got %v, want %v", res.Data[0][0], tt.wantVal)
			}
		})
	}
}

func TestMatrixSubtract(t *testing.T) {
	tests := []struct {
		name     string
		m1, m2   [][]float64
		mismatch bool
	}{
		{"normal", [][]float64{{1, 2}, {3, 4}}, [][]float64{{1, 2}, {3, 4}}, false},
		{"mismatch", [][]float64{{1, 2}, {3, 4}}, [][]float64{{1, 2, 3}, {4, 5, 6}}, true},
		{"zero size", [][]float64{}, [][]float64{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1 := NewMatrix(tt.m1)
			m2 := NewMatrix(tt.m2)
			if tt.mismatch {
				m2 = NewMatrix(tt.m2)
			}
			res := m1.Subtract(m2)
			if tt.mismatch && res != nil {
				t.Error("Subtract should return nil for mismatched sizes")
			}
			if !tt.mismatch && len(tt.m1) > 0 && len(tt.m2) > 0 && res.Data[0][0] != 0 {
				t.Error("Matrix Subtract failed")
			}
		})
	}
}

func TestMatrixGetScalarValue(t *testing.T) {
	tests := []struct {
		name string
		data [][]float64
		want float64
	}{
		{"normal", [][]float64{{1, 2}, {3, 4}}, 2.5},
		{"empty", [][]float64{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.data)
			if got := m.GetScalarValue(); got != tt.want {
				t.Errorf("GetScalarValue failed: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrixNormalize(t *testing.T) {
	tests := []struct {
		name string
		data [][]float64
		set  bool
	}{
		{"normal", [][]float64{{3, 4}, {0, 0}}, true},
		{"empty", [][]float64{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.data)
			m.Normalize()
			norm := 0.0
			for i := 0; i < m.Rows; i++ {
				for j := 0; j < m.Cols; j++ {
					norm += m.Data[i][j] * m.Data[i][j]
				}
			}
			if m.Rows > 0 && m.Cols > 0 && abs(norm-1.0) > 1e-6 {
				t.Errorf("Normalize failed: norm=%v", norm)
			}
		})
	}
}

func TestMatrixNormalize_GuardClauses(t *testing.T) {
	tests := []struct {
		name   string
		matrix *Matrix
	}{
		{"nil matrix", nil},
		{"zero size", NewMatrix([][]float64{})},
		{"all zero", NewMatrix([][]float64{{0, 0}, {0, 0}})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Normalize() panicked: %v", r)
				}
			}()
			if tt.matrix != nil {
				tt.matrix.Normalize()
			} else {
				var m *Matrix
				m.Normalize()
			}
		})
	}
}

func TestMatrix_Copy(t *testing.T) {
	m := NewMatrix([][]float64{{1, 0}, {0, 2}})
	copy := m.Copy()
	if copy.Rows != 2 || copy.Cols != 2 {
		t.Error("Copy: dimension mismatch")
	}
	if copy.Data[0][0] != 1 || copy.Data[1][1] != 2 {
		t.Error("Copy: data mismatch")
	}
	copy.Data[0][0] = 99
	if m.Data[0][0] == 99 {
		t.Error("Copy: not deep copy")
	}
}

func TestMatrixCopy_GuardClauses(t *testing.T) {
	var m *Matrix
	if m.Copy() != nil {
		t.Error("Copy: nil matrix should return nil")
	}
}

func TestMatrix_equal(t *testing.T) {
	tests := []struct {
		name string
		a, b *Matrix
		want bool
	}{
		{"both nil", nil, nil, true},
		{"a nil, b not nil", nil, NewMatrix([][]float64{{0, 0}, {0, 0}}), false},
		{"a not nil, b nil", NewMatrix([][]float64{{0, 0}, {0, 0}}), nil, false},
		{"different size", NewMatrix([][]float64{{0, 0}, {0, 0}}), NewMatrix([][]float64{{0, 0}, {0, 0}, {0, 0}}), false},
		{"same size, different data", NewMatrix([][]float64{{1, 0}, {0, 0}}), NewMatrix([][]float64{{0, 0}, {0, 0}}), false},
		{"same size, same data", NewMatrix([][]float64{{1, 0}, {0, 0}}), NewMatrix([][]float64{{1, 0}, {0, 0}}), true},
		{"zero size", NewMatrix([][]float64{}), NewMatrix([][]float64{}), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := equal(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
