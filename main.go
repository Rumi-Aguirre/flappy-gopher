package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"time"
)

const (
	screenHeight = 600
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		os.Exit(2)
	}
}

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("couldnt intialize SDL: %v", err)
	}

	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("couldnt initialize font: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("window cannot be created: %v", err)
	}
	defer w.Destroy()

	sdl.PumpEvents()
	sdl.PumpEvents()

	drawTitle(r, "Flappy Rumai")
	time.Sleep(1 * time.Second)

	scene, err := newScene(r)
	if err != nil {
		return fmt.Errorf("cannot create scene: %v", err)
	}
	defer scene.destroy()

	events := make(chan sdl.Event)
	errc := scene.run(events, r)

	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

func drawTitle(r *sdl.Renderer, text string) error {
	err := r.Clear()
	if err != nil {
		return err
	}

	f, err := ttf.OpenFont("./resources/SEASRN__.ttf", 18)
	if err != nil {
		return fmt.Errorf("font cannot be opened: %v", err)
	}
	defer f.Close()

	s, err := f.RenderUTF8Solid(text, sdl.Color{111, 230, 16, 255})
	if err != nil {
		return fmt.Errorf("title cannot be rendered: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("surface cannot be created: %v", err)
	}
	defer t.Destroy()

	r.Copy(t, nil, nil)

	r.Present()

	return nil
}
