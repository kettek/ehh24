package game

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return ow, oh
}
