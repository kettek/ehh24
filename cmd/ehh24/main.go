package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game"
	"github.com/kettek/ehh24/pkg/res"
)

func main() {
	if err := res.ReadAssets(); err != nil {
		panic(err)
	}

	game := game.NewGame()
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, 世界")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
