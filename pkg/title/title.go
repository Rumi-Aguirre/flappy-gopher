package title

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	fontPath = "../resources/SEASRN__.ttf"
	fontSize = 64
)

var (
	textColor = sdl.Color{R: 111, G: 230, B: 16, A: 255}
)

type Title struct {
	font *ttf.Font
}

func NewTitle() (*Title, error) {
	font, err := ttf.OpenFont(fontPath, fontSize)
	if err != nil {
		return nil, fmt.Errorf("font cannot be opened: %v", err)
	}
	return &Title{font: font}, nil
}

func (t *Title) Close() {
	if t.font != nil {
		t.font.Close()
	}
}

func (t *Title) Paint(renderer *sdl.Renderer, text string) error {
	if err := renderer.Clear(); err != nil {
		return err
	}

	surface, err := t.font.RenderUTF8Solid(text, textColor)
	if err != nil {
		return fmt.Errorf("title cannot be rendered: %v", err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("surface cannot be created: %v", err)
	}
	defer texture.Destroy()

	if err := renderer.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("error copying texture: %v", err)
	}

	renderer.Present()
	return nil
}
