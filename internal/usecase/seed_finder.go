package usecase

import (
	"axiom_shift/internal/domain"
	"axiom_shift/internal/logic"
	"fmt"
	"math"
	"math/rand"
)

// FindValidSeed: battleMax回数に対し有効なseed/rule/playerPath/enemyPathを返す
func FindValidSeed(battleMax int, initialSeed int64, player *domain.Player, enemy *domain.Enemy) (int64, *logic.RuleMatrix, []float64, []float64) {
	const (
		roughSamples = 2000
		deepSamples  = 20000
		mctsWidth    = 3
		mctsMaxNodes = 100000
	)
	// Wilson score interval (近似)
	betaCI := func(wins, n int, alpha float64) (float64, float64) {
		if n == 0 {
			return 0, 1
		}
		p := float64(wins) / float64(n)
		z := 2.576 // 99%信頼区間
		denom := 1 + z*z/float64(n)
		center := p + z*z/(2*float64(n))
		pm := z * (math.Sqrt(p*(1-p)/float64(n) + z*z/(4*float64(n)*float64(n))))
		low := (center - pm) / denom
		high := (center + pm) / denom
		if low < 0 {
			low = 0
		}
		if high > 1 {
			high = 1
		}
		return low, high
	}
	// サンプリングによる勝率推定
	simulateSamples := func(rule *logic.RuleMatrix, samples int) (int, int) {
		playerWins := 0
		for s := 0; s < samples; s++ {
			player.Reset()
			enemy.Reset()
			service := NewBattleService(player, enemy, &domain.RuleMatrix{Matrix: rule.GetMatrix()})
			inputs := make([]float64, battleMax)
			for i := 0; i < battleMax; i++ {
				inputs[i] = float64(rand.Intn(10)) / 9.0
			}
			var win bool
			for battle := 0; battle < battleMax; battle++ {
				_, win = service.DoBattleTurn(inputs[battle], battle)
			}
			if win {
				playerWins++
			}
		}
		return playerWins, samples
	}
	// ProofPhase: MCTS風DFS（幅3, 最大ノード10万）
	proofPhase := func(rule *logic.RuleMatrix) (bool, []float64, []float64) {
		type node struct {
			depth  int
			inputs []float64
		}
		var foundPlayer, foundEnemy bool
		var playerPath, enemyPath []float64
		var nodes int
		var dfs func(n node)
		dfs = func(n node) {
			if foundPlayer && foundEnemy {
				return
			}
			if nodes >= mctsMaxNodes {
				return
			}
			nodes++
			if n.depth == battleMax {
				player.Reset()
				enemy.Reset()
				service := NewBattleService(player, enemy, &domain.RuleMatrix{Matrix: rule.GetMatrix()})
				var win bool
				for battle := 0; battle < battleMax; battle++ {
					_, win = service.DoBattleTurn(n.inputs[battle], battle)
				}
				if win && !foundPlayer {
					foundPlayer = true
					playerPath = append([]float64{}, n.inputs...)
				}
				if !win && !foundEnemy {
					foundEnemy = true
					enemyPath = append([]float64{}, n.inputs...)
				}
				return
			}
			choices := []float64{0.0, 0.5, 1.0}
			choices = append(choices, float64(rand.Intn(10))/9.0)
			for _, v := range choices {
				dfs(node{n.depth + 1, append(n.inputs, v)})
				if foundPlayer && foundEnemy {
					return
				}
				if nodes >= mctsMaxNodes {
					return
				}
			}
		}
		dfs(node{0, []float64{}})
		return foundPlayer && foundEnemy, playerPath, enemyPath
	}

	rand.Seed(initialSeed)
	var debug_SearchSeedCount = 0
	for {
		seedCandidate := logic.NewSeedManager().GetSeed()
		rule := logic.NewRuleMatrix(seedCandidate, 2)
		// RoughFilter
		playerWins, n := simulateSamples(rule, roughSamples)
		pHat := float64(playerWins) / float64(n)
		low, high := betaCI(playerWins, n, 0.01)
		if !(low < 0.99 && high > 0.01) {
			debug_SearchSeedCount++
			fmt.Printf("[Rough] Seed %d rejected: CI=(%.3f,%.3f)\n", seedCandidate, low, high)
			continue
		}
		// DeepFilter
		playerWins, n = simulateSamples(rule, deepSamples)
		pHat = float64(playerWins) / float64(n)
		if !(pHat > 0 && pHat < 1) {
			debug_SearchSeedCount++
			fmt.Printf("[Deep] Seed %d rejected: p̂=%.3f\n", seedCandidate, pHat)
			continue
		}
		// ProofPhase
		ok, playerPath, enemyPath := proofPhase(rule)
		debug_SearchSeedCount++
		fmt.Printf("[Proof] Seed %d: ok=%v, playerWinPath=%v, enemyWinPath=%v\n", seedCandidate, ok, playerPath, enemyPath)
		if ok {
			return seedCandidate, rule, playerPath, enemyPath
		}
	}
}
