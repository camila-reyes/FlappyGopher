package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const (
	gravity   = 8
	jumpSpeed = 20
)

type bird struct {
	mu sync.RWMutex

	time     int
	textures []*sdl.Texture

	x, y  int32
	w, h  int32
	speed float64
	dead  bool
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("./res/imgs/frame-%d.png", i)
		bird, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not draw background: %v", err)
		}
		textures = append(textures, bird)
	}
	return &bird{textures: textures, x: 10, y: 300, w: 50, h: 43}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity
}

func (b *bird) paint(r *sdl.Renderer) error {
	rect := &sdl.Rect{X: 10, Y: 600 - b.y - b.h/2, W: b.w, H: b.h}

	i := b.time % len(b.textures)
	if err := r.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}
	return nil
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) isDead() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.dead
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = -jumpSpeed
}

func (b *bird) touch(p *pipe) {

	b.mu.Lock()
	defer b.mu.Unlock()

	if p.x > b.x+b.w { // too far right
		return
	} else if p.x+p.w < b.x { // too far left
		return
	} else if !p.inverted && p.h < b.y-b.h/2 { // pipe too low
		return
	} else if p.inverted && 600-p.h > b.y+b.h/2 { // inverted pipe too high
		return
	}
	b.dead = true
}
