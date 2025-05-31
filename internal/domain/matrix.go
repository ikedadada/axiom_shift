package domain

// Matrix represents a mathematical matrix with basic operations.
type Matrix struct {
	Data [][]float64
	Rows int
	Cols int
}

// NewMatrix creates a new Matrix with the given number of rows and columns.
func NewMatrix(data [][]float64) *Matrix {
	rows := len(data)
	cols := 0
	if rows > 0 {
		cols = len(data[0])
	}
	return &Matrix{Data: data, Rows: rows, Cols: cols}
}

// Subtract performs matrix subtraction with another matrix.
func (m *Matrix) Subtract(other *Matrix) *Matrix {
	if m == nil || other == nil || m.Rows == 0 || m.Cols == 0 || other.Rows == 0 || other.Cols == 0 {
		return nil
	}
	if m.Rows != other.Rows || m.Cols != other.Cols {
		return nil
	}
	var data = make([][]float64, m.Rows)
	for i := range data {
		data[i] = make([]float64, m.Cols)
	}
	result := NewMatrix(data)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			result.Data[i][j] = m.Data[i][j] - other.Data[i][j]
		}
	}
	return result
}

// Multiply performs matrix multiplication with another matrix.
func (m *Matrix) Multiply(other *Matrix) *Matrix {
	if m == nil || other == nil || m.Rows == 0 || m.Cols == 0 || other.Rows == 0 || other.Cols == 0 {
		return nil
	}
	if m.Cols != other.Rows {
		return nil
	}
	var data = make([][]float64, m.Rows)
	for i := range data {
		data[i] = make([]float64, m.Cols)
	}
	result := NewMatrix(data)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < other.Cols; j++ {
			for k := 0; k < m.Cols; k++ {
				result.Data[i][j] += m.Data[i][k] * other.Data[k][j]
			}
		}
	}
	return result
}

// GetScalarValue returns a representative scalar value for the matrix (e.g., average of all elements).
func (m *Matrix) GetScalarValue() float64 {
	sum := 0.0
	count := 0
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			sum += m.Data[i][j]
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
	if m == nil || m.Rows == 0 || m.Cols == 0 {
		return
	}
	sumSquares := 0.0
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			sumSquares += m.Data[i][j] * m.Data[i][j]
		}
	}
	if sumSquares == 0 {
		return
	}
	norm := sqrt(sumSquares)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Cols; j++ {
			m.Data[i][j] /= norm
		}
	}
}

// Copy creates a deep copy of the matrix.
func (m *Matrix) Copy() *Matrix {
	if m == nil {
		return nil
	}
	newData := make([][]float64, m.Rows)
	for i := range newData {
		newData[i] = make([]float64, m.Cols)
		copy(newData[i], m.Data[i])
	}
	return &Matrix{Data: newData, Rows: m.Rows, Cols: m.Cols}
}

// sqrt is a helper for square root (for normalization)
func sqrt(x float64) float64 {
	if x == 0 {
		return 0
	}
	if x < 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

func equal(a, b *Matrix) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Rows != b.Rows || a.Cols != b.Cols {
		return false
	}
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			if a.Data[i][j] != b.Data[i][j] {
				return false
			}
		}
	}
	return true
}
