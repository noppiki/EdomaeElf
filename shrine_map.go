package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	shrineScreenWidth  = 1024
	shrineScreenHeight = 768
	shrineTileSize     = 128  // 各タイルは128x128ピクセル
	shrineScaleFactor  = 0.5  // タイルを半分のサイズで表示
	mapWidth           = 16
	mapHeight          = 12
)

// TileID represents a tile by its x,y position in the tileset
type TileID struct {
	X, Y int
}

func (t TileID) ToIndex() int {
	return t.Y*8 + t.X
}

type ShrineGame struct {
	tilemapImage   *ebiten.Image
	tileDesc       map[string]string
	shrineMap      [][]TileID
	cameraX        float64
	cameraY        float64
	selectedTile   TileID
	editMode       bool
}

func NewShrineGame() *ShrineGame {
	// Load the tilemap image
	img, _, err := ebitenutil.NewImageFromFile("assets/tilemap/japanese_town_tileset.png")
	if err != nil {
		log.Fatal(err)
	}

	// Load tile descriptions
	jsonData, err := ioutil.ReadFile("/Users/odakatoshikatsu/Downloads/shrine_tileset_8x8_description.json")
	if err != nil {
		log.Fatal(err)
	}

	var tileDesc map[string]string
	err = json.Unmarshal(jsonData, &tileDesc)
	if err != nil {
		log.Fatal(err)
	}

	// Create the shrine map
	shrineMap := createShrineMap()

	return &ShrineGame{
		tilemapImage: img,
		tileDesc:     tileDesc,
		shrineMap:    shrineMap,
		cameraX:      0,
		cameraY:      0,
		selectedTile: TileID{0, 0},
		editMode:     false,
	}
}

func createShrineMap() [][]TileID {
	// Initialize map with grass
	shrineMap := make([][]TileID, mapHeight)
	for i := range shrineMap {
		shrineMap[i] = make([]TileID, mapWidth)
		for j := range shrineMap[i] {
			// Default to grass
			shrineMap[i][j] = TileID{0, 5} // 草地（1）
		}
	}

	// Create the main path (stone path from bottom to shrine)
	for y := 8; y < mapHeight; y++ {
		shrineMap[y][8] = TileID{0, 1} // 石畳（縦）
	}

	// Place torii gate at entrance
	shrineMap[10][7] = TileID{0, 2} // 鳥居（左半分）
	shrineMap[10][8] = TileID{1, 2} // 鳥居（右半分）

	// Stone lanterns along the path
	shrineMap[8][6] = TileID{2, 2} // 石灯籠（上部）
	shrineMap[9][6] = TileID{3, 2} // 石灯籠（下部）
	shrineMap[8][10] = TileID{2, 2} // 石灯籠（上部）
	shrineMap[9][10] = TileID{3, 2} // 石灯籠（下部）

	// Stairs leading to shrine
	shrineMap[5][8] = TileID{0, 3} // 石階段（上段）
	shrineMap[6][8] = TileID{0, 6} // 石階段（中段）
	shrineMap[7][8] = TileID{1, 6} // 石階段（下段）

	// Create shrine building
	// Roof
	shrineMap[1][6] = TileID{4, 0}  // 拝殿屋根（上端・左）
	shrineMap[1][7] = TileID{5, 0} // 拝殿屋根（上端・中左）
	shrineMap[1][8] = TileID{6, 0} // 拝殿屋根（上端・中右）
	shrineMap[1][9] = TileID{7, 0} // 拝殿屋根（上端・右）

	// Shrine walls
	shrineMap[2][6] = TileID{5, 1}  // 拝殿壁（格子・左）
	shrineMap[2][7] = TileID{6, 1} // 拝殿壁（格子・中央）
	shrineMap[2][8] = TileID{4, 1} // 拝殿入口（正面）
	shrineMap[2][9] = TileID{7, 1} // 拝殿壁（格子・右）

	// Shrine floor
	shrineMap[3][6] = TileID{4, 2}  // 拝殿縁側（左）
	shrineMap[3][7] = TileID{5, 2} // 拝殿縁側（中央）
	shrineMap[3][8] = TileID{5, 2} // 拝殿縁側（中央）
	shrineMap[3][9] = TileID{5, 2} // 拝殿縁側（中央）

	// Place donation box
	shrineMap[4][8] = TileID{1, 4} // 賽銭箱

	// Hand washing basin
	shrineMap[6][5] = TileID{2, 4} // 手水舎（手洗い鉢）

	// Sacred tree
	shrineMap[3][3] = TileID{3, 1} // 御神木（上部）

	// Cherry trees
	// Left cherry tree
	shrineMap[2][1] = TileID{4, 4} // 桜の木（上部・左）
	shrineMap[2][2] = TileID{5, 4} // 桜の木（上部・右）
	shrineMap[3][1] = TileID{6, 5} // 桜の木（中段・左）
	shrineMap[3][2] = TileID{7, 5} // 桜の木（中段・右）
	shrineMap[4][1] = TileID{6, 6} // 桜の木（幹・左）
	shrineMap[4][2] = TileID{7, 6} // 桜の木（幹・右）
	shrineMap[5][1] = TileID{6, 7} // 桜の木（根元・左）
	shrineMap[5][2] = TileID{7, 7} // 桜の木（根元・右）

	// Right cherry tree
	shrineMap[2][13] = TileID{4, 4} // 桜の木（上部・左）
	shrineMap[2][14] = TileID{5, 4} // 桜の木（上部・右）
	shrineMap[3][13] = TileID{6, 5} // 桜の木（中段・左）
	shrineMap[3][14] = TileID{7, 5} // 桜の木（中段・右）
	shrineMap[4][13] = TileID{6, 6} // 桜の木（幹・左）
	shrineMap[4][14] = TileID{7, 6} // 桜の木（幹・右）
	shrineMap[5][13] = TileID{6, 7} // 桜の木（根元・左）
	shrineMap[5][14] = TileID{7, 7} // 桜の木（根元・右）

	// Add some gravel areas around the shrine
	for y := 4; y <= 6; y++ {
		for x := 6; x <= 10; x++ {
			if shrineMap[y][x].X == 0 && shrineMap[y][x].Y == 5 {
				shrineMap[y][x] = TileID{2, 1} // 敷砂利／砂地タイル
			}
		}
	}

	return shrineMap
}

