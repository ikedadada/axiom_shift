package usecase

import (
	"axiom_shift/internal/domain"
	"axiom_shift/internal/logic"
	"fmt"
	"math"
	"math/rand"
)

// FindValidSeed: battleMax 回のバトルで双方に勝ちパターンが存在する seed / rule / playerPath / enemyPath を返す
func FindValidSeed(battleMax int, player *domain.Player, enemy *domain.Enemy) (int64, []int, []int, error) {
	if battleMax <= 0 || player == nil || enemy == nil {
		panic("Invalid parameters: battleMax must be > 0, player and enemy must not be nil")
	}
	// 行列サイズに応じてパラメータ自動調整
	size := 2
	if player.MatrixState != nil && player.MatrixState.Rows > 0 {
		size = player.MatrixState.Rows
	}
	// サンプリング数・ノード数をサイズ依存で調整
	roughSamples := 50 * size * size
	deepSamples := 200 * size * size
	mctsWidth := 3
	mctsMaxNodes := 1000 * size * size
	maxTries := 1000

	// Wilson score interval (近似) で勝率信頼区間を求める
	betaCI := func(wins, n int) (float64, float64) {
		if n == 0 {
			return 0, 1
		}
		p := float64(wins) / float64(n)
		z := 2.576 // 99% 信頼区間
		denom := 1 + z*z/float64(n)
		center := p + z*z/(2*float64(n))
		pm := z * math.Sqrt(p*(1-p)/float64(n)+z*z/(4*float64(n)*float64(n)))
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
	simulateSamples := func(rule *domain.RuleMatrix, samples int) (int, int) {
		playerWins := 0
		for s := 0; s < samples; s++ {
			player.Reset()
			enemy.Reset()
			service := NewBattleService(player, enemy, rule)

			inputs := make([]int, battleMax)
			for i := range inputs {
				inputs[i] = rand.Intn(10)
			}

			var win bool
			for battle := 0; battle < battleMax; battle++ {
				_, win = service.DoBattleTurn(float64(inputs[battle])/9, battle)
			}
			if win {
				playerWins++
			}
		}
		return playerWins, samples
	}

	// ProofPhase: DFS＋ビーム幅＋分岐シャッフルで多様な勝ちパスを探索
	proofPhase := func(rule *domain.RuleMatrix) (bool, []int, []int) {
		type node struct {
			depth  int
			inputs []int
		}

		var (
			playerPaths [][]int
			enemyPaths  [][]int
			nodes       int
		)

		var dfs func(n node)
		dfs = func(n node) {
			if nodes >= mctsMaxNodes {
				return
			}
			nodes++

			// 末端まで到達したら勝敗を判定
			if n.depth == battleMax {
				player.Reset()
				enemy.Reset()
				service := NewBattleService(player, enemy, rule)

				var win bool
				for battle := 0; battle < battleMax; battle++ {
					_, win = service.DoBattleTurn(float64(n.inputs[battle])/9, battle)
				}
				if win {
					// プレイヤー勝利パス
					playerPaths = append(playerPaths, append([]int(nil), n.inputs...))
				} else {
					// 敵勝利パス
					enemyPaths = append(enemyPaths, append([]int(nil), n.inputs...))
				}
				return
			}

			// 0.0〜1.0 を 0.1 刻み 11 通り用意し、毎ノードでシャッフル
			choices := make([]int, 10)
			for i := 0; i < 10; i++ {
				choices[i] = i
			}
			rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })

			// ランダムに mctsWidth 本を採用（幅制限）
			if len(choices) > mctsWidth {
				choices = choices[:mctsWidth]
			}

			for _, v := range choices {
				dfs(node{depth: n.depth + 1, inputs: append(append([]int(nil), n.inputs...), v)})
				if nodes >= mctsMaxNodes {
					return
				}
			}
		}

		dfs(node{depth: 0, inputs: []int{}})

		// 双方に少なくとも 1 パスずつあれば OK
		if len(playerPaths) == 0 || len(enemyPaths) == 0 {
			return false, nil, nil
		}

		// ランダムに 1 本ずつ返す
		playerPath := playerPaths[rand.Intn(len(playerPaths))]
		enemyPath := enemyPaths[rand.Intn(len(enemyPaths))]
		return true, playerPath, enemyPath
	}

	// ——— メインループ ————————————————————————————
	var debugSearchSeedCount int
	for try := 0; try < maxTries; try++ {
		seedCandidate := logic.NewSeedManager().GetSeed()
		rule := domain.NewRuleMatrix(seedCandidate, size)

		// RoughFilter
		playerWins, n := simulateSamples(rule, roughSamples)
		low, high := betaCI(playerWins, n)
		if !(low < 0.99 && high > 0.01) { // ほぼ 0 でも 1 でもない
			debugSearchSeedCount++
			continue
		}

		// DeepFilter
		playerWins, n = simulateSamples(rule, deepSamples)
		pHat := float64(playerWins) / float64(n)
		if pHat == 0 || pHat == 1 {
			debugSearchSeedCount++
			continue
		}

		// ProofPhase
		ok, playerPath, enemyPath := proofPhase(rule)
		debugSearchSeedCount++
		fmt.Printf("[Proof] Seed %d (試行 %d): ok=%v\n", seedCandidate, debugSearchSeedCount, ok)

		if ok {
			fmt.Printf("Found valid seed: %d with playerPath=%v, enemyPath=%v\n", seedCandidate, playerPath, enemyPath)
			return seedCandidate, playerPath, enemyPath, nil
		}

		// 進行が遅いときのデバッグ用出力（任意）
		if debugSearchSeedCount%100 == 0 {
			fmt.Printf("Tried %d seeds so far, still searching...\n", debugSearchSeedCount)
		}
	}
	return 0, nil, nil, fmt.Errorf("valid seed not found after %d tries", maxTries)
}
