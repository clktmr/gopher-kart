package main

import (
	_ "embed"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"strings"
	"time"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/texture"
)

//go:embed assets/road-tile.png
var roadPng string

type Road struct {
	Sprite

	coinTimer time.Duration
}

func NewRoad(children ...Updater) *Road {
	img, err := png.Decode(strings.NewReader(roadPng))
	if err != nil {
		panic(err)
	}
	imgRGBA, ok := img.(*image.NRGBA)
	if !ok {
		panic("wrong image type")
	}
	tex := texture.NewNRGBA32FromImage(imgRGBA)
	road := &Road{
		Sprite: *NewSprite(tex, 1, 1, 0),
	}
	road.relativePos.Y += renderer.Bounds().Dy() - road.frames[0].Dy()

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
			coin.relativePos.X = -p.globalPos.X + renderer.Bounds().Dx() - 1
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

func (p *Road) Render() {
	frame := p.frames[p.frame]
	r := image.Rectangle{Max: frame.Size()}
	r = r.Add(image.Point{p.globalPos.X % frame.Dx(), p.globalPos.Y})
	for {
		renderer.Draw(r, p.sheet, frame.Bounds().Min, draw.Src)
		r = r.Add(image.Point{r.Dx(), 0})
		if !r.Overlaps(renderer.Bounds()) {
			break
		}
	}
}

func (p *Road) Z() int {
	return -1000
}
