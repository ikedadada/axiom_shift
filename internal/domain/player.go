package domain

type Player struct {
	MatrixState *Matrix
	GrowthRate  float64
}

// NewPlayer initializes a new Player with a given initial matrix state and growth rate.
func NewPlayer(initialState *Matrix, growthRate float64) *Player {
	return &Player{
		MatrixState: initialState,
		GrowthRate:  growthRate,
	}
}

// UpdateMatrix updates the player's matrix state based on the input value.
func (p *Player) UpdateMatrix(input float64) {
	if p.MatrixState.Rows == 0 || p.MatrixState.Cols == 0 {
		return
	}
	total := p.MatrixState.Rows * p.MatrixState.Cols
	idx := int(input*float64(total-1) + 0.5)
	targetI := idx / p.MatrixState.Cols
	targetJ := idx % p.MatrixState.Cols
	for i := 0; i < p.MatrixState.Rows; i++ {
		for j := 0; j < p.MatrixState.Cols; j++ {
			if i == targetI && j == targetJ {
				p.MatrixState.Data[i][j] += 1.0 * p.GrowthRate
			} else {
				p.MatrixState.Data[i][j] += 0.1 * p.GrowthRate
			}
		}
	}
}

func (p *Player) GetMatrix() *Matrix {
	return p.MatrixState
}
