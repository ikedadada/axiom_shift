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

// 1ターン分のバトル進行（プレイヤー成長・敵成長含む）
func (b *BattleService) DoBattleTurn(input float64, battleCount int) (float64, bool) {
	// プレイヤー成長率をバトル回数で調整
	b.Player.GrowthRate = 0.5 + 0.1*float64(battleCount)
	result, win := b.ExecuteBattle(input)
	// プレイヤーが勝った場合のみ敵が成長
	if win {
		ruleMatrix := b.Rules.Matrix
		m := domain.NewMatrix(len(ruleMatrix), len(ruleMatrix[0]))
		for i := range ruleMatrix {
			for j := range ruleMatrix[i] {
				m.Data[i][j] = ruleMatrix[i][j]
			}
		}
		maxTry := 10
		for i := 0; i < maxTry; i++ {
			b.Enemy.Grow(input, m)
			// 成長後に再度バトル判定
			_, winTmp := b.ExecuteBattle(input)
			if winTmp {
				continue // まだ勝てない→さらにGrow
			} else {
				break // 勝てなくなったら終了
			}
		}
	}
	return result, win
}

func (b *BattleService) calculateBattleOutcome() float64 {
	playerMatrix := b.Player.GetMatrix()
	enemyMatrix := b.Enemy.GetMatrix()
	ruleMatrix := b.Rules.Matrix // [][]float64
	if playerMatrix == nil || enemyMatrix == nil || ruleMatrix == nil || len(ruleMatrix) == 0 || len(ruleMatrix[0]) == 0 {
		return 0
	}
	m := domain.NewMatrix(len(ruleMatrix), len(ruleMatrix[0]))
	for i := range ruleMatrix {
		for j := range ruleMatrix[i] {
			m.Data[i][j] = ruleMatrix[i][j]
		}
	}
	outcome := playerMatrix.Multiply(m)
	if outcome == nil {
		return 0
	}
	outcome = outcome.Subtract(enemyMatrix)
	if outcome == nil {
		return 0
	}
	return outcome.GetScalarValue()
}
