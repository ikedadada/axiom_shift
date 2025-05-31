package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"axiom_shift/internal/game"
)

func main() {
	seed := time.Now().UnixNano()
	g := game.NewGame(seed)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Axiom Shift")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
