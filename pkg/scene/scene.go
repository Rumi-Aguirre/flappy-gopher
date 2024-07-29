package scene

import (
	"flappy/pkg/bird"
	"flappy/pkg/pipes"
	"flappy/pkg/title"
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
)

const (
	gameOverText = "Game Over"
)

type scene struct {
	time       int
	background *sdl.Texture
	bird       *bird.Bird
	pipes      *pipes.Pipes
	title      *title.Title
}

func NewScene(renderer *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(renderer, "../resources/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not create scene: %v", err)
	}

	bird, err := bird.NewBird(renderer)
	if err != nil {
		return nil, err
	}

	pipes, err := pipes.NewPipes(renderer)
	if err != nil {
		return nil, err
	}

	title, err := title.NewTitle()
	if err != nil {
		return nil, err
	}

	return &scene{background: bg, bird: bird, pipes: pipes, title: title}, nil
}
func (s *scene) update() {
	s.bird.Update()
	s.pipes.Update()
	s.pipes.Touch(s.bird)
}

func (s *scene) DrawTitle(renderer *sdl.Renderer, text string) error {
	return s.title.Paint(renderer, text)
}

func (s *scene) paint(renderer *sdl.Renderer) error {
	s.time++

	renderer.Clear()

	err := renderer.Copy(s.background, nil, nil)
	if err != nil {
		return fmt.Errorf("could not paint scene: %v", err)
	}

	err = s.bird.Paint(renderer)
	if err != nil {
		fmt.Errorf("couldn paint the bird: %v", err)
	}

	err = s.pipes.Paint(renderer)
	if err != nil {
		fmt.Errorf("couldn paint the pip: %v", err)
	}

	renderer.Present()

	return nil
}

func (s *scene) Destroy() {
	s.background.Destroy()
	s.bird.Destroy()
	s.pipes.Destroy()
}

func (s *scene) Run(events <-chan sdl.Event, renderer *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)

		tick := time.Tick(40 * time.Millisecond)
		done := false
		for !done {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()
				if s.bird.IsDead() {
					s.DrawTitle(renderer, gameOverText)
					time.Sleep(2 * time.Second)
					s.restart()
				}
				err := s.paint(renderer)
				if err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.TextInputEvent:
		print(event)
		s.bird.Jump()
		return false
	default:
		log.Printf("unkown event: %T", event)
		return false
	}
}

func (s *scene) restart() {
	s.bird.Restart()
	s.pipes.Restart()
}
