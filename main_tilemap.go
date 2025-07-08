package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	tileSize     = 16
	scaleFactor  = 3 // Scale tiles 3x for better visibility
)

type JapaneseTownGame struct {
	tilemapImage *ebiten.Image
	mapData      [][]int
	cameraX      float64
	cameraY      float64
}

func NewJapaneseTownGame() *JapaneseTownGame {
	// Load the tilemap image
	img, _, err := ebitenutil.NewImageFromFile("assets/tilemap/japanese_town_tileset.png")
	if err != nil {
		log.Fatal(err)
	}

	// Create a Japanese town map layout
	// -1 = empty, other numbers = tile indices
	mapData := [][]int{
		{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2},
		{16, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 17, 18},
		{16, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 18},
		{16, 33, 48, 49, 33, 33, 52, 53, 33, 33, 56, 57, 33, 33, 33, 18},
		{16, 33, 64, 65, 33, 33, 68, 69, 33, 33, 72, 73, 33, 33, 33, 18},
		{16, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 18},
		{16, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 18},
		{16, 35, 36, 35, 36, 35, 36, 35, 36, 35, 36, 35, 36, 35, 36, 18},
		{16, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 18},
		{16, 33, 80, 81, 33, 33, 84, 85, 33, 33, 88, 89, 33, 33, 33, 18},
		{16, 33, 96, 97, 33, 33, 100, 101, 33, 33, 104, 105, 33, 33, 33, 18},
		{16, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 33, 18},
		{32, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34, 34},
	}

	return &JapaneseTownGame{
		tilemapImage: img,
		mapData:      mapData,
		cameraX:      0,
		cameraY:      0,
	}
}

func (g *JapaneseTownGame) Update() error {
	// Camera movement with arrow keys
	speed := 5.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.cameraX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.cameraX += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.cameraY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.cameraY += speed
	}

	// Reset camera with Space
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.cameraX = 0
		g.cameraY = 0
	}

	return nil
}

func (g *JapaneseTownGame) Draw(screen *ebiten.Image) {
	// Clear screen with sky blue
	screen.Fill(color.RGBA{135, 206, 235, 255})

	// Get tilemap dimensions
	tilemapWidth := g.tilemapImage.Bounds().Dx()
	tilesPerRow := tilemapWidth / tileSize

	// Draw the map
	for y, row := range g.mapData {
		for x, tileID := range row {
			if tileID < 0 {
				continue // Skip empty tiles
			}

			// Calculate source position in tilemap
			srcX := (tileID % tilesPerRow) * tileSize
			srcY := (tileID / tilesPerRow) * tileSize

			// Calculate destination position on screen with camera offset
			destX := float64(x*tileSize*scaleFactor) - g.cameraX
			destY := float64(y*tileSize*scaleFactor) - g.cameraY

			// Skip tiles outside the screen
			if destX < -tileSize*scaleFactor || destX > screenWidth ||
				destY < -tileSize*scaleFactor || destY > screenHeight {
				continue
			}

			// Draw the tile
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scaleFactor, scaleFactor)
			op.GeoM.Translate(destX, destY)

			// Create sub-image for the specific tile
			screen.DrawImage(g.tilemapImage.SubImage(
				image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize),
			).(*ebiten.Image), op)
		}
	}

	// Display controls
	controls := fmt.Sprintf(
		"Japanese Town Demo\nFPS: %0.2f\nArrow Keys: Move Camera\nSpace: Reset Camera\nCamera: (%.0f, %.0f)",
		ebiten.ActualFPS(), g.cameraX, g.cameraY,
	)
	ebitenutil.DebugPrint(screen, controls)
}

func (g *JapaneseTownGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Run this with: go run main_tilemap.go
func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EdomaeElf - Japanese Town")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewJapaneseTownGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}