package main

import (
	"container/heap"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	mikoScreenWidth  = 1024
	mikoScreenHeight = 768
	mikoTileSize     = 128
	mikoScaleFactor  = 0.5
	mikoMapWidth     = 16
	mikoMapHeight    = 12
	playerSpeed      = 2.0
	worshipperSpeed  = 1.0
	spawnInterval    = 300 // frames between spawns (5 seconds at 60fps)
	offeringDuration = 120 // frames to stay at donation box (2 seconds)

	// Pathfinding constants
	donationBoxX                  = 8
	donationBoxY                  = 4
	maxSearchRadius               = 5
	pathfindingProximityThreshold = 10.0

)

// TileID represents a tile by its x,y position in the tileset
type TileID struct {
	X, Y int
}

func (t TileID) ToIndex() int {
	return t.Y*8 + t.X
}

// Point represents a tile coordinate
type Point struct {
	X, Y int
}

// Node represents a node in the A* pathfinding algorithm
type Node struct {
	Point  Point
	Parent *Node
	G      float64 // Cost from start to this node
	H      float64 // Heuristic cost from this node to goal
	F      float64 // Total cost (G + H)
}

// NodeList implements a priority queue for A* algorithm
type NodeList []*Node

func (nl NodeList) Len() int           { return len(nl) }
func (nl NodeList) Less(i, j int) bool { return nl[i].F < nl[j].F }
func (nl NodeList) Swap(i, j int)      { nl[i], nl[j] = nl[j], nl[i] }

func (nl *NodeList) Push(x interface{}) {
	*nl = append(*nl, x.(*Node))
}

func (nl *NodeList) Pop() interface{} {
	old := *nl
	n := len(old)
	item := old[n-1]
	*nl = old[0 : n-1]
	return item
}

// Walkable tiles map for O(1) lookups
var walkableTilesMap = map[TileID]bool{
	{0, 1}: true, // Stone path (vertical)
	{0, 5}: true, // Grass
	{2, 1}: true, // Gravel/sand
	{0, 3}: true, // Stone stairs (top)
	{0, 6}: true, // Stone stairs (middle)
	{1, 6}: true, // Stone stairs (bottom)
	{4, 2}: true, // Shrine floor (left)
	{5, 2}: true, // Shrine floor (center)
}

type Player struct {
	X, Y   float64
	Width  float64
	Height float64
	Image  *ebiten.Image
}

// isWalkable checks if a tile at the given coordinates is walkable
func isWalkable(shrineMap [][]TileID, x, y int) bool {
	if x < 0 || x >= mikoMapWidth || y < 0 || y >= mikoMapHeight {
		return false
	}

	tile := shrineMap[y][x]
	return walkableTilesMap[tile]
}

// isValidPosition checks if a position is within map bounds
func isValidPosition(p Point) bool {
	return p.X >= 0 && p.X < mikoMapWidth && p.Y >= 0 && p.Y < mikoMapHeight
}

// manhattanDistance calculates the Manhattan distance between two points
func manhattanDistance(a, b Point) float64 {
	return math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
}

