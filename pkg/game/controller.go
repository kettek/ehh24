package game

type Controller interface {
	Update(ctx *DrawContext, t *Thinger) []Action
}

type PlayerController struct {
}

func NewPlayerController() *PlayerController {
	return &PlayerController{}
}

func (p *PlayerController) Update(ctx *DrawContext, t *Thinger) (a []Action) {
	x, y := ctx.MousePosition()
	w, h := ctx.Size()

	a = append(a, &ActionLook{
		LookX: (x - t.X) / w * 4,
		LookY: (y - t.Y) / h * 4,
	})

	return a
}

type CursorController struct {
}

func NewCursorController() *CursorController {
	return &CursorController{}
}

func (c *CursorController) Update(ctx *DrawContext, t *Thinger) (a []Action) {
	x, y := ctx.MousePosition()

	a = append(a, &ActionPosition{
		X: x,
		Y: y,
	})

	return a
}
