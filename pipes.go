package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
	"time"
)

const (
	pipeTexturePath = "./resources/pipe-3.png"
	initialPipeX    = 800
	pipeWidth       = 43
	minPipeHeight   = 220
	maxPipeHeight   = 340 // minPipeHeight + 120
	pipeSpeed       = 10
)

type pipes struct {
	mu      sync.RWMutex
	texture *sdl.Texture
	speed   int32
	pipes   []*pipe
}

type pipe struct {
	mu       sync.RWMutex
	texture  *sdl.Texture
	x        int32
	h        int32
	w        int32
	inverted bool
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(r, pipeTexturePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load pipe texture: %v", err)
	}

	ps := &pipes{
		texture: texture,
		speed:   pipeSpeed,
	}

	go ps.generatePipes(r)

	return ps, nil
}

func (ps *pipes) generatePipes(r *sdl.Renderer) {
	for {
		pipe, err := newPipe(r)
		if err != nil {
			fmt.Printf("error creating new pipe: %v\n", err)
			continue
		}

		ps.mu.Lock()
		ps.pipes = append(ps.pipes, pipe)
		ps.mu.Unlock()

		time.Sleep(1 * time.Second)
	}
}

func newPipe(r *sdl.Renderer) (*pipe, error) {
	texture, err := img.LoadTexture(r, pipeTexturePath)
	if err != nil {
		return nil, fmt.Errorf("cannot load pipe texture: %v", err)
	}

	return &pipe{
		texture:  texture,
		x:        initialPipeX,
		h:        minPipeHeight + int32(rand.Intn(maxPipeHeight-minPipeHeight)),
		w:        pipeWidth,
		inverted: rand.Float32() > 0.5,
	}, nil
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	y := screenHeight - p.h
	if p.inverted {
		y = 0
	}
	rectangle := &sdl.Rect{X: p.x, Y: y, W: p.w, H: p.h}

	if err := r.Copy(texture, nil, rectangle); err != nil {
		return fmt.Errorf("could not paint pipe: %v", err)
	}

	return nil
}

func (p *pipe) update(speed int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x -= speed
}

func (p *pipe) touch(b *bird) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	b.touch(p)
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(r, ps.texture); err != nil {
			return err
		}
	}

	return nil
}

func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var remainingPipes []*pipe

	for _, p := range ps.pipes {
		p.update(ps.speed)
		if p.x > -30 {
			remainingPipes = append(remainingPipes, p)
		}
	}

	ps.pipes = remainingPipes
}

func (ps *pipes) destroy() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.texture.Destroy()
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	for _, p := range ps.pipes {
		p.touch(b)
	}
}
