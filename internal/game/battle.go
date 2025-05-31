package game

import (
	"axiom_shift/internal/logic"
)

type BattleService struct {
	player *Character
	enemy  *Enemy
	rules  *logic.RuleMatrix
}

func NewBattleService(player *Character, enemy *Enemy, rules *logic.RuleMatrix) *BattleService {
	return &BattleService{
		player: player,
		enemy:  enemy,
		rules:  rules,
	}
}

func (b *BattleService) ExecuteBattle(playerInput float64) (float64, bool) {
	// Update player character's matrix based on input
	b.player.UpdateMatrix(playerInput)

	// 正規化を追加: プレイヤー・敵の行列を正規化
	b.player.GetMatrix().Normalize()
	b.enemy.GetMatrix().Normalize()

	// Perform battle logic using matrix operations
	result := b.calculateBattleOutcome()

	// Determine if the player wins
	playerWins := result > 0 // 0を基準に勝敗判定

	return result, playerWins
}

func (b *BattleService) calculateBattleOutcome() float64 {
	// Perform matrix operations to determine the outcome
	playerMatrix := b.player.GetMatrix()
	enemyMatrix := b.enemy.GetMatrix()

	// ルール行列を*Matrix型で受け取る想定に修正
	ruleMatrix := b.rules.GetMatrix() // [][]float64
	rule := &Matrix{data: ruleMatrix, rows: len(ruleMatrix), cols: len(ruleMatrix[0])}

	outcome := playerMatrix.Multiply(rule).Subtract(enemyMatrix)

	// Return a scalar value representing the battle outcome
	return outcome.GetScalarValue()
}
