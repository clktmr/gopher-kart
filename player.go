package main

import (
	"embed"
	"image"
	"image/png"
	"time"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/texture"
)

type playerVariant string

var (
	//go:embed assets/*-gopher.png
	_playerPngs embed.FS
	playerPngs  cartfs.FS = cartfs.Embed(_playerPngs)
)

const (
	Burgundy playerVariant = "burgundy"
	Beige    playerVariant = "beige"
	Black    playerVariant = "black"
)

type Player struct {
	Sprite

	hidden     bool
	controller int
}

func NewPlayer(v playerVariant, controller int) *Player {
	r, err := playerPngs.Open("assets/" + string(v) + "-gopher.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		panic(err)
	}
	tex := texture.NewTextureFromImage(img)
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
