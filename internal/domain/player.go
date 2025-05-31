package domain

type Player struct {
	initialState Matrix // 初期状態のマトリックス
	MatrixState  *Matrix
	GrowthRate   float64
}

// NewPlayer initializes a new Player with a given initial matrix state and growth rate.
func NewPlayer(initialState *Matrix, growthRate float64) *Player {
	if initialState == nil {
		return &Player{
			initialState: Matrix{},
			MatrixState:  nil,
			GrowthRate:   growthRate,
		}
	}
	return &Player{
		initialState: *initialState.Copy(),
		MatrixState:  initialState.Copy(),
		GrowthRate:   growthRate,
	}
}

// Reset resets the player's matrix state to the initial state.
func (p *Player) Reset() {
	p.MatrixState = p.initialState.Copy()
}

// UpdateMatrix updates the player's matrix state based on the input value.
func (p *Player) UpdateMatrix(input float64) {
	if p.MatrixState == nil || p.MatrixState.Rows == 0 || p.MatrixState.Cols == 0 {
		return
	}
	total := p.MatrixState.Rows * p.MatrixState.Cols
	idx := int(input*float64(total-1) + 0.5)
	if idx < 0 {
		idx = 0
	}
	if idx >= total {
		idx = total - 1
	}
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