func (g *ShrineGame) Update() error {
	// Toggle edit mode
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.editMode = !g.editMode
	}

	// Camera movement
	speed := 5.0
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cameraX -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cameraX += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cameraY -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cameraY += speed
	}

	// Reset camera
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.cameraX = 0
		g.cameraY = 0
	}

	// Edit mode controls
	if g.editMode {
		// Tile selection
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.selectedTile.X--
			if g.selectedTile.X < 0 {
				g.selectedTile.X = 7
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.selectedTile.X++
			if g.selectedTile.X > 7 {
				g.selectedTile.X = 0
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyT) {
			g.selectedTile.Y--
			if g.selectedTile.Y < 0 {
				g.selectedTile.Y = 7
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyY) {
			g.selectedTile.Y++
			if g.selectedTile.Y > 7 {
				g.selectedTile.Y = 0
			}
		}

		// Place tile with mouse
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			// Convert screen coordinates to map coordinates
			mapX := int((float64(mx) + g.cameraX) / (shrineTileSize * shrineScaleFactor))
			mapY := int((float64(my) + g.cameraY) / (shrineTileSize * shrineScaleFactor))

			if mapX >= 0 && mapX < mapWidth && mapY >= 0 && mapY < mapHeight {
				g.shrineMap[mapY][mapX] = g.selectedTile
			}
		}
	}

	return nil
}

func (g *ShrineGame) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{135, 206, 235, 255})

	// Draw the map
	for y := 0; y < mapHeight; y++ {
		for x := 0; x < mapWidth; x++ {
			tile := g.shrineMap[y][x]

			// Calculate source position in tilemap
			srcX := tile.X * shrineTileSize
			srcY := tile.Y * shrineTileSize

			// Calculate destination position with camera offset
			destX := float64(x)*shrineTileSize*shrineScaleFactor - g.cameraX
			destY := float64(y)*shrineTileSize*shrineScaleFactor - g.cameraY

			// Skip tiles outside screen
			if destX < -shrineTileSize*shrineScaleFactor || destX > shrineScreenWidth ||
				destY < -shrineTileSize*shrineScaleFactor || destY > shrineScreenHeight {
				continue
			}

			// Draw tile
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(shrineScaleFactor, shrineScaleFactor)
			op.GeoM.Translate(destX, destY)

			screen.DrawImage(g.tilemapImage.SubImage(
				image.Rect(srcX, srcY, srcX+shrineTileSize, srcY+shrineTileSize),
			).(*ebiten.Image), op)
		}
	}

	// Draw UI
	info := fmt.Sprintf("神社マップ\nFPS: %.2f\n", ebiten.ActualFPS())
	if g.editMode {
		tileKey := fmt.Sprintf("%d,%d", g.selectedTile.X, g.selectedTile.Y)
		tileDesc := g.tileDesc[tileKey]
		info += fmt.Sprintf("\n[編集モード]\n選択タイル: %s\n%s\n", tileKey, tileDesc)
		info += "Q/R: タイルX選択, T/Y: タイルY選択\n左クリック: タイル配置"
	} else {
		info += "\n[表示モード]\n矢印キー/WASD: カメラ移動\nSpace: カメラリセット\nE: 編集モード切替"
	}
	
	ebitenutil.DebugPrint(screen, info)

	// Draw selected tile preview in edit mode
	if g.editMode {
		// Draw preview box
		previewX := shrineScreenWidth - 100
		previewY := 10
		
		// Background
		previewSize := 64.0 // プレビューサイズを64x64に
		ebitenutil.DrawRect(screen, float64(previewX-5), float64(previewY-5), 
			previewSize+10, previewSize+10, color.RGBA{0, 0, 0, 180})
		
		// Selected tile
		op := &ebiten.DrawImageOptions{}
		scale := previewSize / float64(shrineTileSize) // 128px -> 64px
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(previewX), float64(previewY))
		
		srcX := g.selectedTile.X * shrineTileSize
		srcY := g.selectedTile.Y * shrineTileSize
		
		screen.DrawImage(g.tilemapImage.SubImage(
			image.Rect(srcX, srcY, srcX+shrineTileSize, srcY+shrineTileSize),
		).(*ebiten.Image), op)
	}
}

func (g *ShrineGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return shrineScreenWidth, shrineScreenHeight
}

// Run with: go run shrine_map.go
func main() {
	ebiten.SetWindowSize(shrineScreenWidth, shrineScreenHeight)
	ebiten.SetWindowTitle("EdomaeElf - 神社の境内")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewShrineGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}