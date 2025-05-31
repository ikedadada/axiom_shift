package game

import (
	"axiom_shift/internal/domain"
	"axiom_shift/internal/logic"
	"axiom_shift/internal/ui"
	"axiom_shift/internal/usecase"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	battleCount int
	battleMax   int
	player      *domain.Player
	enemy       *domain.Enemy
	rule        *logic.RuleMatrix
	ui          UIInterface
	inputValue  int
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

func NewGame() *Game {
	pm := domain.NewMatrix(3, 3)
	for i := 0; i < 3; i++ {
		pm.Data[i][i] = 2.0
	}
	player := domain.NewPlayer(pm, 0.5)
	enemyMat := domain.NewMatrix(3, 3)
	for i := 0; i < 3; i++ {
		enemyMat.Data[i][2-i] = 2.0
	}
	enemy := domain.NewEnemy("Enemy", enemyMat, 0.5)
	battleMax := 10
	initialSeed := time.Now().UnixNano()
	seed, rule, playerPath, enemyPath, err := usecase.FindValidSeed(battleMax, initialSeed, player, enemy)
	if err != nil {
		panic(fmt.Sprintf("Seed search failed: %v", err))
	}
	player.Reset()
	enemy.Reset()
	_ = playerPath
	_ = enemyPath
	ui := ui.NewUI()
	ui.ClearBattleLog()
	return &Game{
		battleCount: 0,
		battleMax:   battleMax,
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
				g.inputValue = i
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
		service := usecase.NewBattleService(g.player, g.enemy, logicToDomainRuleMatrix(g.rule))
		result, win := service.DoBattleTurn(float64(g.inputValue)/9, g.battleCount)
		g.lastResult = &result
		g.battleCount++
		log := fmt.Sprintf("Battle %d: Input=%d Result=%s Win/Lose=%s", g.battleCount, g.inputValue, formatFloat(result), winLoseStrEN(win))
		g.ui.AddBattleLog(log)
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
		ui.DrawText(screen, "Press 0-9 to input", 10, 460)
	case "confirm":
		g.ui.Draw(screen)
		msg := fmt.Sprintf("Input: %d  [Enter: OK / Backspace: Re-input]", g.inputValue)
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
	ui.DrawText(screen, seedMsg, 485, 460)
	// 画面中央下にResultバーを描画
	if g.lastResult != nil {
		drawResultBar(screen, *g.lastResult)
	}
	// --- Player/Enemy行列のビジュアライズ ---
	startX, startY := 200, 310 // 画面下部のテキストの上
	if g.player != nil && g.player.MatrixState != nil {
		mat := g.player.MatrixState
		cellSize := 18
		margin := 3
		for i := 0; i < mat.Rows; i++ {
			for j := 0; j < mat.Cols; j++ {
				v := mat.Data[i][j]
				if v < 0 {
					v = 0
				}
				if v > 1 {
					v = 1
				}
				clr := color.RGBA{0, 0, uint8(64 + 191*v), 255} // 青の濃さ
				drawRect(screen, float64(startX+j*(cellSize+margin)), float64(startY+i*(cellSize+margin)), float64(cellSize), float64(cellSize), clr)
			}
		}
		ui.DrawText(screen, "Player", startX, startY-18)
	}
	startX += 180
	if g.enemy != nil && g.enemy.MatrixState != nil {
		mat := g.enemy.MatrixState
		cellSize := 18
		margin := 3
		for i := 0; i < mat.Rows; i++ {
			for j := 0; j < mat.Cols; j++ {
				v := mat.Data[i][j]
				if v < 0 {
					v = 0
				}
				if v > 1 {
					v = 1
				}
				clr := color.RGBA{uint8(64 + 191*v), 0, 0, 255} // 赤の濃さ
				drawRect(screen, float64(startX+j*(cellSize+margin)), float64(startY+i*(cellSize+margin)), float64(cellSize), float64(cellSize), clr)
			}
		}
		ui.DrawText(screen, "Enemy", startX, startY-18)
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
	g.player.Reset()
	g.enemy.Reset()
	g.rule = logic.NewRuleMatrix(g.seed, 3) // seedを再利用
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