// findPath uses A* algorithm to find a path from start to goal
func findPath(shrineMap [][]TileID, start, goal Point) ([]Point, error) {
	// Validate input coordinates
	if !isValidPosition(start) {
		return nil, fmt.Errorf("invalid start position: (%d, %d)", start.X, start.Y)
	}
	if !isValidPosition(goal) {
		return nil, fmt.Errorf("invalid goal position: (%d, %d)", goal.X, goal.Y)
	}

	// Check if start and goal are walkable
	if !isWalkable(shrineMap, start.X, start.Y) {
		return nil, fmt.Errorf("start position is not walkable: (%d, %d)", start.X, start.Y)
	}
	if !isWalkable(shrineMap, goal.X, goal.Y) {
		return nil, fmt.Errorf("goal position is not walkable: (%d, %d)", goal.X, goal.Y)
	}

	// If start equals goal, return trivial path
	if start.X == goal.X && start.Y == goal.Y {
		return []Point{start}, nil
	}

	// Initialize open and closed lists
	openList := &NodeList{}
	heap.Init(openList)
	closedList := make(map[Point]*Node)

	startNode := &Node{
		Point:  start,
		Parent: nil,
		G:      0,
		H:      manhattanDistance(start, goal),
		F:      manhattanDistance(start, goal),
	}

	heap.Push(openList, startNode)

	directions := []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} // Down, Right, Up, Left

	for openList.Len() > 0 {
		// Get the node with lowest F value
		current := heap.Pop(openList).(*Node)

		// Add current node to closed list
		closedList[current.Point] = current

		// Check if we've reached the goal
		if current.Point.X == goal.X && current.Point.Y == goal.Y {
			// Reconstruct path efficiently
			path := make([]Point, 0)
			for node := current; node != nil; node = node.Parent {
				path = append(path, node.Point)
			}
			// Reverse path to get correct order (start to goal)
			for i := len(path)/2 - 1; i >= 0; i-- {
				opp := len(path) - 1 - i
				path[i], path[opp] = path[opp], path[i]
			}
			return path, nil
		}

		// Check all neighbors
		for _, dir := range directions {
			neighborPoint := Point{
				X: current.Point.X + dir.X,
				Y: current.Point.Y + dir.Y,
			}

			// Skip if neighbor is not walkable
			if !isWalkable(shrineMap, neighborPoint.X, neighborPoint.Y) {
				continue
			}

			// Skip if neighbor is in closed list
			if _, exists := closedList[neighborPoint]; exists {
				continue
			}

			// Calculate costs
			g := current.G + 1.0
			h := manhattanDistance(neighborPoint, goal)
			f := g + h

			// Check if neighbor is already in open list with better path
			found := false
			for _, node := range *openList {
				if node.Point.X == neighborPoint.X && node.Point.Y == neighborPoint.Y {
					if g < node.G {
						// Update existing node with better path
						node.G = g
						node.F = f
						node.Parent = current
						heap.Fix(openList, 0) // Re-heapify since we modified a node
					}
					found = true
					break
				}
			}

			// Add neighbor to open list if not already there
			if !found {
				neighborNode := &Node{
					Point:  neighborPoint,
					Parent: current,
					G:      g,
					H:      h,
					F:      f,
				}
				heap.Push(openList, neighborNode)
			}
		}
	}

	// No path found
	return nil, fmt.Errorf("no path found from (%d, %d) to (%d, %d)", start.X, start.Y, goal.X, goal.Y)
}

// WorshipperState represents the current state of a worshipper
type WorshipperState int

const (
	StateApproaching WorshipperState = iota
	StateOffering
	StateLeaving
)

// Worshipper represents a shrine visitor
type Worshipper struct {
	X, Y          float64
	Width, Height float64
	Image         *ebiten.Image
	State         WorshipperState
	Timer         int
	StartX        float64 // Starting position for leaving
	TargetX       float64 // Target position for leaving
	Speed         float64
	Color         color.RGBA // Tint color for variety
	Path          []Point    // Path to follow
	PathIndex     int        // Current position in path
	NextTarget    Point      // Next tile to move to
}

// pixelToTile converts pixel coordinates to tile coordinates
func pixelToTile(x, y float64) Point {
	return Point{
		X: int(x / (mikoTileSize * mikoScaleFactor)),
		Y: int(y / (mikoTileSize * mikoScaleFactor)),
	}
}

// tileToPixel converts tile coordinates to pixel coordinates (center of tile)
func tileToPixel(p Point) (float64, float64) {
	return float64(p.X)*mikoTileSize*mikoScaleFactor + (mikoTileSize*mikoScaleFactor)/2,
		float64(p.Y)*mikoTileSize*mikoScaleFactor + (mikoTileSize*mikoScaleFactor)/2
}

// findNearestWalkableTile finds the nearest walkable tile to the given position
func findNearestWalkableTile(shrineMap [][]TileID, x, y float64) (Point, error) {
	tilePos := pixelToTile(x, y)

	// If current position is walkable, return it
	if isWalkable(shrineMap, tilePos.X, tilePos.Y) {
		return tilePos, nil
	}

	// Search in expanding circles
	for radius := 1; radius <= maxSearchRadius; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				checkX := tilePos.X + dx
				checkY := tilePos.Y + dy
				if isWalkable(shrineMap, checkX, checkY) {
					return Point{checkX, checkY}, nil
				}
			}
		}
	}

	// Fallback to bottom center if no walkable tile found
	fallbackPoint := Point{mikoMapWidth / 2, mikoMapHeight - 1}
	if isWalkable(shrineMap, fallbackPoint.X, fallbackPoint.Y) {
		return fallbackPoint, nil
	}

	return Point{}, fmt.Errorf("no walkable tile found near position (%.2f, %.2f)", x, y)
}

