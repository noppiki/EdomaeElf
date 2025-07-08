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
	mikoScreenWidth   = 1024
	mikoScreenHeight  = 768
	mikoTileSize      = 128
	mikoScaleFactor   = 0.5
	mikoMapWidth      = 16
	mikoMapHeight     = 12
	playerSpeed       = 2.0
)

// TileID represents a tile by its x,y position in the tileset
type TileID struct {
	X, Y int
}

func (t TileID) ToIndex() int {
	return t.Y*8 + t.X
}

type Player struct {
	X, Y   float64
	Width  float64
	Height float64
	Image  *ebiten.Image
}

type MikoGame struct {
	tilemapImage   *ebiten.Image
	tileDesc       map[string]string
	shrineMap      [][]TileID
	player         *Player
	cameraX        float64
	cameraY        float64
	editMode       bool
	selectedTile   TileID
}

func NewMikoGame() *MikoGame {
	// Load the tilemap image
	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/tilemap/japanese_town_tileset.png")
	if err != nil {
		log.Fatal(err)
	}

	// Load player image
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/characters/miko_girl.png")
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

	// Create player
	player := &Player{
		X:      float64(mikoMapWidth/2) * mikoTileSize * mikoScaleFactor,
		Y:      float64(mikoMapHeight/2) * mikoTileSize * mikoScaleFactor,
		Width:  32, // キャラクターサイズを半分に
		Height: 32,
		Image:  playerImg,
	}

	// Create the shrine map
	shrineMap := createMikoShrineMap()

	return &MikoGame{
		tilemapImage: tilemapImg,
		tileDesc:     tileDesc,
		shrineMap:    shrineMap,
		player:       player,
		cameraX:      0,
		cameraY:      0,
		editMode:     false,
		selectedTile: TileID{0, 0},
	}
}

func createMikoShrineMap() [][]TileID {
	// Initialize map with grass
	shrineMap := make([][]TileID, mikoMapHeight)
	for i := range shrineMap {
		shrineMap[i] = make([]TileID, mikoMapWidth)
		for j := range shrineMap[i] {
			shrineMap[i][j] = TileID{0, 5} // 草地（1）
		}
	}

	// Create the main path (stone path from bottom to shrine)
	for y := 8; y < mikoMapHeight; y++ {
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
	shrineMap[1][7] = TileID{5, 0}  // 拝殿屋根（上端・中左）
	shrineMap[1][8] = TileID{6, 0}  // 拝殿屋根（上端・中右）
	shrineMap[1][9] = TileID{7, 0}  // 拝殿屋根（上端・右）

	// Shrine walls
	shrineMap[2][6] = TileID{5, 1}  // 拝殿壁（格子・左）
	shrineMap[2][7] = TileID{6, 1}  // 拝殿壁（格子・中央）
	shrineMap[2][8] = TileID{4, 1}  // 拝殿入口（正面）
	shrineMap[2][9] = TileID{7, 1}  // 拝殿壁（格子・右）

	// Shrine floor
	shrineMap[3][6] = TileID{4, 2}  // 拝殿縁側（左）
	shrineMap[3][7] = TileID{5, 2}  // 拝殿縁側（中央）
	shrineMap[3][8] = TileID{5, 2}  // 拝殿縁側（中央）
	shrineMap[3][9] = TileID{5, 2}  // 拝殿縁側（中央）

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

func (g *MikoGame) Update() error {
	// Toggle edit mode
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.editMode = !g.editMode
	}

	if !g.editMode {
		// Player movement with WASD
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
			g.player.Y -= playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
			g.player.Y += playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
			g.player.X -= playerSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
			g.player.X += playerSpeed
		}

		// Keep player within map bounds
		mapWidthPixels := float64(mikoMapWidth) * mikoTileSize * mikoScaleFactor
		mapHeightPixels := float64(mikoMapHeight) * mikoTileSize * mikoScaleFactor
		
		if g.player.X < 0 {
			g.player.X = 0
		}
		if g.player.X > mapWidthPixels-g.player.Width {
			g.player.X = mapWidthPixels - g.player.Width
		}
		if g.player.Y < 0 {
			g.player.Y = 0
		}
		if g.player.Y > mapHeightPixels-g.player.Height {
			g.player.Y = mapHeightPixels - g.player.Height
		}

		// Camera follows player
		g.cameraX = g.player.X - float64(mikoScreenWidth)/2 + g.player.Width/2
		g.cameraY = g.player.Y - float64(mikoScreenHeight)/2 + g.player.Height/2

		// Keep camera within bounds
		maxCameraX := mapWidthPixels - float64(mikoScreenWidth)
		maxCameraY := mapHeightPixels - float64(mikoScreenHeight)
		
		if g.cameraX < 0 {
			g.cameraX = 0
		}
		if g.cameraX > maxCameraX {
			g.cameraX = maxCameraX
		}
		if g.cameraY < 0 {
			g.cameraY = 0
		}
		if g.cameraY > maxCameraY {
			g.cameraY = maxCameraY
		}
	} else {
		// Edit mode controls
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
			mapX := int((float64(mx) + g.cameraX) / (mikoTileSize * mikoScaleFactor))
			mapY := int((float64(my) + g.cameraY) / (mikoTileSize * mikoScaleFactor))

			if mapX >= 0 && mapX < mikoMapWidth && mapY >= 0 && mapY < mikoMapHeight {
				g.shrineMap[mapY][mapX] = g.selectedTile
			}
		}
	}

	// Reset camera with Space (only in edit mode)
	if g.editMode && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.cameraX = 0
		g.cameraY = 0
	}

	return nil
}

