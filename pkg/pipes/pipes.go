package pipes

import (
	"flappy/pkg/bird"
	"flappy/pkg/window"
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
	"time"
)

const (
	pipeTexturePath = "../resources/Pipe-3.png"
	initialPipeX    = 800
	pipeWidth       = 43
	minPipeHeight   = 220
	maxPipeHeight   = 340 // minPipeHeight + 120
	pipeSpeed       = 10
)

type Pipes struct {
	mu      sync.RWMutex
	texture *sdl.Texture
	speed   int32
	pipes   []*Pipe
}

type Pipe struct {
	mu       sync.RWMutex
	texture  *sdl.Texture
	X        int32
	H        int32
	W        int32
	Inverted bool
}

func NewPipes(renderer *sdl.Renderer) (*Pipes, error) {
	texture, err := img.LoadTexture(renderer, pipeTexturePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load Pipe texture: %v", err)
	}

	ps := &Pipes{
		texture: texture,
		speed:   pipeSpeed,
	}

	go ps.generatePipes(renderer)

	return ps, nil
}

func (ps *Pipes) generatePipes(renderer *sdl.Renderer) {
	for {
		pipe, err := NewPipe(renderer)
		if err != nil {
			fmt.Printf("error creating new Pipe: %v\n", err)
			continue
		}

		ps.mu.Lock()
		ps.pipes = append(ps.pipes, pipe)
		ps.mu.Unlock()

		time.Sleep(1 * time.Second)
	}
}

func NewPipe(renderer *sdl.Renderer) (*Pipe, error) {
	texture, err := img.LoadTexture(renderer, pipeTexturePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load Pipe texture: %v", err)
	}

	return &Pipe{
		texture:  texture,
		X:        initialPipeX,
		H:        minPipeHeight + int32(rand.Intn(maxPipeHeight-minPipeHeight)),
		W:        pipeWidth,
		Inverted: rand.Float32() > 0.5,
	}, nil
}

func (p *Pipe) paint(renderer *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	y := window.Height - p.H
	if p.Inverted {
		y = 0
	}
	rectangle := &sdl.Rect{X: p.X, Y: y, W: p.W, H: p.H}

	if err := renderer.Copy(texture, nil, rectangle); err != nil {
		return fmt.Errorf("could not paint Pipe: %v", err)
	}

	return nil
}

func (p *Pipe) update(speed int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.X -= speed
}

func (p *Pipe) touch(bird *bird.Bird) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	bird.Mu.Lock()
	defer bird.Mu.Unlock()

	if p.X > bird.X+bird.W || p.X+p.W < bird.X {
		return
	}

	if p.Inverted {
		if bird.Y < p.H {
			bird.Dead = true
		}
	} else {
		if bird.Y+bird.H > window.Height-p.H {
			bird.Dead = true
		}
	}
}

func (ps *Pipes) Paint(renderer *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(renderer, ps.texture); err != nil {
			return err
		}
	}

	return nil
}

func (ps *Pipes) Update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var remainingPipes []*Pipe

	for _, p := range ps.pipes {
		p.update(ps.speed)
		if p.X > -30 {
			remainingPipes = append(remainingPipes, p)
		}
	}

	ps.pipes = remainingPipes
}

func (ps *Pipes) Destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.texture.Destroy()
}

func (ps *Pipes) Restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *Pipes) Touch(bird *bird.Bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		p.touch(bird)
	}
}
