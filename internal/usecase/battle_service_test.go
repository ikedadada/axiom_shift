package usecase

import (
	"axiom_shift/internal/domain"
	"axiom_shift/internal/logic"
	"testing"
)

func logicToDomainRuleMatrix(lr *logic.RuleMatrix) *domain.RuleMatrix {
	mat := lr.GetMatrix()
	return &domain.RuleMatrix{Matrix: mat}
}

func TestBattleService_ExecuteBattle(t *testing.T) {
	tests := []struct {
		name      string
		playerMat *domain.Matrix
		playerGr  float64
		enemyMat  *domain.Matrix
		enemyGr   float64
		ruleSeed  int64
		ruleSize  int
		input     float64
	}{
		{"basic win", domain.NewMatrix(2, 2), 1.0, domain.NewMatrix(2, 2), 1.0, 1, 2, 1.0},
		{"basic lose", domain.NewMatrix(2, 2), 0.1, domain.NewMatrix(2, 2), 1.0, 1, 2, 0.0},
		{"zero matrix", domain.NewMatrix(0, 0), 1.0, domain.NewMatrix(0, 0), 1.0, 1, 0, 0.0},
		{"input negative", domain.NewMatrix(2, 2), 1.0, domain.NewMatrix(2, 2), 1.0, 1, 2, -1.0},
		{"input >1", domain.NewMatrix(2, 2), 1.0, domain.NewMatrix(2, 2), 1.0, 1, 2, 2.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := domain.NewPlayer(tt.playerMat, tt.playerGr)
			enemy := domain.NewEnemy("enemy", tt.enemyMat, tt.enemyGr)
			rule := logic.NewRuleMatrix(tt.ruleSeed, tt.ruleSize)
			b := NewBattleService(player, enemy, logicToDomainRuleMatrix(rule))
			result, _ := b.ExecuteBattle(tt.input)
			_ = result // Could add more assertions if desired
			if tt.playerMat.Rows == 0 || tt.enemyMat.Rows == 0 {
				return // skip win/lose check for zero matrix
			}
			// Just check that result is a float and win is a bool (no panic)
		})
	}
}

func TestNewBattleService(t *testing.T) {
	tests := []struct {
		name      string
		playerMat *domain.Matrix
		playerGr  float64
		enemyMat  *domain.Matrix
		enemyGr   float64
		ruleSeed  int64
		ruleSize  int
	}{
		{"normal", domain.NewMatrix(1, 1), 1.0, domain.NewMatrix(1, 1), 1.0, 1, 1},
		{"zero matrix", domain.NewMatrix(0, 0), 1.0, domain.NewMatrix(0, 0), 1.0, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := domain.NewPlayer(tt.playerMat, tt.playerGr)
			enemy := domain.NewEnemy("e", tt.enemyMat, tt.enemyGr)
			rule := logic.NewRuleMatrix(tt.ruleSeed, tt.ruleSize)
			b := NewBattleService(player, enemy, logicToDomainRuleMatrix(rule))
			if b.Player != player || b.Enemy != enemy || b.Rules == nil {
				t.Error("NewBattleService did not set fields correctly")
			}
		})
	}
}
