package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type InputManager struct {
	inputValue float64
}

func NewInputManager() *InputManager {
	return &InputManager{
		inputValue: 0.5, // Default value in the range [0, 1]
	}
}

func (im *InputManager) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		im.inputValue += 0.01
		if im.inputValue > 1 {
			im.inputValue = 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		im.inputValue -= 0.01
		if im.inputValue < 0 {
			im.inputValue = 0
		}
	}
}

func (im *InputManager) GetInputValue() float64 {
	return im.inputValue
}

func (im *InputManager) Draw(screen *ebiten.Image) {
	text := "Input Value: " + fmt.Sprintf("%.2f", im.inputValue)
	ebitenutil.DebugPrint(screen, text)
}
