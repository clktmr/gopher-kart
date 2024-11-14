package main

import (
	"embedded/rtos"
	"image"
	"image/color"
	"image/draw"
	"slices"
	"time"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/drivers/display"
	n64draw "github.com/clktmr/n64/drivers/draw"
)

type Updater interface {
	Update(delta time.Duration, input [4]controller.Controller)
	Children() []Updater
	AddChild(Updater)
	RemoveChild(Updater)
	setParent(Updater)
}

type Renderer interface {
	Render()
	Z() int
}

type GameLoop struct {
	root     Updater
	renderer *n64draw.Rdp
	display  *display.Display
}

func NewGameLoop(disp *display.Display, renderer *n64draw.Rdp, root Updater) *GameLoop {
	return &GameLoop{root, renderer, disp}
}

var ClearColor = color.RGBA{0xb9, 0xff, 0xfd, 0xff}

func (p *GameLoop) Run() {
	gamepad := make(chan [4]controller.Controller)
	last := rtos.Nanotime()
	go func() {
		for {
			controller.States.Poll()
			gamepad <- controller.States
		}
	}()
	clearImg := image.Uniform{ClearColor}
	renderNodes := make([]Renderer, 0, 64)
	for {
		// finish current frame
		write := p.display.Swap()

		// setup next frame
		now := rtos.Nanotime()
		input := <-gamepad
		p.renderer.SetFramebuffer(write)

		p.renderer.Draw(p.renderer.Bounds(), &clearImg, image.Point{}, draw.Src)

		p.root.Update(now-last, input)
		last = now

		renderNodes = appendChildren(renderNodes[:0], p.root)

		slices.SortFunc(renderNodes, func(a, b Renderer) int { return a.Z() - b.Z() })
		for _, node := range renderNodes {
			node.Render()
		}

		p.renderer.Flush()
	}
}

func appendChildren(nodes []Renderer, root Updater) []Renderer {
	for _, child := range root.Children() {
		if child, ok := child.(Renderer); ok {
			nodes = append(nodes, child)
		}
		nodes = appendChildren(nodes, child)
	}

	return nodes
}
