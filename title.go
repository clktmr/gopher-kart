package main

import (
	"embed"
	"image/png"
	"io"
	"math"
	"time"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/drivers/rspq/mixer"
	"github.com/clktmr/n64/rcp/serial/joybus"
	"github.com/clktmr/n64/rcp/texture"
)

var (
	//go:embed assets/menu-animation.png assets/main-menu-buttons.png
	//go:embed "assets/8bit Bossa.pcm_s16be"
	_menuPngs embed.FS
	menuFiles cartfs.FS = cartfs.Embed(_menuPngs)
)

var menuMusic *mixer.Source

type Title struct {
	Sprite

	button *Sprite
	seek   time.Duration

	next Updater
}

func NewTitle(next Updater) *Title {
	r, err := menuFiles.Open("assets/menu-animation.png")
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

	r, err = menuFiles.Open("assets/main-menu-buttons.png")
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

	r, err = menuFiles.Open("assets/8bit Bossa.pcm_s16be")
	if err != nil {
		panic(err)
	}
	menuMusic = mixer.NewSource(mixer.Loop(r.(io.ReadSeeker)), 8000)
	mixer.SetSource(0, menuMusic)

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
