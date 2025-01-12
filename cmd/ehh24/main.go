package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ehh24/pkg/editor"
	"github.com/kettek/ehh24/pkg/game"
	"github.com/kettek/ehh24/pkg/intro"
	"github.com/kettek/ehh24/pkg/res"
	"github.com/kettek/ehh24/pkg/splash"
	"github.com/kettek/ehh24/pkg/statemachine"
)

func main() {
	if err := res.ReadAssets(); err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Hello, 世界")

	m := statemachine.NewMachine(game.NewState())
	m.AddCheck(func() {
		if inpututil.IsKeyJustReleased(ebiten.KeyF1) {
			m.SetState(game.NewState())
		} else if inpututil.IsKeyJustReleased(ebiten.KeyF2) {
			m.SetState(editor.NewState())
		} else if inpututil.IsKeyJustReleased(ebiten.KeyF3) {
			m.SetState(splash.NewState())
		} else if inpututil.IsKeyJustReleased(ebiten.KeyF4) {
			m.SetState(intro.NewState())
		}
	})

	if err := ebiten.RunGame(m); err != nil {
		panic(err)
	}
}
