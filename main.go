package main

import (
	"embedded/arch/r4000/systim"
	"embedded/rtos"
	"image"
	"log"
	"os"
	"syscall"

	"github.com/clktmr/n64/drivers/carts"
	"github.com/clktmr/n64/drivers/display"
	"github.com/clktmr/n64/drivers/rspq/mixer"
	_ "github.com/clktmr/n64/machine"
	"github.com/clktmr/n64/rcp/audio"
	"github.com/clktmr/n64/rcp/cpu"
	"github.com/clktmr/n64/rcp/video"

	"github.com/embeddedgo/fs/termfs"
)

func init() {
	systim.Setup(cpu.ClockSpeed)

	var err error
	var cart carts.Cart

	// Redirect stdout and stderr either to cart's logger
	if cart = carts.ProbeAll(); cart == nil {
		return
	}

	devConsole := termfs.NewLight("termfs", nil, cart)
	rtos.Mount(devConsole, "/dev/console")
	os.Stdout, err = os.OpenFile("/dev/console", syscall.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
	os.Stderr = os.Stdout

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
}

var players []*Player

var worldbounds image.Rectangle

func main() {
	video.Setup(false)
	resolution := video.NativeResolution()
	resolution.X /= 2
	log.Println("Enabling video:", resolution)
	disp := display.NewDisplay(resolution, video.BPP16)
	worldbounds = disp.Swap().Bounds()

	sampleRate := 48000
	log.Println("Enabling audio:", sampleRate, "Hz")
	audio.Start(sampleRate)
	mixer.Init()

	players = []*Player{
		NewPlayer(Burgundy, 0),
		NewPlayer(Beige, 1),
		NewPlayer(Black, 2),
	}
	game := NewNode(
		NewRoad(),
		players[0],
		players[1],
		players[2],
	)
	root := NewNode(
		NewDebug(disp),
		NewTitle(game),
	)

	gameloop := NewGameLoop(disp, root)

	log.Println("Starting gameloop")
	gameloop.Run()
}
