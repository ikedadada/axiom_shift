package usecase

import (
	"testing"

	"axiom_shift/internal/domain"
)

func TestFindValidSeed_Basic(t *testing.T) {
	tests := []struct {
		name        string
		battleMax   int
		initialSeed int64
		playerMat   *domain.Matrix
		playerGr    float64
		enemyMat    *domain.Matrix
		enemyGr     float64
	}{
		{"basic", 5, 42, domain.NewMatrix(2, 2), 0.5, domain.NewMatrix(2, 2), 0.5},
		{"different seed", 5, 99, domain.NewMatrix(2, 2), 0.5, domain.NewMatrix(2, 2), 0.5},
		// TODO: larger matrix does not work yet
		// {"larger matrix", 5, 42, domain.NewMatrix(3, 3), 0.5, domain.NewMatrix(3, 3), 0.5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := tt.playerMat
			pm.Data[0][0] = 2.0
			if pm.Rows > 1 && pm.Cols > 1 {
				pm.Data[1][1] = 2.0
			}
			player := domain.NewPlayer(pm, tt.playerGr)
			enemy := domain.NewEnemy("Enemy", tt.enemyMat, tt.enemyGr)
			seed, rule, playerPath, enemyPath := FindValidSeed(tt.battleMax, tt.initialSeed, player, enemy)
			if rule == nil {
				t.Fatal("rule should not be nil")
			}
			if len(playerPath) != tt.battleMax {
				t.Errorf("playerPath length = %d, want %d", len(playerPath), tt.battleMax)
			}
			if len(enemyPath) != tt.battleMax {
				t.Errorf("enemyPath length = %d, want %d", len(enemyPath), tt.battleMax)
			}
			if seed == 0 {
				t.Error("seed should not be zero (very unlikely)")
			}
			_ = rule.GetMatrix()
		})
	}
}

func TestFindValidSeed_GuardCases(t *testing.T) {
	tests := []struct {
		name      string
		battleMax int
		player    *domain.Player
		enemy     *domain.Enemy
		seed      int64
		wantPanic bool
	}{
		{"zero battleMax", 0, domain.NewPlayer(domain.NewMatrix(2, 2), 0.5), domain.NewEnemy("E", domain.NewMatrix(2, 2), 0.5), 42, true},
		{"nil player", 5, nil, domain.NewEnemy("E", domain.NewMatrix(2, 2), 0.5), 42, true},
		{"nil enemy", 5, domain.NewPlayer(domain.NewMatrix(2, 2), 0.5), nil, 42, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tt.wantPanic && r == nil {
					t.Errorf("Expected panic but did not panic")
				}
				if !tt.wantPanic && r != nil {
					t.Errorf("Unexpected panic: %v", r)
				}
			}()
			_, _, _, _ = FindValidSeed(tt.battleMax, tt.seed, tt.player, tt.enemy)
		})
	}
}

func TestFindValidSeed_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		battleMax int
		seed      int64
		player    *domain.Player
		enemy     *domain.Enemy
	}{
		{"battleMax 1", 1, 1, domain.NewPlayer(domain.NewMatrix(2, 2), 0.5), domain.NewEnemy("E", domain.NewMatrix(2, 2), 0.5)},
		{"large battleMax", 6, 42, domain.NewPlayer(domain.NewMatrix(2, 2), 0.5), domain.NewEnemy("E", domain.NewMatrix(2, 2), 0.5)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, rule, playerPath, enemyPath := FindValidSeed(tt.battleMax, tt.seed, tt.player, tt.enemy)
			if rule == nil {
				t.Fatal("rule should not be nil")
			}
			if len(playerPath) != tt.battleMax {
				t.Errorf("playerPath length = %d, want %d", len(playerPath), tt.battleMax)
			}
			if len(enemyPath) != tt.battleMax {
				t.Errorf("enemyPath length = %d, want %d", len(enemyPath), tt.battleMax)
			}
		})
	}
}
