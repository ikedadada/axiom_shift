package usecase

import (
	"testing"

	"axiom_shift/internal/domain"
)

func TestFindValidSeed_Basic(t *testing.T) {
	battleMax := 5
	initialSeed := int64(42)
	pm := domain.NewMatrix(2, 2)
	pm.Data[0][0] = 2.0
	pm.Data[1][1] = 2.0
	player := domain.NewPlayer(pm, 0.5)
	enemy := domain.NewEnemy("Enemy", domain.NewMatrix(2, 2), 0.5)
	seed, rule, playerPath, enemyPath := FindValidSeed(battleMax, initialSeed, player, enemy)
	if rule == nil {
		t.Fatal("rule should not be nil")
	}
	if len(playerPath) != battleMax {
		t.Errorf("playerPath length = %d, want %d", len(playerPath), battleMax)
	}
	if len(enemyPath) != battleMax {
		t.Errorf("enemyPath length = %d, want %d", len(enemyPath), battleMax)
	}
	if seed == 0 {
		t.Error("seed should not be zero (very unlikely)")
	}
	_ = rule.GetMatrix()
}
