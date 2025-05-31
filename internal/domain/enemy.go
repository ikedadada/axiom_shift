package domain

type Enemy struct {
	Name       string
	Matrix     *Matrix
	GrowthRate float64
}

// NewEnemy creates a new enemy with the given name and initial matrix state.
func NewEnemy(name string, initialMatrix *Matrix, growthRate float64) *Enemy {
	return &Enemy{
		Name:       name,
		Matrix:     initialMatrix,
		GrowthRate: growthRate,
	}
}

// Grow updates the enemy's matrix based on the input value and rule matrix.
func (e *Enemy) Grow(input float64, rule *Matrix) {
	if e.Matrix.Rows == 0 || e.Matrix.Cols == 0 {
		return
	}
	total := e.Matrix.Rows * e.Matrix.Cols
	idx := int(input*float64(total-1) + 0.5)
	targetI := idx / e.Matrix.Cols
	targetJ := idx % e.Matrix.Cols
	for i := 0; i < e.Matrix.Rows; i++ {
		for j := 0; j < e.Matrix.Cols; j++ {
			if i == targetI && j == targetJ {
				e.Matrix.Data[i][j] += 1.0 * e.GrowthRate
			} else {
				e.Matrix.Data[i][j] += 0.1 * e.GrowthRate
			}
		}
	}
	// ルール行列ベース：最大値の要素も強化
	maxVal := rule.Data[0][0]
	maxI, maxJ := 0, 0
	for x := 0; x < rule.Rows; x++ {
		for y := 0; y < rule.Cols; y++ {
			if rule.Data[x][y] > maxVal {
				maxVal = rule.Data[x][y]
				maxI, maxJ = x, y
			}
		}
	}
	e.Matrix.Data[maxI][maxJ] += 1.0 * e.GrowthRate
}

func (e *Enemy) GetMatrix() *Matrix {
	return e.Matrix
}
