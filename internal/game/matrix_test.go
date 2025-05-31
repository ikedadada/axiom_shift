package game

import "testing"

func TestNewMatrix(t *testing.T) {
	m := NewMatrix(2, 3)
	if m.rows != 2 || m.cols != 3 {
		t.Error("NewMatrix did not set dimensions correctly")
	}
}

func TestMatrixMultiply(t *testing.T) {
	m1 := NewMatrix(2, 2)
	m2 := NewMatrix(2, 2)
	m1.data[0][0], m1.data[0][1] = 1, 2
	m1.data[1][0], m1.data[1][1] = 3, 4
	m2.data[0][0], m2.data[0][1] = 5, 6
	m2.data[1][0], m2.data[1][1] = 7, 8
	res := m1.Multiply(m2)
	if res.data[0][0] != 19 || res.data[0][1] != 22 || res.data[1][0] != 43 || res.data[1][1] != 50 {
		t.Error("Matrix Multiply failed")
	}
}

func TestMatrixMultiply_ErrorCase(t *testing.T) {
	m1 := NewMatrix(2, 3)
	m2 := NewMatrix(2, 2)
	if m1.Multiply(m2) != nil {
		t.Error("Multiply should return nil for mismatched sizes")
	}
}

func TestMatrixSubtract(t *testing.T) {
	m1 := NewMatrix(2, 2)
	m2 := NewMatrix(2, 2)
	m1.data[0][0], m2.data[0][0] = 5, 2
	m1.data[1][1], m2.data[1][1] = 7, 4
	res := m1.Subtract(m2)
	if res.data[0][0] != 3 || res.data[1][1] != 3 {
		t.Error("Matrix Subtract failed")
	}
}

func TestMatrixSubtract_ErrorCase(t *testing.T) {
	m1 := NewMatrix(2, 2)
	m2 := NewMatrix(2, 3)
	if m1.Subtract(m2) != nil {
		t.Error("Subtract should return nil for mismatched sizes")
	}
}

func TestMatrixGetScalarValue(t *testing.T) {
	m := NewMatrix(2, 2)
	m.data[0][0], m.data[0][1] = 1, 2
	m.data[1][0], m.data[1][1] = 3, 4
	if m.GetScalarValue() != 2.5 {
		t.Errorf("GetScalarValue failed: got %v", m.GetScalarValue())
	}
}

func TestMatrixGetScalarValue_Empty(t *testing.T) {
	m := NewMatrix(0, 0)
	if m.GetScalarValue() != 0 {
		t.Error("GetScalarValue should return 0 for empty matrix")
	}
}

func TestMatrixNormalize(t *testing.T) {
	m := NewMatrix(2, 2)
	m.data[0][0], m.data[0][1] = 3, 4
	m.data[1][0], m.data[1][1] = 0, 0
	m.Normalize()
	norm := 0.0
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			norm += m.data[i][j] * m.data[i][j]
		}
	}
	if abs(norm-1.0) > 1e-6 {
		t.Errorf("Normalize failed: norm=%v", norm)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
