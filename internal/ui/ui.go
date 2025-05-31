package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type UI struct {
	battleLog []string
}

func NewUI() *UI {
	return &UI{
		battleLog: make([]string, 0),
	}
}

func (u *UI) Update() {
	// Update logic for the UI can be added here
}

func (u *UI) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// バトルログを上から下へ表示
	for i, log := range u.battleLog {
		y := 10 + i*20
		if y > 460 {
			break // 画面下部にはみ出さない
		}
		ebitenutil.DebugPrintAt(screen, log, 10, y)
	}
}

func (u *UI) AddBattleLog(log string) {
	u.battleLog = append(u.battleLog, log)
}

func (u *UI) ClearBattleLog() {
	u.battleLog = make([]string, 0)
}

// 任意座標にテキストを描画するユーティリティ
func DrawText(screen *ebiten.Image, msg string, x, y int) {
	ebitenutil.DebugPrintAt(screen, msg, x, y)
}
