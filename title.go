package main

import (
	_ "embed"
	"image"
	"image/png"
	"math"
	"strings"
	"time"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/serial/joybus"
	"github.com/clktmr/n64/rcp/texture"
)

var (
	//go:embed assets/menu-animation.png
	titlePng string

	//go:embed assets/main-menu-buttons.png
	pressStartPng string
)

type Title struct {
	Sprite

	button *Sprite
	seek   time.Duration

	next Updater
}

func NewTitle(next Updater) *Title {
	img, err := png.Decode(strings.NewReader(titlePng))
	if err != nil {
		panic(err)
	}
	imgRGBA, ok := img.(*image.NRGBA)
	if !ok {
		panic("wrong image type")
	}
	tex := texture.NewTextureFromImage(imgRGBA)
	node := &Title{
		Sprite: *NewSprite(tex, 1, 2, 100*time.Millisecond),
		next:   next,
	}

	node.relativePos.X += renderer.Bounds().Dx()/2 - node.Sprite.Size().X/2
	node.relativePos.Y += renderer.Bounds().Dy()/2 - node.Sprite.Size().Y/2

	img, err = png.Decode(strings.NewReader(pressStartPng))
	if err != nil {
		panic(err)
	}
	imgRGBA, ok = img.(*image.NRGBA)
	if !ok {
		panic("wrong image type")
	}
	tex = texture.NewTextureFromImage(img)
	node.button = NewSprite(tex, 1, 2, 1*time.Second)
	node.button.relativePos.Y = node.Sprite.Size().Y + 5
	node.button.relativePos.X += node.Sprite.Size().X/2 - node.button.Size().X/2
	node.AddChild(node.button)
	return node
}

func (p *Title) Update(delta time.Duration, input [4]controller.Controller) {
	p.Sprite.Update(delta, input)

	p.seek += delta
	p.button.relativePos.Y = p.Sprite.Size().Y + 5 - int(math.Abs(math.Sin(float64(p.seek)/5e8)*16))

	if input[0].Pressed()&joybus.ButtonStart != 0 {
		p.parent.AddChild(p.next)
		p.parent.RemoveChild(p)
	}
}