func (g *MikoGame) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(color.RGBA{135, 206, 235, 255})

	// Draw the map
	for y := 0; y < mikoMapHeight; y++ {
		for x := 0; x < mikoMapWidth; x++ {
			tile := g.shrineMap[y][x]

			// Calculate source position in tilemap
			srcX := tile.X * mikoTileSize
			srcY := tile.Y * mikoTileSize

			// Calculate destination position with camera offset
			destX := float64(x)*mikoTileSize*mikoScaleFactor - g.cameraX
			destY := float64(y)*mikoTileSize*mikoScaleFactor - g.cameraY

			// Skip tiles outside screen
			if destX < -mikoTileSize*mikoScaleFactor || destX > mikoScreenWidth ||
				destY < -mikoTileSize*mikoScaleFactor || destY > mikoScreenHeight {
				continue
			}

			// Draw tile
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(mikoScaleFactor, mikoScaleFactor)
			op.GeoM.Translate(destX, destY)

			screen.DrawImage(g.tilemapImage.SubImage(
				image.Rect(srcX, srcY, srcX+mikoTileSize, srcY+mikoTileSize),
			).(*ebiten.Image), op)
		}
	}

	// Draw player
	if !g.editMode {
		playerOp := &ebiten.DrawImageOptions{}
		playerScale := 0.125 // プレイヤーを更に小さく表示（0.25の半分）
		playerOp.GeoM.Scale(playerScale, playerScale)
		playerOp.GeoM.Translate(g.player.X-g.cameraX, g.player.Y-g.cameraY)
		screen.DrawImage(g.player.Image, playerOp)
	}

	// Draw UI
	info := fmt.Sprintf("巫女さんの神社探索\nFPS: %.2f\n", ebiten.ActualFPS())
	if g.editMode {
		tileKey := fmt.Sprintf("%d,%d", g.selectedTile.X, g.selectedTile.Y)
		tileDesc := g.tileDesc[tileKey]
		info += fmt.Sprintf("\n[編集モード]\n選択タイル: %s\n%s\n", tileKey, tileDesc)
		info += "Q/R: タイルX選択, T/Y: タイルY選択\n左クリック: タイル配置\nSpace: カメラリセット"
	} else {
		info += fmt.Sprintf("\nプレイヤー位置: (%.0f, %.0f)\n", g.player.X, g.player.Y)
		info += "WASD/矢印キー: 移動\nE: 編集モード切替"
	}
	
	ebitenutil.DebugPrint(screen, info)

	// Draw selected tile preview in edit mode
	if g.editMode {
		// Draw preview box
		previewX := mikoScreenWidth - 100
		previewY := 10
		previewSize := 64.0
		
		// Background
		ebitenutil.DrawRect(screen, float64(previewX-5), float64(previewY-5), 
			previewSize+10, previewSize+10, color.RGBA{0, 0, 0, 180})
		
		// Selected tile
		op := &ebiten.DrawImageOptions{}
		scale := previewSize / float64(mikoTileSize)
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(previewX), float64(previewY))
		
		srcX := g.selectedTile.X * mikoTileSize
		srcY := g.selectedTile.Y * mikoTileSize
		
		screen.DrawImage(g.tilemapImage.SubImage(
			image.Rect(srcX, srcY, srcX+mikoTileSize, srcY+mikoTileSize),
		).(*ebiten.Image), op)
	}
}

func (g *MikoGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return mikoScreenWidth, mikoScreenHeight
}

// Run with: go run miko_game.go
func main() {
	ebiten.SetWindowSize(mikoScreenWidth, mikoScreenHeight)
	ebiten.SetWindowTitle("EdomaeElf - 巫女さんの神社探索")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewMikoGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}