package game

import (
	"axiom_shift/internal/domain"
	"axiom_shift/internal/logic"
	"axiom_shift/internal/ui"
	"axiom_shift/internal/usecase"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	battleCount int
	battleMax   int
	player      *domain.Player
	enemy       *domain.Enemy
	rule        *logic.RuleMatrix
	ui          UIInterface
	inputValue  float64
	phase       string   // "input", "confirm", "battle", "end"
	lastWin     bool     // 最終戦の勝敗記録
	seed        int64    // ルール生成用シード値
	lastResult  *float64 // 直近バトルの結果値（-1.0〜+1.0想定）
}

type UIInterface interface {
	AddBattleLog(string)
	ClearBattleLog()
	Draw(screen *ebiten.Image)
}

func logicToDomainRuleMatrix(lr *logic.RuleMatrix) *domain.RuleMatrix {
	mat := lr.GetMatrix()
	return &domain.RuleMatrix{Matrix: mat}
}

func NewGame(seed int64) *Game {
	// プレイヤーの初期行列を有利に設定（対角成分を大きく）
	pm := domain.NewMatrix(2, 2)
	pm.Data[0][0] = 2.0
	pm.Data[1][1] = 2.0
	player := domain.NewPlayer(pm, 0.5)
	// 敵は従来通り
	enemy := domain.NewEnemy("Enemy", domain.NewMatrix(2, 2), 0.5)
	rule := logic.NewRuleMatrix(seed, 2)
	ui := ui.NewUI()
	ui.ClearBattleLog()
	return &Game{
		battleCount: 0,
		battleMax:   10,
		player:      player,
		enemy:       enemy,
		rule:        rule,
		ui:          ui,
		phase:       "input",
		lastWin:     false,
		seed:        seed,
	}
}

// formatFloat: 全ての数値出力を統一的に整形できるメリットがあるため利用
func (g *Game) Update() error {
	switch g.phase {
	case "input":
		// キー入力受付: 0-9キーで0.0-1.0にマッピング
		for i := 0; i <= 9; i++ {
			if ebiten.IsKeyPressed(ebiten.Key0 + ebiten.Key(i)) {
				g.inputValue = float64(i) / 9.0
				g.phase = "confirm"
				break
			}
		}
	case "confirm":
		// Enterでバトルへ、Backspaceで再入力
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.phase = "battle"
		} else if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			g.phase = "input"
		}
	case "battle":
		// バトル回数に応じて成長率を増加
		g.player.GrowthRate = 0.5 + 0.1*float64(g.battleCount)
		service := usecase.NewBattleService(g.player, g.enemy, logicToDomainRuleMatrix(g.rule))
		result, win := service.ExecuteBattle(g.inputValue)
		g.lastResult = &result // 直近のバトル結果を保存
		g.battleCount++
		log := fmt.Sprintf("Battle %d: Input=%s Result=%s Win/Lose=%s", g.battleCount, formatFloat(g.inputValue), formatFloat(result), winLoseStrEN(win))
		g.ui.AddBattleLog(log)
		// プレイヤーが勝った場合のみ敵が成長
		if win {
			ruleMatrix := g.rule.GetMatrix()
			m := domain.NewMatrix(len(ruleMatrix), len(ruleMatrix[0]))
			for i := range ruleMatrix {
				for j := range ruleMatrix[i] {
					m.Data[i][j] = ruleMatrix[i][j]
				}
			}
			// 敵がプレイヤーに勝てるまで最大10回Grow
			maxTry := 10
			for i := 0; i < maxTry; i++ {
				g.enemy.Grow(g.inputValue, m)
				// 成長後に再度バトル判定
				serviceTmp := usecase.NewBattleService(g.player, g.enemy, logicToDomainRuleMatrix(g.rule))
				_, winTmp := serviceTmp.ExecuteBattle(g.inputValue)
				if winTmp {
					continue // まだ勝てない→さらにGrow
				} else {
					break // 勝てなくなったら終了
				}
			}
		}
		if g.battleCount >= g.battleMax {
			g.phase = "end"
			g.ui.AddBattleLog("---")
			g.lastWin = win
			if win {
				g.ui.AddBattleLog("[GAME WIN] Congratulations!")
			} else {
				g.ui.AddBattleLog("[GAME LOSE] Try again!")
			}
		} else {
			g.phase = "input"
		}
	case "end":
		// Rキーでリトライ
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.Reset()
		}
	}
	return nil
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.phase {
	case "input":
		g.ui.Draw(screen)
		// 指示文を画面下部に表示
		ui.DrawText(screen, "Press 0-9 to input (0.0-1.0)", 10, 460)
	case "confirm":
		g.ui.Draw(screen)
		msg := fmt.Sprintf("Input: %s  [Enter: OK / Backspace: Re-input]", formatFloat(g.inputValue))
		ui.DrawText(screen, msg, 10, 460)
	case "battle":
		g.ui.Draw(screen)
		ui.DrawText(screen, "Battle processing...", 10, 460)
	case "end":
		g.ui.Draw(screen)
		if g.lastWin {
			ui.DrawText(screen, "GAME WIN! Press R to retry same rule.", 10, 460)
		} else {
			ui.DrawText(screen, "GAME LOSE! Press R to retry same rule.", 10, 460)
		}
	}
	// 画面右下にSeed値を表示
	seedMsg := fmt.Sprintf("Seed: %d", g.seed)
	ui.DrawText(screen, seedMsg, 540, 460)
	// 画面中央下にResultバーを描画
	if g.lastResult != nil {
		drawResultBar(screen, *g.lastResult)
	}
}

