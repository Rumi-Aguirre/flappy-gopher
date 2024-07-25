package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const gravity = 1

type bird struct {
	mu         sync.RWMutex
	time       int
	textures   []*sdl.Texture
	x, y, h, w int32
	speed      float64
	dead       bool
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("./resources/Frame-%v.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not create scene: %v", err)
		}

		textures = append(textures, texture)
	}

	return &bird{time: 0, textures: textures, y: 320, speed: 8, x: 40, h: 50, w: 43}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	i := b.time % len(b.textures)
	rectangle := &sdl.Rect{b.x, b.y, b.h, b.w}
	err := r.Copy(b.textures[i], nil, rectangle)
	if err != nil {
		return fmt.Errorf("could not paint scene: %v", err)
	}

	return nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.time++
	b.y += int32(b.speed)
	if b.y > 600 || b.y < 0 {
		b.dead = true
	}

	b.speed += gravity
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.speed = -14
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dead
}

func (b *bird) restart() {
	b.mu.RLock()
	defer b.mu.RUnlock()

	b.y = 320
	b.speed = 8
	b.dead = false
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if p.x > b.x+b.w { // too far right
		return
	}

	if p.x+p.w < b.x { // too far left
		return
	}

	if p.inverted {
		if p.h < b.y-b.h/2 { // pipe is too low
			return
		}
	} else {
		if p.h > b.y-b.h/2 { // pipe is too low
			return
		}
	}

	b.dead = true
}
