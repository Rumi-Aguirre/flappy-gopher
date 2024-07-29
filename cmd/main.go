package main

import (
	"flappy/pkg/scene"
	"flappy/pkg/window"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"time"
)

const (
	gameTitle = "Flappy Gopher"
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

	w, r, err := sdl.CreateWindowAndRenderer(window.Width, window.Height, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("window cannot be created: %v", err)
	}
	defer w.Destroy()

	sdl.PumpEvents()
	sdl.PumpEvents()

	scene, err := scene.NewScene(r)
	if err != nil {
		return fmt.Errorf("cannot create scene: %v", err)
	}
	scene.DrawTitle(r, gameTitle)

	time.Sleep(1 * time.Second)

	defer scene.Destroy()

	events := make(chan sdl.Event)
	errc := scene.Run(events, r)

	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}