// NewWorshipper creates a new worshipper at a random spawn position
func NewWorshipper(image *ebiten.Image, shrineMap [][]TileID) *Worshipper {
	// Spawn from random side of screen
	var startX, startY float64
	var targetX float64

	side := rand.Intn(2) // 0 = left, 1 = right
	if side == 0 {
		// Spawn from left
		startX = -50
		targetX = float64(mikoMapWidth)*mikoTileSize*mikoScaleFactor + 50
	} else {
		// Spawn from right
		startX = float64(mikoMapWidth)*mikoTileSize*mikoScaleFactor + 50
		targetX = -50
	}

	startY = float64(mikoMapHeight-1) * mikoTileSize * mikoScaleFactor // Bottom of screen

	// Random color tint for variety
	colors := []color.RGBA{
		{255, 255, 255, 255}, // White (no tint)
		{255, 200, 200, 255}, // Light red
		{200, 255, 200, 255}, // Light green
		{200, 200, 255, 255}, // Light blue
		{255, 255, 200, 255}, // Light yellow
	}

	worshipper := &Worshipper{
		X:         startX,
		Y:         startY,
		Width:     32,
		Height:    32,
		Image:     image,
		State:     StateApproaching,
		Timer:     0,
		StartX:    startX,
		TargetX:   targetX,
		Speed:     worshipperSpeed + rand.Float64()*0.5, // Random speed variation
		Color:     colors[rand.Intn(len(colors))],
		Path:      []Point{},
		PathIndex: 0,
	}

	// Find nearest walkable tile to starting position
	startTile, err := findNearestWalkableTile(shrineMap, startX, startY)
	if err != nil {
		// Log error but continue with fallback behavior
		log.Printf("Warning: Could not find walkable starting tile: %v", err)
		startTile = Point{mikoMapWidth / 2, mikoMapHeight - 1}
	}

	// Donation box is at predefined position
	donationBoxTile := Point{donationBoxX, donationBoxY}

	// Calculate path to donation box
	path, err := findPath(shrineMap, startTile, donationBoxTile)
	if err != nil {
		// Log error but continue with empty path (fallback behavior)
		log.Printf("Warning: Could not find path to donation box: %v", err)
		path = []Point{}
	}

	if len(path) > 0 {
		worshipper.Path = path
		worshipper.PathIndex = 0
		worshipper.NextTarget = path[0]
	}

	return worshipper
}

