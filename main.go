package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Game implements ebiten.Game interface
type Game struct {
	counter int
}

// Update proceeds the game state
// Update is called every tick (1/60 [s] by default)
func (g *Game) Update() error {
	g.counter++
	return nil
}

// Draw draws the game screen
// Draw is called every frame (typically 60 times per second)
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with a nice blue color
	screen.Fill(color.RGBA{0x1a, 0x1c, 0x2e, 0xff})
	
	// Draw "Hello, Ebitengine!" text
	ebitenutil.DebugPrint(screen, "Hello, Ebitengine!")
	ebitenutil.DebugPrintAt(screen, "Welcome to EdomaeElf!", 10, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Frame: %d", g.counter), 10, 50)
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size
// If you don't have to adjust the screen size with the outside size, just return a fixed size
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("EdomaeElf - Ebitengine Example")
	
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}