package res

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/stax"
)

//go:embed *.png
//go:embed places/*.json
var f embed.FS

type StaxImage struct {
	Image    image.Image
	EbiImage *ebiten.Image
	Stax     stax.Stax
}

var Staxii map[string]StaxImage = make(map[string]StaxImage)

var Images map[string]*ebiten.Image = make(map[string]*ebiten.Image)

var Places map[string]Place = make(map[string]Place)

func GetStax(name string) (StaxImage, error) {
	st, ok := Staxii[name]
	if !ok {
		return StaxImage{}, errors.New("Stax not found")
	}
	return st, nil
}

func ReadAssets() error {
	entries, err := ReadDir(".", "")
	if err != nil {
		return err
	}
	for _, e := range entries {
		if strings.HasSuffix(e, ".png") {
			data, err := f.ReadFile(e)
			if err != nil {
				return err
			}
			st, err := stax.ReadStaxFromPNG(data)
			if err != nil {
				fmt.Println(err, "for", e)
				//return err
			}
			png, _, err := image.Decode(strings.NewReader(string(data)))
			if err != nil {
				return err
			}
			eimg := ebiten.NewImageFromImage(png)
			if st == nil {
				Images[e[:len(e)-len(".png")]] = eimg
			} else {
				Staxii[e[:len(e)-len(".png")]] = StaxImage{
					Stax:     *st,
					Image:    png,
					EbiImage: eimg,
				}
			}
		} else if strings.HasSuffix(e, ".json") {
			data, err := f.ReadFile(e)
			if err != nil {
				return err
			}
			var place Place
			if err := json.Unmarshal(data, &place); err != nil {
				return err
			}
			Places[e[:len(e)-len(".json")]] = place
		}
	}
	return nil
}

func ReadDir(name string, prepend string) ([]string, error) {
	entries, err := f.ReadDir(name)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			n2, err := ReadDir(e.Name(), e.Name()+"/")
			fmt.Println(n2, err)
			if err != nil {
				continue
			}
			names = append(names, n2...)
			continue
		}
		names = append(names, prepend+e.Name())
	}
	return names, nil
}