// Update updates the worshipper's state and position
func (w *Worshipper) Update(shrineMap [][]TileID) {
	switch w.State {
	case StateApproaching:
		// Follow the path to the donation box
		if len(w.Path) > 0 && w.PathIndex < len(w.Path) {
			targetX, targetY := tileToPixel(w.NextTarget)

			// Calculate direction to next path point
			dx := targetX - w.X
			dy := targetY - w.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			// If close enough to current target, move to next path point
			if distance < pathfindingProximityThreshold {
				w.PathIndex++
				if w.PathIndex < len(w.Path) {
					w.NextTarget = w.Path[w.PathIndex]
				} else {
					// Reached destination
					w.State = StateOffering
					w.Timer = 0
					return
				}
			} else {
				// Move towards current target
				if distance > 0 {
					w.X += (dx / distance) * w.Speed
					w.Y += (dy / distance) * w.Speed
				}
			}
		} else {
			// Fallback to direct movement if no path
			donationBoxPixelX := float64(donationBoxX) * mikoTileSize * mikoScaleFactor
			donationBoxPixelY := float64(donationBoxY) * mikoTileSize * mikoScaleFactor

			dx := donationBoxPixelX - w.X
			dy := donationBoxPixelY - w.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance < 20 {
				w.State = StateOffering
				w.Timer = 0
				return
			}

			if distance > 0 {
				w.X += (dx / distance) * w.Speed
				w.Y += (dy / distance) * w.Speed
			}
		}

	case StateOffering:
		// Stay at donation box for a while
		w.Timer++
		if w.Timer >= offeringDuration {
			w.State = StateLeaving
			w.Timer = 0

			// Calculate path to exit
			currentTile := pixelToTile(w.X, w.Y)
			exitTile, err := findNearestWalkableTile(shrineMap, w.TargetX, float64(mikoMapHeight-1)*mikoTileSize*mikoScaleFactor)
			if err != nil {
				// Log error but continue with fallback behavior
				log.Printf("Warning: Could not find walkable exit tile: %v", err)
				exitTile = Point{mikoMapWidth / 2, mikoMapHeight - 1}
			}

			exitPath, err := findPath(shrineMap, currentTile, exitTile)
			if err != nil {
				// Log error but continue with empty path (fallback behavior)
				log.Printf("Warning: Could not find path to exit: %v", err)
				exitPath = []Point{}
			}

			if len(exitPath) > 0 {
				w.Path = exitPath
				w.PathIndex = 0
				w.NextTarget = exitPath[0]
			}
		}

	case StateLeaving:
		// Follow the path to the exit
		if len(w.Path) > 0 && w.PathIndex < len(w.Path) {
			targetX, targetY := tileToPixel(w.NextTarget)

			// Calculate direction to next path point
			dx := targetX - w.X
			dy := targetY - w.Y
			distance := math.Sqrt(dx*dx + dy*dy)

			// If close enough to current target, move to next path point
			if distance < pathfindingProximityThreshold {
				w.PathIndex++
				if w.PathIndex < len(w.Path) {
					w.NextTarget = w.Path[w.PathIndex]
				} else {
					// Reached exit path, now move off-screen
					w.Path = []Point{}
				}
			} else {
				// Move towards current target
				if distance > 0 {
					w.X += (dx / distance) * w.Speed
					w.Y += (dy / distance) * w.Speed
				}
			}
		} else {
			// Move off-screen after following path
			dx := w.TargetX - w.X
			dy := (float64(mikoMapHeight) * mikoTileSize * mikoScaleFactor) - w.Y

			distance := math.Sqrt(dx*dx + dy*dy)

			if distance > 0 {
				w.X += (dx / distance) * w.Speed
				w.Y += (dy / distance) * w.Speed
			}
		}
	}
}

// IsOffScreen checks if worshipper is off screen and should be removed
func (w *Worshipper) IsOffScreen() bool {
	return w.State == StateLeaving && (w.X < -100 || w.X > float64(mikoMapWidth)*mikoTileSize*mikoScaleFactor+100)
}

type MikoGameWithWorshippers struct {
	tilemapImage    *ebiten.Image
	shrineMap       [][]TileID
	player          *Player
	cameraX         float64
	cameraY         float64
	editMode        bool
	selectedTile    TileID
	worshippers     []*Worshipper
	spawnTimer      int
	worshipperImage *ebiten.Image
	donationCount   int
	totalDonations  int
	usingFallback   bool // Track if we're using fallback images
}

