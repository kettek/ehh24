package res

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/png" // I am justifying this as the linter so demands. Look at this justifying, it's unbelievable.
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/stax"
)

//go:embed *.png
//go:embed places/*.json
//go:embed places/*.txt
var f embed.FS

// StaxImage is a convenience struct that stores Stax and ebiten.Image goodies.
type StaxImage struct {
	EbiImage *ebiten.Image
	Stax     stax.Stax
}

// Staxii is a cache of our stax images.
var Staxii map[string]StaxImage = make(map[string]StaxImage)

// Images is a cache of our non-stax images.
var Images map[string]*ebiten.Image = make(map[string]*ebiten.Image)

// Places is a cache of our places.
var Places map[string]Place = make(map[string]Place)

// Scripts is a cache of place scripts.
var Scripts map[string]string = make(map[string]string)

// GetStax gets the StaxImage associated with the given name, if possible.
func GetStax(name string) (StaxImage, error) {
	st, ok := Staxii[name]
	if !ok {
		return StaxImage{}, errors.New("Stax not found")
	}
	return st, nil
}

// ReadAssets reads all images, staxii, and places from disk.
func ReadAssets() error {
	entries, err := ReadDirs(".", "")
	if err != nil {
		return err
	}
	for _, e := range entries {
		if strings.HasSuffix(e, ".png") {
			data, err := ReadFile(e)
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
					EbiImage: eimg,
				}
			}
		} else if strings.HasSuffix(e, ".json") {
			data, err := ReadFile(e)
			if err != nil {
				return err
			}
			var place Place
			if err := json.Unmarshal(data, &place); err != nil {
				return err
			}
			Places[e[:len(e)-len(".json")]] = place
		} else if strings.HasSuffix(e, ".txt") {
			data, err := ReadFile(e)
			if err != nil {
				return err
			}
			Scripts[e[:len(e)-len(".txt")]] = string(data)
		}
	}
	return nil
}

// RefreshAssets clears the asset caches and reloads all assets.
func RefreshAssets() error {
	Staxii = make(map[string]StaxImage)
	Images = make(map[string]*ebiten.Image)
	Places = make(map[string]Place)
	Scripts = make(map[string]string)
	return ReadAssets()
}

// ReadDirs recursively reads directory contents from both embedded and on-disk locations.
func ReadDirs(name string, prepend string) ([]string, error) {
	entries, err := ReadDir(name)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			n2, err := ReadDirs(e.Name(), e.Name()+"/")
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
