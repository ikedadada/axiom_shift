package usecase

import (
	"axiom_shift/internal/domain"
	"testing"
)

func TestBattleService_ExecuteBattle(t *testing.T) {
	tests := []struct {
		name      string
		playerMat [][]float64
		playerGr  float64
		enemyMat  [][]float64
		enemyGr   float64
		ruleSeed  int64
		ruleSize  int
		input     float64
		wantWin   bool
	}{
		{"basic win", [][]float64{{2, 2}, {2, 2}}, 1.0, [][]float64{{0, 0}, {0, 0}}, 1.0, 1, 2, 1.0, true},
		{"basic lose", [][]float64{{0, 0}, {0, 0}}, 0.1, [][]float64{{2, 2}, {2, 2}}, 1.0, 1, 2, 0.0, false},
		{"zero matrix", [][]float64{}, 1.0, [][]float64{}, 1.0, 1, 0, 0.0, false},
		{"input negative", [][]float64{{2, 2}, {2, 2}}, 1.0, [][]float64{{0, 0}, {0, 0}}, 1.0, 1, 2, -1.0, true},
		{"input >1", [][]float64{{2, 2}, {2, 2}}, 1.0, [][]float64{{0, 0}, {0, 0}}, 1.0, 1, 2, 2.0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := domain.NewPlayer(domain.NewMatrix(tt.playerMat), tt.playerGr)
			enemy := domain.NewEnemy("enemy", domain.NewMatrix(tt.enemyMat), tt.enemyGr)
			rule := domain.NewRuleMatrix(tt.ruleSeed, tt.ruleSize)
			b := NewBattleService(player, enemy, rule)
			result, win := b.ExecuteBattle(tt.input)
			if len(tt.playerMat) == 0 || len(tt.enemyMat) == 0 {
				if win {
					t.Error("win should be false for zero matrix")
				}
				return
			}
			if win != tt.wantWin {
				t.Errorf("win = %v, want %v", win, tt.wantWin)
			}
			if result == 0 && tt.wantWin {
				t.Error("result should not be zero for win case")
			}
		})
	}
}

func TestNewBattleService(t *testing.T) {
	tests := []struct {
		name      string
		playerMat [][]float64
		playerGr  float64
		enemyMat  [][]float64
		enemyGr   float64
		ruleSeed  int64
		ruleSize  int
	}{
		{"normal", [][]float64{{0}}, 1.0, [][]float64{{0}}, 1.0, 1, 1},
		{"zero matrix", [][]float64{}, 1.0, [][]float64{}, 1.0, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := domain.NewPlayer(domain.NewMatrix(tt.playerMat), tt.playerGr)
			enemy := domain.NewEnemy("e", domain.NewMatrix(tt.enemyMat), tt.enemyGr)
			rule := domain.NewRuleMatrix(tt.ruleSeed, tt.ruleSize)
			b := NewBattleService(player, enemy, rule)
			if b.Player != player || b.Enemy != enemy || b.Rules == nil {
				t.Error("NewBattleService did not set fields correctly")
			}
		})
	}
}

func TestBattleService_calculateBattleOutcome_GuardClauses(t *testing.T) {
	tests := []struct {
		name      string
		playerMat [][]float64
		enemyMat  [][]float64
		ruleMat   [][]float64
		wantZero  bool
	}{
		{"all nil", nil, nil, nil, true},
		{"player nil", nil, [][]float64{{0, 0}, {0, 0}}, [][]float64{{1, 2}, {3, 4}}, true},
		{"enemy nil", [][]float64{{0, 0}, {0, 0}}, nil, [][]float64{{1, 2}, {3, 4}}, true},
		{"rule nil", [][]float64{{0, 0}, {0, 0}}, [][]float64{{0, 0}, {0, 0}}, nil, true},
		{"rule empty", [][]float64{{0, 0}, {0, 0}}, [][]float64{{0, 0}, {0, 0}}, [][]float64{}, true},
		{"rule row empty", [][]float64{{0, 0}, {0, 0}}, [][]float64{{0, 0}, {0, 0}}, [][]float64{{}}, true},
		{"multiply nil", [][]float64{{0, 0}, {0, 0}}, [][]float64{{0, 0}, {0, 0}}, [][]float64{{1}, {2}}, true},
		{"subtract nil", [][]float64{{0, 0}, {0, 0}}, [][]float64{{0, 0}, {0, 0}, {0, 0}}, [][]float64{{1, 2}, {3, 4}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BattleService{
				Player: &domain.Player{MatrixState: domain.NewMatrix(tt.playerMat)},
				Enemy:  &domain.Enemy{MatrixState: domain.NewMatrix(tt.enemyMat)},
				Rules:  &domain.RuleMatrix{Matrix: domain.NewMatrix(tt.ruleMat)},
			}
			got := b.calculateBattleOutcome()
			if tt.wantZero && got != 0 {
				t.Errorf("Expected 0 for guard clause, got %v", got)
			}
		})
	}
}

func TestBattleService_DoBattleTurn_Branches(t *testing.T) {
	tests := []struct {
		name      string
		playerMat [][]float64
		playerGr  float64
		enemyMat  [][]float64
		enemyGr   float64
		ruleSeed  int64
		ruleSize  int
		input     float64
		battles   int
	}{
		{"win false branch", [][]float64{{0, 0}, {0, 0}}, 0.1, [][]float64{{0, 0}, {0, 0}}, 1.0, 1, 2, 0.0, 0},
		{"maxTry loop", [][]float64{{0, 0}, {0, 0}}, 1.0, [][]float64{{0, 0}, {0, 0}}, 0.1, 1, 2, 1.0, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := domain.NewPlayer(domain.NewMatrix(tt.playerMat), tt.playerGr)
			enemy := domain.NewEnemy("enemy", domain.NewMatrix(tt.enemyMat), tt.enemyGr)
			rule := domain.NewRuleMatrix(tt.ruleSeed, tt.ruleSize)
			b := NewBattleService(player, enemy, rule)
			for i := 0; i < tt.battles; i++ {
				_, _ = b.DoBattleTurn(tt.input, i)
			}
		})
	}
}
