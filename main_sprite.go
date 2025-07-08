package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

// Game implements ebiten.Game interface
type Game struct {
	playerImage *ebiten.Image
	playerX     float64
	playerY     float64
	time        float64
}

// NewGame creates a new game instance
func NewGame() *Game {
	g := &Game{
		playerX: screenWidth / 2,
		playerY: screenHeight / 2,
	}
	
	// Create a simple sprite (a colored square for now)
	g.playerImage = ebiten.NewImage(32, 32)
	g.playerImage.Fill(color.RGBA{0x00, 0xff, 0x00, 0xff})
	
	// Draw a simple face on the sprite
	faceImg := g.playerImage
	
	// Eyes
	ebitenutil.DrawRect(faceImg, 8, 8, 4, 4, color.Black)
	ebitenutil.DrawRect(faceImg, 20, 8, 4, 4, color.Black)
	
	// Mouth
	ebitenutil.DrawRect(faceImg, 10, 20, 12, 2, color.Black)
	
	return g
}

// Update proceeds the game state
func (g *Game) Update() error {
	g.time += 0.02
	
	// Simple keyboard controls
	speed := 3.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.playerX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.playerX += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.playerY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.playerY += speed
	}
	
	// Keep player on screen
	g.playerX = math.Max(0, math.Min(g.playerX, float64(screenWidth-32)))
	g.playerY = math.Max(0, math.Min(g.playerY, float64(screenHeight-32)))
	
	return nil
}

// Draw draws the game screen
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill background
	screen.Fill(color.RGBA{0x87, 0xce, 0xeb, 0xff}) // Sky blue
	
	// Draw ground
	ebitenutil.DrawRect(screen, 0, float64(screenHeight-100), screenWidth, 100, 
		color.RGBA{0x22, 0x8b, 0x22, 0xff}) // Forest green
	
	// Draw some clouds
	for i := 0; i < 3; i++ {
		cloudX := float64(i*200) + math.Sin(g.time+float64(i))*20
		cloudY := 50.0 + float64(i*30)
		g.drawCloud(screen, cloudX, cloudY)
	}
	
	// Draw player with bobbing animation
	op := &ebiten.DrawImageOptions{}
	bobY := math.Sin(g.time*5) * 2
	op.GeoM.Translate(g.playerX, g.playerY+bobY)
	screen.DrawImage(g.playerImage, op)
	
	// Draw UI
	ebitenutil.DebugPrintAt(screen, "EdomaeElf - Sprite Example", 10, 10)
	ebitenutil.DebugPrintAt(screen, "Use Arrow Keys or WASD to move", 10, 30)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Position: (%.0f, %.0f)", g.playerX, g.playerY), 10, 50)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()), 10, 70)
}

// drawCloud draws a simple cloud shape
func (g *Game) drawCloud(screen *ebiten.Image, x, y float64) {
	cloudColor := color.RGBA{0xff, 0xff, 0xff, 0xcc}
	// Draw three circles to form a cloud
	DrawCircle(screen, x, y, 20, cloudColor)
	DrawCircle(screen, x+15, y, 25, cloudColor)
	DrawCircle(screen, x+30, y, 20, cloudColor)
}

// Layout returns the game's logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// DrawCircle is a helper function since ebitenutil doesn't have DrawCircle
func DrawCircle(dst *ebiten.Image, cx, cy, r float64, clr color.Color) {
	// Simple circle drawing using lines
	segments := 32
	for i := 0; i < segments; i++ {
		angle1 := float64(i) * 2 * math.Pi / float64(segments)
		angle2 := float64(i+1) * 2 * math.Pi / float64(segments)
		
		x1 := cx + r*math.Cos(angle1)
		y1 := cy + r*math.Sin(angle1)
		x2 := cx + r*math.Cos(angle2)
		y2 := cy + r*math.Sin(angle2)
		
		ebitenutil.DrawLine(dst, x1, y1, x2, y2, clr)
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EdomaeElf - Sprite Movement Demo")
	
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}