func NewMikoGameWithWorshippers() *MikoGameWithWorshippers {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	usingFallback := false

	// Load the tilemap image with WebGL-compatible error handling
	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/tilemap/japanese_town_tileset.png")
	if err != nil {
		log.Printf("Failed to load tilemap image: %v", err)
		// Create a fallback colored image instead of fatal error
		tilemapImg = createFallbackTilemapImage()
		usingFallback = true
	}

	// Load player image with WebGL-compatible error handling
	playerImg, _, err := ebitenutil.NewImageFromFile("assets/characters/miko_girl.png")
	if err != nil {
		log.Printf("Failed to load player image: %v", err)
		// Create a fallback colored image instead of fatal error
		playerImg = createFallbackPlayerImage()
		usingFallback = true
	}

	// Note: Tile descriptions are not loaded to ensure WebGL compatibility

	// Create player
	player := &Player{
		X:      float64(mikoMapWidth/2) * mikoTileSize * mikoScaleFactor,
		Y:      float64(mikoMapHeight/2) * mikoTileSize * mikoScaleFactor,
		Width:  32,
		Height: 32,
		Image:  playerImg,
	}

	// Create the shrine map
	shrineMap := createMikoShrineMap()

	return &MikoGameWithWorshippers{
		tilemapImage:    tilemapImg,
		shrineMap:       shrineMap,
		player:          player,
		cameraX:         0,
		cameraY:         0,
		editMode:        false,
		selectedTile:    TileID{0, 0},
		worshippers:     make([]*Worshipper, 0),
		spawnTimer:      0,
		worshipperImage: playerImg, // Use same image as player for now
		donationCount:   0,
		totalDonations:  0,
		usingFallback:   usingFallback,
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
	shrineMap[8][6] = TileID{2, 2}  // 石灯籠（上部）
	shrineMap[9][6] = TileID{3, 2}  // 石灯籠（下部）
	shrineMap[8][10] = TileID{2, 2} // 石灯籠（上部）
	shrineMap[9][10] = TileID{3, 2} // 石灯籠（下部）

	// Stairs leading to shrine
	shrineMap[5][8] = TileID{0, 3} // 石階段（上段）
	shrineMap[6][8] = TileID{0, 6} // 石階段（中段）
	shrineMap[7][8] = TileID{1, 6} // 石階段（下段）

	// Create shrine building
	// Roof
	shrineMap[1][6] = TileID{4, 0} // 拝殿屋根（上端・左）
	shrineMap[1][7] = TileID{5, 0} // 拝殿屋根（上端・中左）
	shrineMap[1][8] = TileID{6, 0} // 拝殿屋根（上端・中右）
	shrineMap[1][9] = TileID{7, 0} // 拝殿屋根（上端・右）

	// Shrine walls
	shrineMap[2][6] = TileID{5, 1} // 拝殿壁（格子・左）
	shrineMap[2][7] = TileID{6, 1} // 拝殿壁（格子・中央）
	shrineMap[2][8] = TileID{4, 1} // 拝殿入口（正面）
	shrineMap[2][9] = TileID{7, 1} // 拝殿壁（格子・右）

	// Shrine floor
	shrineMap[3][6] = TileID{4, 2} // 拝殿縁側（左）
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

// createFallbackTilemapImage creates a fallback tilemap image when asset loading fails
func createFallbackTilemapImage() *ebiten.Image {
	// Create a 8x8 tilemap of 128x128 pixel tiles
	img := ebiten.NewImage(8*mikoTileSize, 8*mikoTileSize)
	
	// Define colors for different tile types
	colors := []color.RGBA{
		{50, 100, 50, 255},   // Green grass
		{100, 100, 100, 255}, // Gray stone
		{139, 69, 19, 255},   // Brown dirt
		{128, 128, 128, 255}, // Light gray
		{255, 182, 193, 255}, // Light pink
		{152, 251, 152, 255}, // Light green
		{176, 196, 222, 255}, // Light steel blue
		{255, 228, 181, 255}, // Light orange
	}
	
	// Fill tiles with different colors
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			colorIndex := (x + y*8) % len(colors)
			tileColor := colors[colorIndex]
			
			// Create a tile-sized image with the color
			tileImg := ebiten.NewImage(mikoTileSize, mikoTileSize)
			tileImg.Fill(tileColor)
			
			// Add a border for tile visibility
			borderColor := color.RGBA{0, 0, 0, 100}
			for i := 0; i < 2; i++ {
				// Top and bottom borders
				for tx := 0; tx < mikoTileSize; tx++ {
					tileImg.Set(tx, i, borderColor)
					tileImg.Set(tx, mikoTileSize-1-i, borderColor)
				}
				// Left and right borders
				for ty := 0; ty < mikoTileSize; ty++ {
					tileImg.Set(i, ty, borderColor)
					tileImg.Set(mikoTileSize-1-i, ty, borderColor)
				}
			}
			
			// Draw the tile to the main image
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*mikoTileSize), float64(y*mikoTileSize))
			img.DrawImage(tileImg, op)
		}
	}
	
	return img
}

// createFallbackPlayerImage creates a fallback player image when asset loading fails
func createFallbackPlayerImage() *ebiten.Image {
	// Create a simple colored rectangle to represent the player
	img := ebiten.NewImage(256, 256) // Use a reasonable size
	
	// Fill with a distinctive color (red for player)
	img.Fill(color.RGBA{255, 100, 100, 255})
	
	// Add a simple face or pattern
	centerX, centerY := 128, 128
	
	// Draw eyes
	eyeSize := 20
	eyeColor := color.RGBA{0, 0, 0, 255}
	
	// Left eye
	eyeImg := ebiten.NewImage(eyeSize, eyeSize)
	eyeImg.Fill(eyeColor)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(centerX-40), float64(centerY-20))
	img.DrawImage(eyeImg, op)
	
	// Right eye
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(float64(centerX+20), float64(centerY-20))
	img.DrawImage(eyeImg, op2)
	
	// Draw mouth
	mouthImg := ebiten.NewImage(40, 10)
	mouthImg.Fill(eyeColor)
	op3 := &ebiten.DrawImageOptions{}
	op3.GeoM.Translate(float64(centerX-20), float64(centerY+20))
	img.DrawImage(mouthImg, op3)
	
	return img
}

