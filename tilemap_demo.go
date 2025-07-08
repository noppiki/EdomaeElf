package main

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	tileSize     = 16 // Each tile is 16x16 pixels
)

type TilemapGame struct {
	tilemapImage *ebiten.Image
	mapData      [][]int
}

func NewTilemapGame() *TilemapGame {
	// Load the tilemap image
	img, _, err := ebitenutil.NewImageFromFile("assets/tilemap/japanese_town_tileset.png")
	if err != nil {
		log.Fatal(err)
	}

	// Sample map data (16x12 tiles for 800x600 screen at 50x50 tile size)
	// Numbers represent different tiles in the tileset
	mapData := [][]int{
		{0, 0, 0, 0, 0, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 4, 5, 5, 5, 5, 5, 4, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 4, 5, 5, 5, 5, 5, 4, 1, 1, 1, 1},
		{2, 2, 2, 2, 2, 4, 5, 5, 5, 5, 5, 4, 2, 2, 2, 2},
		{2, 2, 2, 2, 2, 4, 4, 4, 4, 4, 4, 4, 2, 2, 2, 2},
		{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
		{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
		{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	return &TilemapGame{
		tilemapImage: img,
		mapData:      mapData,
	}
}

func (g *TilemapGame) Update() error {
	// Game update logic here
	return nil
}

func (g *TilemapGame) Draw(screen *ebiten.Image) {
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

			// Calculate destination position on screen
			// Scale tiles to 50x50 for better visibility
			destX := x * 50
			destY := y * 50

			// Draw the tile
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(50.0/float64(tileSize), 50.0/float64(tileSize))
			op.GeoM.Translate(float64(destX), float64(destY))

			// Create sub-image for the specific tile
			tileRect := ebiten.NewImageFromImage(g.tilemapImage.SubImage(
				image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize),
			).(*ebiten.Image))

			screen.DrawImage(tileRect, op)
		}
	}

	// Display FPS
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

func (g *TilemapGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func RunTilemapDemo() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("EdomaeElf - Japanese Town Tilemap Demo")
	
	game := NewTilemapGame()
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// To run this demo, call RunTilemapDemo() from your main function