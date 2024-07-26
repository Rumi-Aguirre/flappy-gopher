package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const (
	gravity      = 1
	initialY     = 320
	initialSpeed = 8
	birdHeight   = 50
	birdWidth    = 43
	jumpSpeed    = -14
)

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
		path := fmt.Sprintf("./resources/Frame-%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not create bird textures: %v", err)
		}
		textures = append(textures, texture)
	}

	return &bird{
		time:     0,
		textures: textures,
		x:        40,
		y:        initialY,
		h:        birdHeight,
		w:        birdWidth,
		speed:    initialSpeed,
	}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	i := b.time % len(b.textures)
	rectangle := &sdl.Rect{X: b.x, Y: b.y, W: b.w, H: b.h}
	if err := r.Copy(b.textures[i], nil, rectangle); err != nil {
		return fmt.Errorf("could not paint bird: %v", err)
	}

	return nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.y += int32(b.speed)
	b.speed += gravity

	if b.y > screenHeight || b.y < 0 {
		b.dead = true
	}
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

	b.speed = jumpSpeed
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dead
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = initialY
	b.speed = initialSpeed
	b.dead = false
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if p.x > b.x+b.w || p.x+p.w < b.x {
		return
	}

	if p.inverted {
		if b.y < p.h {
			b.dead = true
		}
	} else {
		if b.y+b.h > screenHeight-p.h {
			b.dead = true
		}
	}
}
