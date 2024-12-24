package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game"
)

func main() {
	game := game.NewGame()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, 世界")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
