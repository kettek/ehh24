package res

import (
	"embed"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/stax"
)

//go:embed *.png
var f embed.FS

type StaxImage struct {
	Image    image.Image
	EbiImage *ebiten.Image
	Stax     stax.Stax
}

var Staxii map[string]StaxImage = make(map[string]StaxImage)

var Images map[string]*ebiten.Image = make(map[string]*ebiten.Image)

func GetStax(name string) (StaxImage, error) {
	st, ok := Staxii[name]
	if !ok {
		return StaxImage{}, errors.New("Stax not found")
	}
	return st, nil
}

func ReadAssets() error {
	entries, err := f.ReadDir(".")
	if err != nil {
		return err
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".png") {
			data, err := f.ReadFile(e.Name())
			if err != nil {
				return err
			}
			st, err := stax.ReadStaxFromPNG(data)
			if err != nil {
				fmt.Println(err, "for", e.Name())
				//return err
			}
			png, _, err := image.Decode(strings.NewReader(string(data)))
			if err != nil {
				return err
			}
			eimg := ebiten.NewImageFromImage(png)
			if st == nil {
				Images[e.Name()[:len(e.Name())-len(".png")]] = eimg
			} else {
				Staxii[e.Name()[:len(e.Name())-len(".png")]] = StaxImage{
					Stax:     *st,
					Image:    png,
					EbiImage: eimg,
				}
			}
		}
	}
	return nil
}
