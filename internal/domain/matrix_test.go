package domain

import "testing"

func TestNewMatrix(t *testing.T) {
	m := NewMatrix(2, 3)
	if m.Rows != 2 || m.Cols != 3 {
		t.Error("NewMatrix did not set dimensions correctly")
	}
}

func TestMatrixMultiply(t *testing.T) {
	m1 := NewMatrix(2, 2)
	m2 := NewMatrix(2, 2)
	m1.Data[0][0], m1.Data[0][1] = 1, 2
	m1.Data[1][0], m1.Data[1][1] = 3, 4
	m2.Data[0][0], m2.Data[0][1] = 5, 6
	m2.Data[1][0], m2.Data[1][1] = 7, 8
	res := m1.Multiply(m2)
	if res.Data[0][0] != 19 || res.Data[0][1] != 22 || res.Data[1][0] != 43 || res.Data[1][1] != 50 {
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
	m1.Data[0][0], m2.Data[0][0] = 5, 2
	m1.Data[1][1], m2.Data[1][1] = 7, 4
	res := m1.Subtract(m2)
	if res.Data[0][0] != 3 || res.Data[1][1] != 3 {
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
	m.Data[0][0], m.Data[0][1] = 1, 2
	m.Data[1][0], m.Data[1][1] = 3, 4
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
	m.Data[0][0], m.Data[0][1] = 3, 4
	m.Data[1][0], m.Data[1][1] = 0, 0
	m.Normalize()
	norm := 0.0
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			norm += m.Data[i][j] * m.Data[i][j]
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
