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
	player := domain.NewPlayer(domain.NewMatrix(2, 2), 1.0)
	enemy := domain.NewEnemy("enemy", domain.NewMatrix(2, 2), 1.0)
	rule := logic.NewRuleMatrix(1, 2)
	b := NewBattleService(player, enemy, logicToDomainRuleMatrix(rule))
	result, win := b.ExecuteBattle(1.0)
	if result == 0 && win {
		t.Error("BattleService ExecuteBattle logic may be incorrect")
	}
}

func TestNewBattleService(t *testing.T) {
	player := domain.NewPlayer(domain.NewMatrix(1, 1), 1.0)
	enemy := domain.NewEnemy("e", domain.NewMatrix(1, 1), 1.0)
	rule := logic.NewRuleMatrix(1, 1)
	b := NewBattleService(player, enemy, logicToDomainRuleMatrix(rule))
	if b.Player != player || b.Enemy != enemy || b.Rules == nil {
		t.Error("NewBattleService did not set fields correctly")
	}
}
