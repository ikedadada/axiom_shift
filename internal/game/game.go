package game

import (
	"axiom_shift/internal/logic"
	"axiom_shift/internal/ui"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	battleCount int
	battleMax   int
	player      *Character
	enemy       *Enemy
	rule        *logic.RuleMatrix
	ui          UIInterface
	inputValue  float64
	phase       string // "input", "confirm", "battle", "end"
	lastWin     bool   // 最終戦の勝敗記録
}

type UIInterface interface {
	AddBattleLog(string)
	ClearBattleLog()
	Draw(screen *ebiten.Image)
}

func NewGame() *Game {
	player := NewCharacter(NewMatrix(2, 2), 0.5)
	enemy := NewEnemy("Enemy", NewMatrix(2, 2), 0.5)
	rule := logic.NewRuleMatrix(42, 2)
	ui := ui.NewUI()
	ui.ClearBattleLog() // ゲーム開始時にログをクリア（ClearBattleLogの活用）
	return &Game{
		battleCount: 0,
		battleMax:   10, // バトル回数を10回に変更
		player:      player,
		enemy:       enemy,
		rule:        rule,
		ui:          ui,
		phase:       "input",
		lastWin:     false,
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
		service := NewBattleService(g.player, g.enemy, g.rule)
		result, win := service.ExecuteBattle(g.inputValue)
		g.battleCount++
		log := fmt.Sprintf("Battle %d: Input=%s Result=%s Win/Lose=%s", g.battleCount, formatFloat(g.inputValue), formatFloat(result), winLoseStrEN(win))
		g.ui.AddBattleLog(log)
		// プレイヤーが勝った場合のみ敵が成長
		if win {
			ruleMatrix := g.rule.GetMatrix()
			m := &Matrix{data: ruleMatrix, rows: len(ruleMatrix), cols: len(ruleMatrix[0])}
			// 敵がプレイヤーに勝てるまで最大10回Grow
			maxTry := 10
			for i := 0; i < maxTry; i++ {
				g.enemy.Grow(g.inputValue, m)
				// 成長後に再度バトル判定
				serviceTmp := NewBattleService(g.player, g.enemy, g.rule)
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480 // Set the game window size
}

func (g *Game) Reset() {
	g.battleCount = 0
	g.player = NewCharacter(NewMatrix(2, 2), 0.5)
	g.enemy = NewEnemy("Enemy", NewMatrix(2, 2), 0.5)
	// ルール行列は同じものを再利用
	g.ui.ClearBattleLog()
	g.phase = "input"
	g.lastWin = false
}

// winLoseStrEN: 英語表記の勝敗判定は今後も使うため残す
func winLoseStrEN(win bool) string {
	if win {
		return "WIN"
	}
	return "LOSE"
}