func (g *MikoGameWithWorshippers) Update() error {
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

	// Worshipper system updates
	g.updateWorshippers()

	return nil
}

func (g *MikoGameWithWorshippers) updateWorshippers() {
	// Increment spawn timer
	g.spawnTimer++

	// Spawn new worshipper randomly
	if g.spawnTimer >= spawnInterval && rand.Float64() < 0.3 { // 30% chance every spawn interval
		g.worshippers = append(g.worshippers, NewWorshipper(g.worshipperImage, g.shrineMap))
		g.spawnTimer = 0
	}

	// Update existing worshippers
	for i := 0; i < len(g.worshippers); i++ {
		worshipper := g.worshippers[i]
		oldState := worshipper.State

		worshipper.Update(g.shrineMap)

		// Check if worshipper just started offering (for donation counting)
		if oldState == StateApproaching && worshipper.State == StateOffering {
			g.donationCount++
			g.totalDonations++
		}

		// Remove worshippers that are off screen
		if worshipper.IsOffScreen() {
			g.worshippers = append(g.worshippers[:i], g.worshippers[i+1:]...)
			i--
		}
	}
}

func (g *MikoGameWithWorshippers) Draw(screen *ebiten.Image) {
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

	// Draw worshippers
	for _, worshipper := range g.worshippers {
		g.drawWorshipper(screen, worshipper)
	}

	// Draw player
	if !g.editMode {
		playerOp := &ebiten.DrawImageOptions{}
		playerScale := 0.125
		playerOp.GeoM.Scale(playerScale, playerScale)
		playerOp.GeoM.Translate(g.player.X-g.cameraX, g.player.Y-g.cameraY)
		screen.DrawImage(g.player.Image, playerOp)
	}

	// Draw UI
	info := fmt.Sprintf("巫女さんの神社探索 - 参拝客システム\nFPS: %.2f\n", ebiten.ActualFPS())
	info += fmt.Sprintf("参拝客数: %d\n", len(g.worshippers))
	info += fmt.Sprintf("現在の賽銭: %d\n", g.donationCount)
	info += fmt.Sprintf("総賽銭: %d\n", g.totalDonations)

	if g.usingFallback {
		info += "\n[WebGL互換モード] フォールバック画像を使用中\n"
	}

	if g.editMode {
		tileKey := fmt.Sprintf("%d,%d", g.selectedTile.X, g.selectedTile.Y)
		info += fmt.Sprintf("\n[編集モード]\n選択タイル: %s\n", tileKey)
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

func (g *MikoGameWithWorshippers) drawWorshipper(screen *ebiten.Image, worshipper *Worshipper) {
	op := &ebiten.DrawImageOptions{}

	// Scale worshipper
	scale := 0.1 // Smaller than player
	op.GeoM.Scale(scale, scale)

	// Position with camera offset
	op.GeoM.Translate(worshipper.X-g.cameraX, worshipper.Y-g.cameraY)

	// Apply color tint
	op.ColorScale.ScaleWithColor(worshipper.Color)

	// Add special effects for different states
	switch worshipper.State {
	case StateOffering:
		// Add a slight bounce effect while offering
		bounceOffset := math.Sin(float64(worshipper.Timer)*0.3) * 2
		op.GeoM.Translate(0, bounceOffset)
	case StateLeaving:
		// Fade out when leaving
		alpha := 1.0 - float64(worshipper.Timer)/300.0
		if alpha < 0.3 {
			alpha = 0.3
		}
		op.ColorScale.Scale(1, 1, 1, float32(alpha))
	}

	screen.DrawImage(worshipper.Image, op)
}

func (g *MikoGameWithWorshippers) Layout(outsideWidth, outsideHeight int) (int, int) {
	return mikoScreenWidth, mikoScreenHeight
}

// Run with: go run miko_game_with_worshippers.go
func main() {
	ebiten.SetWindowSize(mikoScreenWidth, mikoScreenHeight)
	ebiten.SetWindowTitle("EdomaeElf - 巫女さんの神社探索（参拝客システム）")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewMikoGameWithWorshippers()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