// ebitenutil.DrawRectの代替
func drawRect(screen *ebiten.Image, x, y, w, h float64, clr color.Color) {
	img := ebiten.NewImage(int(w), int(h))
	img.Fill(clr)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(img, op)
}

// 結果値をバーでビジュアライズ
func drawResultBar(screen *ebiten.Image, result float64) {
	barY := 410
	barX := 120
	barW := 400
	barH := 18
	// 背景バー
	drawRect(screen, float64(barX), float64(barY), float64(barW), float64(barH), color.RGBA{80, 80, 80, 255})
	// 中央線
	drawRect(screen, float64(barX+barW/2-1), float64(barY), 2, float64(barH), color.White)
	// 結果バー
	clamped := result
	if clamped < -1 {
		clamped = -1
	}
	if clamped > 1 {
		clamped = 1
	}
	barLen := int(float64(barW/2) * clamped)
	if barLen > 0 {
		// プレイヤー有利（右側・緑）
		drawRect(screen, float64(barX+barW/2), float64(barY), float64(barLen), float64(barH), color.RGBA{80, 200, 80, 255})
	} else if barLen < 0 {
		// 敵有利（左側・赤）
		drawRect(screen, float64(barX+barW/2+barLen), float64(barY), float64(-barLen), float64(barH), color.RGBA{200, 80, 80, 255})
	}
	// ラベル
	ui.DrawText(screen, "Result", barX-70, barY+2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480 // Set the game window size
}

func (g *Game) Reset() {
	g.battleCount = 0
	// プレイヤーの初期行列も有利に再設定
	pm := domain.NewMatrix(2, 2)
	pm.Data[0][0] = 2.0
	pm.Data[1][1] = 2.0
	g.player = domain.NewPlayer(pm, 0.5)
	g.enemy = domain.NewEnemy("Enemy", domain.NewMatrix(2, 2), 0.5)
	g.rule = logic.NewRuleMatrix(g.seed, 2) // seedを再利用
	g.ui.ClearBattleLog()
	g.phase = "input"
	g.lastWin = false
	g.lastResult = nil
}

// winLoseStrEN: 英語表記の勝敗判定は今後も使うため残す
func winLoseStrEN(win bool) string {
	if win {
		return "WIN"
	}
	return "LOSE"
}
