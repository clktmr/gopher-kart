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

//go:embed assets/coin-shadow.png
var coinPng string
var coinImg *image.NRGBA

type Coin struct {
	Sprite

	hidden bool
}

func init() {
	img, err := png.Decode(strings.NewReader(coinPng))
	if err != nil {
		panic(err)
	}

	var ok bool
	coinImg, ok = img.(*image.NRGBA)
	if !ok {
		panic("wrong image type")
	}
}

func NewCoin() *Coin {
	tex := texture.NewNRGBA32FromImage(coinImg)
	return &Coin{
		Sprite: *NewSprite(tex, 6, 1, 100*time.Millisecond),
	}
}

func (p *Coin) Update(delta time.Duration, input [4]controller.Controller) {
	p.Sprite.Update(delta, input)

	frame := p.frames[p.frame]
	r := image.Rectangle{Max: frame.Size()}.Add(p.globalPos)

	if !r.Overlaps(renderer.Bounds()) {
		p.hidden = true
	}
	for _, player := range players {
		if r.Overlaps(player.Bounds()) {
			p.hidden = true
		}
	}
}

func (p *Coin) Render() {
	if !p.hidden {
		p.Sprite.Render()
	}
}
