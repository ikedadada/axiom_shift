package game

type Character struct {
	MatrixState *Matrix // Matrix型に変更
	GrowthRate  float64
}

// NewCharacter initializes a new Character with a given initial matrix state and growth rate.
func NewCharacter(initialState *Matrix, growthRate float64) *Character {
	return &Character{
		MatrixState: initialState,
		GrowthRate:  growthRate,
	}
}

// UpdateMatrix updates the character's matrix state based on the input value.
func (c *Character) UpdateMatrix(input float64) {
	if c.MatrixState.rows == 0 || c.MatrixState.cols == 0 {
		return
	}
	total := c.MatrixState.rows * c.MatrixState.cols
	idx := int(input*float64(total-1) + 0.5)
	targetI := idx / c.MatrixState.cols
	targetJ := idx % c.MatrixState.cols
	for i := 0; i < c.MatrixState.rows; i++ {
		for j := 0; j < c.MatrixState.cols; j++ {
			if i == targetI && j == targetJ {
				c.MatrixState.data[i][j] += 1.0 * c.GrowthRate // 入力値で決まるセルは大きく成長
			} else {
				c.MatrixState.data[i][j] += 0.1 * c.GrowthRate // 他セルは微小成長
			}
		}
	}
}

// GetMatrix returns the current matrix state of the character.
func (c *Character) GetMatrix() *Matrix {
	return c.MatrixState
}
