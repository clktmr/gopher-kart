package main

import (
	_ "embed"
	"image"
	"image/png"
	"strings"
	"time"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/texture"
)

type playerVariant string

var (
	//go:embed assets/burgundy-gopher.png
	Burgundy playerVariant

	//go:embed assets/beige-gopher.png
	Beige playerVariant

	//go:embed assets/black-gopher.png
	Black playerVariant
)

type Player struct {
	Sprite

	hidden     bool
	controller int
}

func NewPlayer(v playerVariant, controller int) *Player {
	img, err := png.Decode(strings.NewReader(string(v)))
	if err != nil {
		panic(err)
	}
	imgRGBA, ok := img.(*image.NRGBA)
	if !ok {
		panic("wrong image type")
	}
	tex := texture.NewTextureFromImage(imgRGBA)
	player := &Player{
		Sprite:     *NewSprite(tex, 2, 1, 100*time.Millisecond),
		controller: controller,
	}
	player.relativePos.Y += renderer.Bounds().Dy() / 2
	return player
}

func (p *Player) Update(delta time.Duration, input [4]controller.Controller) {
	p.Sprite.Update(delta, input)

	stick := image.Point{int(input[p.controller].X()), -int(input[p.controller].Y())}
	stick = stick.Mul(int(delta.Microseconds()))
	stick = stick.Div(5e5)

	p.relativePos = p.relativePos.Add(stick)

	p.hidden = !input[p.controller].Present()
}

func (p *Player) Render() {
	if p.hidden {
		return
	}
	p.Sprite.Render()
}
