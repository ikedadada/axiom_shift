package usecase

import (
	"axiom_shift/internal/domain"
)

type BattleService struct {
	Player *domain.Player
	Enemy  *domain.Enemy
	Rules  *domain.RuleMatrix
}

func NewBattleService(player *domain.Player, enemy *domain.Enemy, rules *domain.RuleMatrix) *BattleService {
	return &BattleService{
		Player: player,
		Enemy:  enemy,
		Rules:  rules,
	}
}

func (b *BattleService) ExecuteBattle(playerInput float64) (float64, bool) {
	b.Player.UpdateMatrix(playerInput)
	b.Player.GetMatrix().Normalize()
	b.Enemy.GetMatrix().Normalize()
	result := b.calculateBattleOutcome()
	playerWins := result > 0
	return result, playerWins
}

func (b *BattleService) calculateBattleOutcome() float64 {
	playerMatrix := b.Player.GetMatrix()
	enemyMatrix := b.Enemy.GetMatrix()
	ruleMatrix := b.Rules.Matrix // [][]float64
	m := domain.NewMatrix(len(ruleMatrix), len(ruleMatrix[0]))
	for i := range ruleMatrix {
		for j := range ruleMatrix[i] {
			m.Data[i][j] = ruleMatrix[i][j]
		}
	}
	outcome := playerMatrix.Multiply(m).Subtract(enemyMatrix)
	return outcome.GetScalarValue()
}
