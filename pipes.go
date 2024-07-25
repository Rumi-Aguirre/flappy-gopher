package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
	"time"
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
	texture, err := img.LoadTexture(r, "./resources/pipe-3.png")
	if err != nil {
		fmt.Errorf("canno load pipe texture: %v", err)
	}

	ps := &pipes{
		texture: texture,
		speed:   10,
	}

	go func() {
		for {
			ps.mu.Lock()
			pipe, _ := newPipe(r)
			ps.pipes = append(ps.pipes, pipe)
			ps.mu.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()

	return ps, nil
}

func newPipe(r *sdl.Renderer) (*pipe, error) {
	texture, err := img.LoadTexture(r, "./resources/pipe-3.png")
	if err != nil {
		fmt.Errorf("canno load pipe texture: %v", err)
	}

	return &pipe{
		texture:  texture,
		x:        800,
		h:        220 + int32(rand.Intn(120)),
		w:        43,
		inverted: rand.Float32() > 0.5,
	}, nil
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rectangle := &sdl.Rect{p.x, 600 - p.h, p.w, p.h}

	if p.inverted {
		rectangle.Y = 0
	}

	err := r.Copy(texture, nil, rectangle)
	if err != nil {
		return fmt.Errorf("could not paint scene: %v", err)
	}

	return nil
}

func (p *pipe) update(speed int32) {
	p.mu.RLock()
	defer p.mu.RUnlock()

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

	var rem []*pipe

	for _, p := range ps.pipes {
		p.update(ps.speed)
		if p.x > -30 {
			rem = append(rem, p)
		}
	}

	ps.pipes = rem
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
