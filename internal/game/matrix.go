package game

type Matrix struct {
	data [][]float64
	rows int
	cols int
}

// NewMatrix creates a new Matrix with the given number of rows and columns.
func NewMatrix(rows, cols int) *Matrix {
	data := make([][]float64, rows)
	for i := range data {
		data[i] = make([]float64, cols)
	}
	return &Matrix{data: data, rows: rows, cols: cols}
}

// Subtract performs matrix subtraction with another matrix.
func (m *Matrix) Subtract(other *Matrix) *Matrix {
	if m.rows != other.rows || m.cols != other.cols {
		return nil // or handle error
	}
	result := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] - other.data[i][j]
		}
	}
	return result
}

// Multiply performs matrix multiplication with another matrix.
func (m *Matrix) Multiply(other *Matrix) *Matrix {
	if m.cols != other.rows {
		return nil // or handle error
	}
	result := NewMatrix(m.rows, other.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < other.cols; j++ {
			for k := 0; k < m.cols; k++ {
				result.data[i][j] += m.data[i][k] * other.data[k][j]
			}
		}
	}
	return result
}

// GetScalarValue returns a representative scalar value for the matrix (e.g., average of all elements).
func (m *Matrix) GetScalarValue() float64 {
	sum := 0.0
	count := 0
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			sum += m.data[i][j]
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

// Normalize normalizes the matrix so that its L2 norm becomes 1 (unless norm is 0).
func (m *Matrix) Normalize() {
	sumSquares := 0.0
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			sumSquares += m.data[i][j] * m.data[i][j]
		}
	}
	if sumSquares == 0 {
		return
	}
	norm := sqrt(sumSquares)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			m.data[i][j] /= norm
		}
	}
}

// sqrt is a helper for square root (for normalization)
func sqrt(x float64) float64 {
	// Use Newton's method for simplicity
	if x == 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
