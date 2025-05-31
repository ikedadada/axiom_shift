package game

import (
	"axiom_shift/internal/logic"
	"testing"
)

func TestBattleService_ExecuteBattle(t *testing.T) {
	player := NewCharacter(NewMatrix(2, 2), 1.0)
	enemy := NewEnemy("enemy", NewMatrix(2, 2), 1.0)
	rule := logic.NewRuleMatrix(1, 2)
	b := &BattleService{player: player, enemy: enemy, rules: rule}
	result, win := b.ExecuteBattle(1.0)
	if result == 0 && win {
		t.Error("BattleService ExecuteBattle logic may be incorrect")
	}
}

func TestNewBattleService(t *testing.T) {
	player := NewCharacter(NewMatrix(1, 1), 1.0)
	enemy := NewEnemy("e", NewMatrix(1, 1), 1.0)
	rule := logic.NewRuleMatrix(1, 1)
	b := NewBattleService(player, enemy, rule)
	if b.player != player || b.enemy != enemy || b.rules != rule {
		t.Error("NewBattleService did not set fields correctly")
	}
}
