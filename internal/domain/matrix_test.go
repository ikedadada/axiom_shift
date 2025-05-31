package domain

import "testing"

func TestNewMatrix(t *testing.T) {
	tests := []struct {
		name    string
		rows    int
		cols    int
		wantNil bool
	}{
		{"normal", 2, 3, false},
		{"zero size", 0, 0, false},
		{"one by one", 1, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.rows, tt.cols)
			if m.Rows != tt.rows || m.Cols != tt.cols {
				t.Errorf("NewMatrix did not set dimensions correctly: got %dx%d", m.Rows, m.Cols)
			}
		})
	}
}

func TestMatrixMultiply(t *testing.T) {
	tests := []struct {
		name      string
		r1, c1    int
		r2, c2    int
		setValues bool
		wantNil   bool
	}{
		{"normal", 2, 2, 2, 2, true, false},
		{"mismatch", 2, 3, 2, 2, false, true},
		{"zero size", 0, 0, 0, 0, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1 := NewMatrix(tt.r1, tt.c1)
			m2 := NewMatrix(tt.r2, tt.c2)
			if tt.setValues && tt.r1 > 0 && tt.c1 > 0 && tt.r2 > 0 && tt.c2 > 0 {
				m1.Data[0][0], m1.Data[0][1] = 1, 2
				m1.Data[1][0], m1.Data[1][1] = 3, 4
				m2.Data[0][0], m2.Data[0][1] = 5, 6
				m2.Data[1][0], m2.Data[1][1] = 7, 8
			}
			res := m1.Multiply(m2)
			if tt.wantNil && res != nil {
				t.Error("Multiply should return nil for mismatched sizes")
			}
			if !tt.wantNil && tt.setValues && res != nil && res.Data[0][0] != 19 {
				t.Error("Matrix Multiply failed")
			}
		})
	}
}

func TestMatrixSubtract(t *testing.T) {
	tests := []struct {
		name     string
		r, c     int
		mismatch bool
	}{
		{"normal", 2, 2, false},
		{"mismatch", 2, 3, true},
		{"zero size", 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m1 := NewMatrix(tt.r, tt.c)
			m2 := NewMatrix(tt.r, tt.c)
			if tt.mismatch {
				m2 = NewMatrix(tt.r, tt.c+1)
			}
			res := m1.Subtract(m2)
			if tt.mismatch && res != nil {
				t.Error("Subtract should return nil for mismatched sizes")
			}
			if !tt.mismatch && m1.Rows > 0 && m1.Cols > 0 && res.Data[0][0] != 0 {
				t.Error("Matrix Subtract failed")
			}
		})
	}
}

func TestMatrixGetScalarValue(t *testing.T) {
	tests := []struct {
		name string
		r, c int
		set  bool
		want float64
	}{
		{"normal", 2, 2, true, 2.5},
		{"empty", 0, 0, false, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.r, tt.c)
			if tt.set && tt.r > 0 && tt.c > 0 {
				m.Data[0][0], m.Data[0][1] = 1, 2
				m.Data[1][0], m.Data[1][1] = 3, 4
			}
			if got := m.GetScalarValue(); got != tt.want {
				t.Errorf("GetScalarValue failed: got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrixNormalize(t *testing.T) {
	tests := []struct {
		name string
		r, c int
		set  bool
	}{
		{"normal", 2, 2, true},
		{"empty", 0, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix(tt.r, tt.c)
			if tt.set && tt.r > 0 && tt.c > 0 {
				m.Data[0][0], m.Data[0][1] = 3, 4
				m.Data[1][0], m.Data[1][1] = 0, 0
			}
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
		{"zero size", NewMatrix(0, 0)},
		{"all zero", func() *Matrix { m := NewMatrix(2, 2); return m }()},
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
	m := NewMatrix(2, 2)
	m.Data[0][0] = 1
	m.Data[1][1] = 2
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

func TestMatrix_sqrt(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"sqrt(4)", 4, 2},
		{"sqrt(0)", 0, 0},
		{"sqrt(-1)", -1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sqrt(tt.input); got != tt.want {
				t.Errorf("sqrt(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMatrix_equal(t *testing.T) {
	tests := []struct {
		name string
		a, b *Matrix
		want bool
	}{
		{"both nil", nil, nil, true},
		{"a nil, b not nil", nil, NewMatrix(2, 2), false},
		{"a not nil, b nil", NewMatrix(2, 2), nil, false},
		{"different size", NewMatrix(2, 2), NewMatrix(3, 2), false},
		{"same size, different data", func() *Matrix { m := NewMatrix(2, 2); m.Data[0][0] = 1; return m }(), NewMatrix(2, 2), false},
		{"same size, same data", func() *Matrix { m := NewMatrix(2, 2); m.Data[0][0] = 1; return m }(), func() *Matrix { m := NewMatrix(2, 2); m.Data[0][0] = 1; return m }(), true},
		{"zero size", NewMatrix(0, 0), NewMatrix(0, 0), true},
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
