package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
)

type scene struct {
	time       int
	background *sdl.Texture
	bird       *bird
	pipes      *pipes
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "./resources/background.jpg")
	if err != nil {
		return nil, fmt.Errorf("could not create scene: %v", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, err
	}

	pipes, err := newPipes(r)
	if err != nil {
		return nil, err
	}

	return &scene{background: bg, bird: bird, pipes: pipes}, nil
}
func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++

	r.Clear()

	err := r.Copy(s.background, nil, nil)
	if err != nil {
		return fmt.Errorf("could not paint scene: %v", err)
	}

	err = s.bird.paint(r)
	if err != nil {
		fmt.Errorf("couldn paint the bird: %v", err)
	}

	err = s.pipes.paint(r)
	if err != nil {
		fmt.Errorf("couldn paint the pip: %v", err)
	}

	r.Present()

	return nil
}

func (s *scene) destroy() {
	s.background.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
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
				if s.bird.isDead() {
					drawTitle(r, "Game over")
					time.Sleep(2 * time.Second)
					s.restart()
				}
				err := s.paint(r)
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
		s.bird.jump()
		return false
	default:
		log.Printf("unkown event: %T", event)
		return false
	}
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}
