package main

import (
	"embed"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"time"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/clktmr/n64/drivers/controller"
	n64draw "github.com/clktmr/n64/drivers/draw"
	"github.com/clktmr/n64/rcp/texture"
)

//go:embed assets/road-tile.png
var _roadPng embed.FS
var roadPng cartfs.FS = cartfs.Embed(_roadPng)

type Road struct {
	Sprite

	coinTimer time.Duration
}

func NewRoad(children ...Updater) *Road {
	r, err := roadPng.Open("assets/road-tile.png")
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(r)
	if err != nil {
		panic(err)
	}
	tex := texture.NewTextureFromImage(img)
	road := &Road{
		Sprite: *NewSprite(tex, 1, 1, 0),
	}
	road.relativePos.Y += worldbounds.Dy() - road.frames[0].Dy()

	for range 5 {
		road.AddChild(NewCoin())
	}
	return road
}

func (p *Road) Update(delta time.Duration, input [4]controller.Controller) {
	p.relativePos.X -= int(delta.Milliseconds() >> 2)

	p.coinTimer += delta
	if p.coinTimer >= 250*time.Millisecond {
		p.coinTimer = 0
		if coin := p.getFreeCoin(); coin != nil {
			coin.relativePos.X = -p.globalPos.X + worldbounds.Dx() - 1
			coin.relativePos.Y = rand.Intn(159) + 34 - coin.Size().Y
			coin.hidden = false
		}
	}
	p.Sprite.Update(delta, input)
}

func (p *Road) getFreeCoin() *Coin {
	for _, coin := range p.Children() {
		if coin, ok := coin.(*Coin); ok {
			if coin.hidden {
				return coin
			}
		}
	}
	return nil
}

func (p *Road) Render(dst draw.Image) {
	frame := p.frames[p.frame]
	r := image.Rectangle{Max: frame.Size()}
	r = r.Add(image.Point{p.globalPos.X % frame.Dx(), p.globalPos.Y})
	for {
		n64draw.Draw(dst, r, p.sheet, frame.Bounds().Min, draw.Src)
		r = r.Add(image.Point{r.Dx(), 0})
		if !r.Overlaps(dst.Bounds()) {
			break
		}
	}
}

func (p *Road) Z() int {
	return -1000
}
