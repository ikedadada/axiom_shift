package game

type Enemy struct {
	Name       string
	Matrix     *Matrix // *Matrix型に変更
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
// プレイヤーが勝った場合、入力値ベースとルール行列ベースの両方で成長する
func (e *Enemy) Grow(input float64, rule *Matrix) {
	if e.Matrix.rows == 0 || e.Matrix.cols == 0 {
		return
	}
	total := e.Matrix.rows * e.Matrix.cols
	idx := int(input*float64(total-1) + 0.5)
	targetI := idx / e.Matrix.cols
	targetJ := idx % e.Matrix.cols
	for i := 0; i < e.Matrix.rows; i++ {
		for j := 0; j < e.Matrix.cols; j++ {
			if i == targetI && j == targetJ {
				e.Matrix.data[i][j] += 1.0 * e.GrowthRate
			} else {
				e.Matrix.data[i][j] += 0.1 * e.GrowthRate
			}
		}
	}
	// ルール行列ベース：最大値の要素も強化（従来通り）
	maxVal := rule.data[0][0]
	maxI, maxJ := 0, 0
	for x := 0; x < rule.rows; x++ {
		for y := 0; y < rule.cols; y++ {
			if rule.data[x][y] > maxVal {
				maxVal = rule.data[x][y]
				maxI, maxJ = x, y
			}
		}
	}
	e.Matrix.data[maxI][maxJ] += 1.0 * e.GrowthRate
}

// GetMatrix returns the current matrix state of the enemy.
func (e *Enemy) GetMatrix() *Matrix {
	return e.Matrix
}
