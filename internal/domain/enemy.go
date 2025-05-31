package domain

type Enemy struct {
	Name         string
	initialState Matrix // 初期状態のマトリックス
	MatrixState  *Matrix
	GrowthRate   float64
}

// NewEnemy creates a new enemy with the given name and initial matrix state.
func NewEnemy(name string, initialMatrix *Matrix, growthRate float64) *Enemy {
	if initialMatrix == nil {
		return &Enemy{
			Name:         name,
			initialState: Matrix{},
			MatrixState:  nil,
			GrowthRate:   growthRate,
		}
	}
	return &Enemy{
		Name:         name,
		initialState: *initialMatrix, // Store the initial state
		MatrixState:  initialMatrix,
		GrowthRate:   growthRate,
	}
}

// Reset resets the enemy's matrix to its initial state.
func (e *Enemy) Reset() {
	e.MatrixState = e.initialState.Copy()
}

// Grow updates the enemy's matrix based on the input value and rule matrix.
func (e *Enemy) Grow(input float64, rule *Matrix) {
	if e.MatrixState == nil || rule == nil || e.MatrixState.Rows == 0 || e.MatrixState.Cols == 0 || rule.Rows == 0 || rule.Cols == 0 {
		return
	}
	total := e.MatrixState.Rows * e.MatrixState.Cols
	idx := int(input*float64(total-1) + 0.5)
	if idx < 0 {
		idx = 0
	}
	if idx >= total {
		idx = total - 1
	}
	targetI := idx / e.MatrixState.Cols
	targetJ := idx % e.MatrixState.Cols
	for i := 0; i < e.MatrixState.Rows; i++ {
		for j := 0; j < e.MatrixState.Cols; j++ {
			if i == targetI && j == targetJ {
				e.MatrixState.Data[i][j] += 1.0 * e.GrowthRate
			} else {
				e.MatrixState.Data[i][j] += 0.1 * e.GrowthRate
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
	e.MatrixState.Data[maxI][maxJ] += 0.5 * e.GrowthRate
}

func (e *Enemy) GetMatrix() *Matrix {
	return e.MatrixState
}
