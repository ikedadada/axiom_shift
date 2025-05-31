package main

import (
	"log"

	"axiom_shift/internal/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Axiom Shift")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
