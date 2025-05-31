package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"axiom_shift/internal/game"
)

func main() {
	g := game.NewGame()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Axiom Shift")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
