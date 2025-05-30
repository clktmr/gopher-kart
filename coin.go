package main

import (
	"embed"
	"image"
	"image/png"
	"math/rand"
	"time"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/texture"
)

//go:embed assets/coin-shadow.png
var _coinPng embed.FS
var coinPng cartfs.FS = cartfs.Embed(_coinPng)
var coinImg image.Image

type Coin struct {
	Sprite

	hidden bool
}

func init() {
	r, err := coinPng.Open("assets/coin-shadow.png")
	if err != nil {
		panic(err)
	}
	coinImg, err = png.Decode(r)
	if err != nil {
		panic(err)
	}
}

func NewCoin() *Coin {
	tex := texture.NewTextureFromImage(coinImg)
	coin := &Coin{
		Sprite: *NewSprite(tex, 6, 1, 100*time.Millisecond),
	}
	coin.frame = rand.Intn(len(coin.frames))
	return coin
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
