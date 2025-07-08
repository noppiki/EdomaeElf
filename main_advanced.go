package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

// Game implements ebiten.Game interface
type Game struct {
	counter float64
}

// Update proceeds the game state
func (g *Game) Update() error {
	g.counter += 0.02
	return nil
}

// Draw draws the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill the screen with a gradient background
	for y := 0; y < screenHeight; y++ {
		ratio := float32(y) / float32(screenHeight)
		r := uint8(0x1a + ratio*30)
		g := uint8(0x1c + ratio*30)
		b := uint8(0x2e + ratio*30)
		vector.DrawFilledRect(screen, 0, float32(y), screenWidth, 1, color.RGBA{r, g, b, 0xff}, false)
	}
	
	// Draw animated circles
	centerX := float32(screenWidth / 2)
	centerY := float32(screenHeight / 2)
	
	// Draw multiple rotating circles
	for i := 0; i < 8; i++ {
		angle := g.counter + float64(i)*math.Pi/4
		x := centerX + float32(math.Cos(angle)*100)
		y := centerY + float32(math.Sin(angle)*100)
		
		// Create pulsing effect
		radius := float32(20 + math.Sin(g.counter*2+float64(i))*5)
		
		// Create color variation
		hue := float64(i) / 8.0
		r, g, b := hslToRGB(hue, 0.8, 0.6)
		
		vector.DrawFilledCircle(screen, x, y, radius, color.RGBA{r, g, b, 0xff}, false)
	}
	
	// Draw a central shape
	vector.DrawFilledCircle(screen, centerX, centerY, 40, color.RGBA{0xff, 0xff, 0xff, 0x80}, false)
	
	// Draw title and info
	ebitenutil.DebugPrintAt(screen, "EdomaeElf - Ebitengine Advanced Example", 10, 10)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), 10, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 10, 50)
	ebitenutil.DebugPrintAt(screen, "Press ESC to exit", 10, 70)
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// hslToRGB converts HSL to RGB color values
func hslToRGB(h, s, l float64) (uint8, uint8, uint8) {
	var r, g, b float64

	if s == 0 {
		r, g, b = l, l, l
	} else {
		hue2rgb := func(p, q, t float64) float64 {
			if t < 0 {
				t += 1
			}
			if t > 1 {
				t -= 1
			}
			if t < 1.0/6.0 {
				return p + (q-p)*6*t
			}
			if t < 1.0/2.0 {
				return q
			}
			if t < 2.0/3.0 {
				return p + (q-p)*(2.0/3.0-t)*6
			}
			return p
		}

		q := l + s - l*s
		if l < 0.5 {
			q = l * (1 + s)
		}
		p := 2*l - q
		r = hue2rgb(p, q, h+1.0/3.0)
		g = hue2rgb(p, q, h)
		b = hue2rgb(p, q, h-1.0/3.0)
	}

	return uint8(r * 255), uint8(g * 255), uint8(b * 255)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EdomaeElf - Advanced Ebitengine Demo")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}