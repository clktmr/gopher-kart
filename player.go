package main

import (
	"embed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"time"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/drivers/rspq/mixer"
	"github.com/clktmr/n64/rcp/texture"
)

type playerVariant string

var enginePCM io.ReaderAt

var (
	//go:embed assets/*-gopher.png
	//go:embed assets/engine-loop-1-normalized.pcm_s16be
	_playerFiles embed.FS
	playerFiles  cartfs.FS = cartfs.Embed(_playerFiles)
)

const (
	Burgundy playerVariant = "burgundy"
	Beige    playerVariant = "beige"
	Black    playerVariant = "black"
)

func init() {
	r, err := playerFiles.Open("assets/engine-loop-1-normalized.pcm_s16be")
	if err != nil {
		panic(err)
	}
	enginePCM = r.(io.ReaderAt)
}

type Player struct {
	Sprite

	hidden     bool
	controller int

	// audio
	sfxSource *mixer.Source
	sfxReader io.ReadSeeker
}

func NewPlayer(v playerVariant, controller int) *Player {
	r, err := playerFiles.Open("assets/" + string(v) + "-gopher.png")
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
	player.relativePos.Y += worldbounds.Dy() / 2

	player.sfxReader = mixer.Loop(io.NewSectionReader(enginePCM, 0, (1<<63 - 1)))
	player.sfxSource = mixer.NewSource(player.sfxReader, 8000)
	player.sfxSource.SetVolume(0.0, 0.5)
	mixer.SetSource(player.controller+6, player.sfxSource)

	return player
}

func (p *Player) Update(delta time.Duration, input [4]controller.Controller) {
	p.Sprite.Update(delta, input)

	stick := image.Point{int(input[p.controller].X()), -int(input[p.controller].Y())}
	stick = stick.Mul(int(delta.Microseconds()))

	p.relativePos = p.relativePos.Add(stick.Div(5e5))

	if !p.hidden {
		p.sfxSource.SetVolume(0.2, 0.5)
		hz := 8000.0 + (2000.0 * float32(input[p.controller].X()) / 85.0)
		p.sfxSource.SetSampleRate(uint(hz))
	} else {
		p.sfxSource.SetVolume(0.0, 0.5)
	}

	p.hidden = !input[p.controller].Present()
}

func (p *Player) Render(dst draw.Image) {
	if p.hidden {
		return
	}
	p.Sprite.Render(dst)
}
