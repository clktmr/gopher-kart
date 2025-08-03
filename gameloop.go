package main

import (
	"embedded/rtos"
	"image"
	"image/color"
	"image/draw"
	"io"
	"slices"
	"time"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/drivers/display"
	n64draw "github.com/clktmr/n64/drivers/draw"
	"github.com/clktmr/n64/drivers/rspq/mixer"
	"github.com/clktmr/n64/rcp/audio"
)

type Updater interface {
	Update(delta time.Duration, input [4]controller.Controller)
	Children() []Updater
	AddChild(Updater)
	RemoveChild(Updater)
	setParent(Updater)
}

type Renderer interface {
	Render(dst draw.Image)
	Z() int
}

type GameLoop struct {
	root    Updater
	display *display.Display
}

func NewGameLoop(disp *display.Display, root Updater) *GameLoop {
	return &GameLoop{root, disp}
}

var ClearColor = color.RGBA{0xb9, 0xff, 0xfd, 0xff}

func (p *GameLoop) Run() {
	gamepad := make(chan [4]controller.Controller)
	last := rtos.Nanotime()
	go func() {
		inputs := [4]controller.Controller{}
		for {
			controller.Poll(&inputs)
			gamepad <- inputs
		}
	}()
	go func() { io.Copy(audio.Buffer, mixer.Output) }()
	clearImg := image.Uniform{ClearColor}
	renderNodes := make([]Renderer, 0, 64)
	for {
		// finish current frame
		write := p.display.Swap()

		// setup next frame
		now := rtos.Nanotime()
		input := <-gamepad

		n64draw.Src.Draw(write, write.Bounds(), &clearImg, image.Point{})

		p.root.Update(now-last, input)
		last = now

		renderNodes = appendChildren(renderNodes[:0], p.root)

		slices.SortFunc(renderNodes, func(a, b Renderer) int { return a.Z() - b.Z() })
		for _, node := range renderNodes {
			node.Render(write)
		}

		n64draw.Flush()
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
