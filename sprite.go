package main

import (
	"image"
	"image/draw"
	"time"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/rcp/texture"
)

type Sprite struct {
	Node

	sheet          *texture.Texture
	frames         []image.Rectangle
	frame          int
	seek, duration time.Duration
}

func NewSprite(sheet *texture.Texture, xCnt, yCnt int, duration time.Duration) *Sprite {
	bounds := sheet.Bounds()
	size := image.Point{
		bounds.Size().X / xCnt,
		bounds.Size().Y / yCnt,
	}

	frames := make([]image.Rectangle, 0)
	frame := image.Rectangle{Max: size}.Add(bounds.Min)
	for range yCnt {
		for range xCnt {
			frames = append(frames, frame)
			frame = frame.Add(image.Point{size.X, 0})
		}
		frame = frame.Add(image.Point{-frame.Min.X, size.Y})
	}

	return &Sprite{
		sheet:    sheet,
		frames:   frames,
		duration: duration,
	}
}

func (p *Sprite) Size() image.Point {
	return p.frames[p.frame].Size()
}

func (p *Sprite) Bounds() image.Rectangle {
	return image.Rectangle{Max: p.Size()}.Add(p.globalPos)
}

func (p *Sprite) Update(delta time.Duration, input [4]controller.Controller) {
	p.Node.Update(delta, input)

	p.seek += delta
	if p.seek > p.duration {
		p.frame = (p.frame + 1) % len(p.frames)
		p.seek -= p.duration
	}
}

func (p *Sprite) Render() {
	frame := p.frames[p.frame]
	r := p.Bounds()
	renderer.Draw(r, p.sheet, frame.Min, draw.Over)
}

func (p *Sprite) Z() int {
	frame := p.frames[p.frame]
	return p.globalPos.Y + frame.Dy()
}
