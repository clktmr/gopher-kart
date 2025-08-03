package main

import (
	"embed"
	"image"
	"image/draw"
	"image/png"
	"io"
	"math/rand"
	"time"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/clktmr/n64/drivers/controller"
	"github.com/clktmr/n64/drivers/rspq/mixer"
	"github.com/clktmr/n64/rcp/texture"
)

//go:embed assets/coin-shadow.png
//go:embed assets/sfx_coin_double1.pcm_s16be
var _coinPng embed.FS
var coinFiles cartfs.FS = cartfs.Embed(_coinPng)
var coinImg image.Image
var coinPCM io.ReaderAt

type Coin struct {
	Sprite

	hidden bool

	// audio
	sfxSource  *mixer.Source
	sfxReader  *io.SectionReader
	sfxChannel int
}

func init() {
	r, err := coinFiles.Open("assets/coin-shadow.png")
	if err != nil {
		panic(err)
	}
	coinImg, err = png.Decode(r)
	if err != nil {
		panic(err)
	}
	r, err = coinFiles.Open("assets/sfx_coin_double1.pcm_s16be")
	if err != nil {
		panic(err)
	}
	coinPCM = r.(io.ReaderAt)
}

func NewCoin(mixerChannel int) *Coin {
	tex := texture.NewTextureFromImage(coinImg)
	coin := &Coin{
		Sprite: *NewSprite(tex, 6, 1, 100*time.Millisecond),
	}
	coin.frame = rand.Intn(len(coin.frames))
	coin.hidden = true
	coin.sfxReader = io.NewSectionReader(coinPCM, 0, (1<<63 - 1))
	coin.sfxSource = mixer.NewSource(coin.sfxReader, 16000)
	coin.sfxChannel = mixerChannel
	return coin
}

func (p *Coin) Update(delta time.Duration, input [4]controller.Controller) {
	p.Sprite.Update(delta, input)

	if p.hidden {
		return
	}

	frame := p.frames[p.frame]
	r := image.Rectangle{Max: frame.Size()}.Add(p.globalPos)

	if r.Max.X < 0 {
		p.hidden = true
	}
	for _, player := range players {
		if !player.hidden && r.Overlaps(player.Bounds()) {
			p.hidden = true
			p.sfxReader.Seek(0, io.SeekStart)
			mixer.SetSource(p.sfxChannel, p.sfxSource)
			break
		}
	}
}

func (p *Coin) Render(dst draw.Image) {
	if !p.hidden {
		p.Sprite.Render(dst)
	}
}
