package main

import (
	"embed"
	"image/png"
	"math"
	"time"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/serial/joybus"
	"github.com/clktmr/n64/rcp/texture"
)

var (
	//go:embed assets/menu-animation.png assets/main-menu-buttons.png
	_menuPngs embed.FS
	menuPngs  cartfs.FS = cartfs.Embed(_menuPngs)
)

type Title struct {
	Sprite

	button *Sprite
	seek   time.Duration

	next Updater
}

func NewTitle(next Updater) *Title {
	r, err := menuPngs.Open("assets/menu-animation.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		panic(err)
	}
	tex := texture.NewTextureFromImage(img)
	node := &Title{
		Sprite: *NewSprite(tex, 1, 2, 100*time.Millisecond),
		next:   next,
	}

	node.relativePos.X += worldbounds.Dx()/2 - node.Sprite.Size().X/2
	node.relativePos.Y += worldbounds.Dy()/2 - node.Sprite.Size().Y/2

	r, err = menuPngs.Open("assets/main-menu-buttons.png")
	if err != nil {
		panic(err)
	}
	img, err = png.Decode(r)
	if err != nil {
		panic(err)
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
