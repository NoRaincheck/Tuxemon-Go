package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Common game configuration constants
const (
	ScreenWidth  = 240
	ScreenHeight = 160
	GridSize     = 16
	GridXSize    = ScreenWidth / GridSize
	GridYSize    = ScreenHeight / GridSize
)

// Game represents the main game state
type Game struct {
	count int
}

func (g *Game) drawStripeTransition(screen *ebiten.Image, count int) {
	// draw a pulsing retangle that moves from the top left to the bottom right
	ticksToEmpty := 50
	ticksToFill := 100
	ticksStatic := 50

	state := count % (ticksToEmpty + ticksToFill + ticksStatic)
	color := color.RGBA{0x4c, 0x4f, 0x69, 0xf5}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("state: %d", state))

	var width float32
	if state < ticksToEmpty {
		width = 0
		ebitenutil.DebugPrint(screen, "\nEmpty")
	} else if state < ticksToEmpty+ticksToFill {
		width = (float32(state-ticksToEmpty) / float32(ticksToFill)) * ScreenWidth
		ebitenutil.DebugPrint(screen, "\nFilling LR")
	} else if state < ticksToEmpty+ticksToFill+ticksStatic {
		width = ScreenWidth
		ebitenutil.DebugPrint(screen, "\nFull")
	}

	gridSize := GridSize / 2

	for i := 0; i < ScreenHeight/gridSize; i++ {
		if i%2 == 0 {
			vector.DrawFilledRect(screen, 0, float32(i*gridSize), width, float32(gridSize), color, true)
		} else {
			vector.DrawFilledRect(screen, ScreenWidth-width, float32(i*gridSize), width, float32(gridSize), color, true)
		}
	}

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

// Update updates the game state
func (g *Game) Update() error {
	g.count++
	return nil
}

// Draw renders the game
func (g *Game) Draw(screen *ebiten.Image) {
	DrawGrid(screen, GridSize, GridXSize, GridYSize)
	g.drawStripeTransition(screen, g.count)
}

// Layout returns the logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// NewGame creates a new game instance
func NewGame() *Game {
	return &Game{}
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("Transition")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Failed to run game:", err)
	}
}
