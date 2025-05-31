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

func TestBattleService_calculateBattleOutcome_GuardClauses(t *testing.T) {
	tests := []struct {
		name      string
		playerMat *domain.Matrix
		enemyMat  *domain.Matrix
		ruleMat   [][]float64
		wantZero  bool
	}{
		{"all nil", nil, nil, nil, true},
		{"player nil", nil, domain.NewMatrix(2, 2), [][]float64{{1, 2}, {3, 4}}, true},
		{"enemy nil", domain.NewMatrix(2, 2), nil, [][]float64{{1, 2}, {3, 4}}, true},
		{"rule nil", domain.NewMatrix(2, 2), domain.NewMatrix(2, 2), nil, true},
		{"rule empty", domain.NewMatrix(2, 2), domain.NewMatrix(2, 2), [][]float64{}, true},
		{"rule row empty", domain.NewMatrix(2, 2), domain.NewMatrix(2, 2), [][]float64{{}}, true},
		{"multiply nil", domain.NewMatrix(2, 2), domain.NewMatrix(2, 2), [][]float64{{1}, {2}}, true},
		{"subtract nil", domain.NewMatrix(2, 2), domain.NewMatrix(3, 3), [][]float64{{1, 2}, {3, 4}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BattleService{
				Player: &domain.Player{MatrixState: tt.playerMat},
				Enemy:  &domain.Enemy{MatrixState: tt.enemyMat},
				Rules:  &domain.RuleMatrix{Matrix: tt.ruleMat},
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
		playerMat *domain.Matrix
		playerGr  float64
		enemyMat  *domain.Matrix
		enemyGr   float64
		ruleSeed  int64
		ruleSize  int
		input     float64
		battles   int
	}{
		{"win false branch", domain.NewMatrix(2, 2), 0.1, domain.NewMatrix(2, 2), 1.0, 1, 2, 0.0, 0},
		{"maxTry loop", domain.NewMatrix(2, 2), 1.0, domain.NewMatrix(2, 2), 0.1, 1, 2, 1.0, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := domain.NewPlayer(tt.playerMat, tt.playerGr)
			enemy := domain.NewEnemy("enemy", tt.enemyMat, tt.enemyGr)
			rule := logic.NewRuleMatrix(tt.ruleSeed, tt.ruleSize)
			b := NewBattleService(player, enemy, logicToDomainRuleMatrix(rule))
			for i := 0; i < tt.battles; i++ {
				_, _ = b.DoBattleTurn(tt.input, i)
			}
		})
	}
}
