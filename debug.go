package main

import (
	"image/color"
	"runtime"
	"strconv"
	"time"

	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/drivers/display"
	"github.com/clktmr/n64/fonts/gomono12"
	"github.com/clktmr/n64/rcp/serial/joybus"
)

type Debug struct {
	Node

	hidden bool
	disp   *display.Display
	stats  []byte
	alloc  uint64
}

func NewDebug(disp *display.Display) *Debug {
	dbg := &Debug{
		disp:   disp,
		hidden: true,
		stats:  make([]byte, 0, 64),
	}
	go func() {
		for {
			memstats := runtime.MemStats{}
			runtime.ReadMemStats(&memstats)
			dbg.alloc = memstats.HeapAlloc
			time.Sleep(1 * time.Second)
		}
	}()
	return dbg
}

var gomono = gomono12.NewFace()

func (p *Debug) Update(delta time.Duration, input [4]controller.Controller) {
	p.Node.Update(delta, input)

	if input[0].Pressed()&joybus.ButtonL != 0 {
		p.hidden = !p.hidden
	}
}

func (p *Debug) Render() {
	bounds := renderer.Bounds().Inset(renderer.Bounds().Dy() / 10)
	if !p.hidden {
		p.stats = p.stats[:0]
		p.stats = append(p.stats, []byte("Draw:  ")...)
		p.stats = strconv.AppendInt(p.stats, p.disp.Duration().Microseconds(), 10)
		p.stats = append(p.stats, []byte("us\nFPS:   ")...)
		p.stats = strconv.AppendFloat(p.stats, float64(p.disp.FPS()), 'f', 3, 32)
		p.stats = append(p.stats, []byte("\nAlloc: ")...)
		p.stats = strconv.AppendUint(p.stats, p.alloc, 10)

		renderer.DrawText(bounds, gomono, bounds.Min, color.RGBA{A: 0xff}, nil, p.stats)
	}
}

func (p *Debug) Z() int {
	return 1000
}
