package bird

import (
	"flappy/pkg/window"
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

type Bird struct {
	Mu         sync.RWMutex
	time       int
	textures   []*sdl.Texture
	X, Y, H, W int32
	speed      float64
	Dead       bool
}

func NewBird(renderer *sdl.Renderer) (*Bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("../resources/Frame-%d.png", i)
		texture, err := img.LoadTexture(renderer, path)
		if err != nil {
			return nil, fmt.Errorf("could not create Bird textures: %v", err)
		}
		textures = append(textures, texture)
	}

	return &Bird{
		time:     0,
		textures: textures,
		X:        40,
		Y:        initialY,
		H:        birdHeight,
		W:        birdWidth,
		speed:    initialSpeed,
	}, nil
}

func (b *Bird) Paint(renderer *sdl.Renderer) error {
	b.Mu.RLock()
	defer b.Mu.RUnlock()

	i := b.time % len(b.textures)
	rectangle := &sdl.Rect{X: b.X, Y: b.Y, W: b.W, H: b.H}
	if err := renderer.Copy(b.textures[i], nil, rectangle); err != nil {
		return fmt.Errorf("could not paint Bird: %v", err)
	}

	return nil
}

func (b *Bird) Update() {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	b.time++
	b.Y += int32(b.speed)
	b.speed += gravity

	if b.Y > window.Height || b.Y < 0 {
		b.Dead = true
	}
}

func (b *Bird) Destroy() {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b *Bird) Jump() {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	b.speed = jumpSpeed
}

func (b *Bird) IsDead() bool {
	b.Mu.RLock()
	defer b.Mu.RUnlock()

	return b.Dead
}

func (b *Bird) Restart() {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	b.Y = initialY
	b.speed = initialSpeed
	b.Dead = false
}
