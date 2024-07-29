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

func NewPipes(r *sdl.Renderer) (*Pipes, error) {
	texture, err := img.LoadTexture(r, pipeTexturePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load Pipe texture: %v", err)
	}

	ps := &Pipes{
		texture: texture,
		speed:   pipeSpeed,
	}

	go ps.generatePipes(r)

	return ps, nil
}

func (ps *Pipes) generatePipes(r *sdl.Renderer) {
	for {
		pipe, err := NewPipe(r)
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

func NewPipe(r *sdl.Renderer) (*Pipe, error) {
	texture, err := img.LoadTexture(r, pipeTexturePath)
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

func (p *Pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	y := window.Height - p.H
	if p.Inverted {
		y = 0
	}
	rectangle := &sdl.Rect{X: p.X, Y: y, W: p.W, H: p.H}

	if err := r.Copy(texture, nil, rectangle); err != nil {
		return fmt.Errorf("could not paint Pipe: %v", err)
	}

	return nil
}

func (p *Pipe) update(speed int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.X -= speed
}

func (p *Pipe) touch(b *bird.Bird) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	//b.Touch(p)

	b.Mu.Lock()
	defer b.Mu.Unlock()

	if p.X > b.X+b.W || p.X+p.W < b.X {
		return
	}

	if p.Inverted {
		if b.Y < p.H {
			b.Dead = true
		}
	} else {
		if b.Y+b.H > window.Height-p.H {
			b.Dead = true
		}
	}
}

func (ps *Pipes) Paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
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

func (ps *Pipes) Touch(b *bird.Bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		p.touch(b)
	}
}
