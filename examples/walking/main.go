package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"tuxemon/mods/tuxemon/spritesheets"
)

// Common game configuration constants
const (
	ScreenWidth  = 240
	ScreenHeight = 160
	GridSize     = 16
	GridXSize    = ScreenWidth / GridSize
	GridYSize    = ScreenHeight / GridSize
)

// Game configuration
const (
	SpriteWidth  = 16
	SpriteHeight = 32
	FrameOX      = 0
	FrameOY      = 32
	AnimTicks    = 10
	FrameCount   = 3
)

// Movement directions
var directions = map[ebiten.Key]struct {
	dx, dy, frameY int
}{
	ebiten.KeyLeft:  {-1, 0, 2},
	ebiten.KeyRight: {1, 0, 3},
	ebiten.KeyUp:    {0, -1, 1},
	ebiten.KeyDown:  {0, 1, 0},
}

// Adventurer represents the player character
type Adventurer struct {
	X, Y       int
	dx, dy     int
	frameX     int
	frameY     int
	moving     bool
	tickOffset int
	prevKey    ebiten.Key
}

// Game represents the main game state
type Game struct {
	count      int
	adventurer *Adventurer
	spriteImg  *ebiten.Image
}

// DrawGrid renders a checkerboard grid pattern
func DrawGrid(screen *ebiten.Image, gridSize, gridXSize, gridYSize int) {
	for x := 0; x < gridXSize; x++ {
		for y := 0; y <= gridYSize; y++ {
			if x%2 == y%2 {
				vector.DrawFilledRect(
					screen,
					float32(x*gridSize),
					float32(y*gridSize)-float32(gridSize/2),
					float32(gridSize),
					float32(gridSize),
					color.RGBA{0x80, 0x80, 0x80, 0xc0},
					true,
				)
			}
		}
	}
}

// NewAdventurer creates a new adventurer instance
func NewAdventurer() *Adventurer {
	return &Adventurer{
		X: GridXSize / 2,
		Y: GridYSize / 2,
	}
}

// inBounds checks if position is within game bounds
func inBounds(x, y int) bool {
	return x >= 0 && x < GridXSize && y >= 0 && y < (GridYSize-1)
}

// Update handles movement and animation
func (a *Adventurer) Update(count int) {
	if !a.moving {
		a.handleWalkingInput(count)
	} else {
		a.updateAnimation(count)
	}
}

// handleWalkingInput processes input and starts movement if valid
func (a *Adventurer) handleWalkingInput(count int) {
	keys := inpututil.AppendPressedKeys(nil)

	// Check for previous key first, then any pressed key
	var key ebiten.Key
	if slices.Contains(keys, a.prevKey) {
		key = a.prevKey
	} else if len(keys) > 0 {
		key = keys[0]
	} else {
		return
	}

	dir, ok := directions[key]
	if !ok {
		return
	}

	if !inBounds(a.X+dir.dx, a.Y+dir.dy) {
		return
	}

	a.moving = true
	a.dx = dir.dx
	a.dy = dir.dy
	a.tickOffset = count
	a.frameX = 0
	a.frameY = dir.frameY
	a.prevKey = key
}

// updateAnimation updates animation frames during movement
func (a *Adventurer) updateAnimation(count int) {
	tickDiff := count - a.tickOffset
	if tickDiff < AnimTicks {
		a.frameX = 1
	} else if tickDiff < AnimTicks*2 {
		a.frameX = 2
	} else {
		a.frameX = 0
		a.moving = false
		a.X += a.dx
		a.Y += a.dy
		a.dx, a.dy = 0, 0
	}
}

// Update updates the game state
func (g *Game) Update() error {
	g.count++
	g.adventurer.Update(g.count)
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	DrawGrid(screen, GridSize, GridXSize, GridYSize)
	g.drawAdventurer(screen)
}

// drawAdventurer renders the adventurer sprite
func (g *Game) drawAdventurer(screen *ebiten.Image) {
	a := g.adventurer

	// Calculate offset for smooth movement
	offsetX := GridSize/2 + a.frameX*(GridSize/3)*a.dx
	offsetY := GridSize/2 + a.frameX*(GridSize/3)*a.dy

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(SpriteWidth)/2, -float64(SpriteHeight)/2)
	op.GeoM.Translate(
		float64(a.X*GridSize+offsetX),
		float64(a.Y*GridSize+offsetY),
	)

	sx := FrameOX + a.frameX*SpriteWidth
	sy := FrameOY * a.frameY
	spriteRect := image.Rect(sx, sy, sx+SpriteWidth, sy+SpriteHeight)
	screen.DrawImage(g.spriteImg.SubImage(spriteRect).(*ebiten.Image), op)
}

// Layout returns the logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// NewGame creates a new game instance
func NewGame() *Game {
	return &Game{
		adventurer: NewAdventurer(),
	}
}

// loadSprites loads the sprite assets
func loadSprites() (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(spritesheets.Adventurer_walk_png))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func main() {
	game := NewGame()
	spriteImg, err := loadSprites()
	if err != nil {
		log.Fatal("Failed to load sprites:", err)
	}
	game.spriteImg = spriteImg

	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("Walking")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Failed to run game:", err)
	}
}
