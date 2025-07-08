# EdomaeElf - Ebitengine Setup

This project has been set up with Ebitengine (Ebiten), a simple 2D game engine for Go.

## Installation

Ebitengine has been installed via:
```bash
go get -u github.com/hajimehoshi/ebiten/v2
```

## Example Files

Three example files have been created to demonstrate different aspects of Ebitengine:

### 1. main.go - Basic Example
A simple "Hello, Ebitengine!" application that demonstrates:
- Basic game loop structure
- Window setup
- Simple text rendering
- Frame counter

Run with:
```bash
go run main.go
```

### 2. main_advanced.go - Advanced Graphics Example
Demonstrates more advanced graphics features:
- Gradient backgrounds
- Animated shapes using vector graphics
- Color manipulation (HSL to RGB conversion)
- FPS/TPS display
- Window resizing

Run with:
```bash
go run main_advanced.go
```

### 3. main_sprite.go - Sprite & Input Example
A simple game-like example featuring:
- Sprite creation and rendering
- Keyboard input handling (Arrow keys or WASD)
- Simple animation (bobbing effect)
- Basic scene composition (sky, ground, clouds)
- Player movement with boundary checking

Run with:
```bash
go run main_sprite.go
```

## Project Structure
```
EdomaeElf/
├── go.mod              # Go module file
├── go.sum              # Go dependencies
├── main.go             # Basic Ebitengine example
├── main_advanced.go    # Advanced graphics example
├── main_sprite.go      # Sprite and input example
└── README.md          # This file
```

## Next Steps
You can now:
1. Choose one of the examples as a starting point for your game
2. Add image assets and load them using `ebitenutil.NewImageFromFile()`
3. Implement game logic, physics, collision detection, etc.
4. Add audio support using Ebitengine's audio package
5. Build for different platforms (Windows, macOS, Linux, Web, Mobile)

## Resources
- [Ebitengine Official Documentation](https://ebitengine.org/)
- [Ebitengine Examples](https://github.com/hajimehoshi/ebiten/tree/main/examples)
- [Ebitengine API Reference](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2)