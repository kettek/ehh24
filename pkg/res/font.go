package res

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Font Is Our One And Only Font That Has Meaning Only To Me And Those Who Can Learn.
var Font *text.GoTextFaceSource

func init() {
	ff, err := ReadFile("nokore.ttf")
	if err != nil {
		panic(err)
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(ff))
	if err != nil {
		panic(err)
	}
	Font = s
}
