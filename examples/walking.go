package walking

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

// Walking game configuration constants
const (
	WalkingScreenWidth           = 240
	WalkingScreenHeight          = 160
	WalkingGridSize              = 16
	WalkingSpriteWidth           = 16
	WalkingSpriteHeight          = 32
	WalkingFrameOX               = 0
	WalkingFrameOY               = 32
	WalkingAnimationTickInterval = 10
	WalkingFrameCount            = 3
	WalkingGridXSize             = WalkingScreenWidth / WalkingGridSize
	WalkingGridYSize             = WalkingScreenHeight / WalkingGridSize
)

// WalkingMovementState represents the current movement state of the adventurer
type WalkingMovementState struct {
	DeltaX   int
	DeltaY   int
	FrameX   int
	FrameY   int
	IsMoving bool
}

// WalkingAdventurer represents the player character
type WalkingAdventurer struct {
	X, Y           int
	DeltaX, DeltaY int
	FrameX, FrameY int
	IsMoving       bool
	TickOffset     int
	PrevKey        ebiten.Key
}

// WalkingGame represents the main game state
type WalkingGame struct {
	count      int
	adventurer *WalkingAdventurer
	runnerImg  *ebiten.Image
}

// NewWalkingAdventurer creates a new adventurer instance
func NewWalkingAdventurer() *WalkingAdventurer {
	return &WalkingAdventurer{
		X: WalkingGridXSize / 2,
		Y: WalkingGridYSize / 2,
	}
}

// getWalkingMovementState returns the movement state for a given key
func getWalkingMovementState(key ebiten.Key) WalkingMovementState {
	switch key {
	case ebiten.KeyLeft:
		return WalkingMovementState{
			DeltaX:   -1,
			DeltaY:   0,
			FrameX:   0,
			FrameY:   2,
			IsMoving: true,
		}
	case ebiten.KeyRight:
		return WalkingMovementState{
			DeltaX:   1,
			DeltaY:   0,
			FrameX:   0,
			FrameY:   3,
			IsMoving: true,
		}
	case ebiten.KeyUp:
		return WalkingMovementState{
			DeltaX:   0,
			DeltaY:   -1,
			FrameX:   0,
			FrameY:   1,
			IsMoving: true,
		}
	case ebiten.KeyDown:
		return WalkingMovementState{
			DeltaX:   0,
			DeltaY:   1,
			FrameX:   0,
			FrameY:   0,
			IsMoving: true,
		}
	}
	return WalkingMovementState{}
}

// isWalkingWithinBounds checks if the given position is within the game bounds
func isWalkingWithinBounds(x, y int) bool {
	return x >= 0 && x < WalkingGridXSize && y >= 0 && y < (WalkingGridYSize-1)
}

// updateWalkingMovement handles the movement logic for the adventurer
func (a *WalkingAdventurer) updateWalkingMovement(count int) {
	if !a.IsMoving {
		a.handleWalkingInput(count)
	} else {
		a.updateWalkingAnimation(count)
	}
}

// handleWalkingInput processes input and starts movement if valid
func (a *WalkingAdventurer) handleWalkingInput(count int) {
	keys := inpututil.AppendPressedKeys(nil)
	var movementState WalkingMovementState
	if slices.Contains(keys, a.PrevKey) {
		movementState = getWalkingMovementState(a.PrevKey)
	} else {
		for _, key := range keys {
			movementState = getWalkingMovementState(key)
			if movementState.IsMoving {
				break
			}
		}
	}
	if !isWalkingWithinBounds(a.X+movementState.DeltaX, a.Y+movementState.DeltaY) {
		return
	}
	if movementState.IsMoving {
		a.IsMoving = true
		a.DeltaX = movementState.DeltaX
		a.DeltaY = movementState.DeltaY
		a.TickOffset = count
		a.FrameX = 0
		a.FrameY = movementState.FrameY
	}
}

// updateWalkingAnimation updates the animation frames during movement
func (a *WalkingAdventurer) updateWalkingAnimation(count int) {
	tickDiff := count - a.TickOffset
	if tickDiff < WalkingAnimationTickInterval {
		a.FrameX = 1
	} else if tickDiff < WalkingAnimationTickInterval*2 {
		a.FrameX = 2
	} else {
		a.FrameX = 0
		a.IsMoving = false
		a.X += a.DeltaX
		a.Y += a.DeltaY
		a.DeltaX = 0
		a.DeltaY = 0
	}
}

// Update updates the game state
func (g *WalkingGame) Update() error {
	g.count++
	g.adventurer.updateWalkingMovement(g.count)
	return nil
}

// Draw renders the game
func (g *WalkingGame) Draw(screen *ebiten.Image) {
	g.drawWalkingGrid(screen)
	g.drawWalkingAdventurer(screen)
}

// drawWalkingAdventurer renders the adventurer sprite
func (g *WalkingGame) drawWalkingAdventurer(screen *ebiten.Image) {
	offsetX := WalkingGridSize/2 + g.adventurer.FrameX*(WalkingGridSize/3)*g.adventurer.DeltaX
	offsetY := WalkingGridSize/2 + g.adventurer.FrameX*(WalkingGridSize/3)*g.adventurer.DeltaY
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(WalkingSpriteWidth)/2, -float64(WalkingSpriteHeight)/2)
	op.GeoM.Translate(
		float64(g.adventurer.X*WalkingGridSize+offsetX),
		float64(g.adventurer.Y*WalkingGridSize+offsetY),
	)
	sx := WalkingFrameOX + g.adventurer.FrameX*WalkingSpriteWidth
	sy := WalkingFrameOY * g.adventurer.FrameY
	spriteRect := image.Rect(sx, sy, sx+WalkingSpriteWidth, sy+WalkingSpriteHeight)
	screen.DrawImage(g.runnerImg.SubImage(spriteRect).(*ebiten.Image), op)
}

// drawWalkingGrid renders the background grid
func (g *WalkingGame) drawWalkingGrid(screen *ebiten.Image) {
	for x := 0; x < WalkingGridXSize; x++ {
		for y := 0; y <= WalkingGridYSize; y++ {
			if x%2 == y%2 {
				vector.DrawFilledRect(
					screen,
					float32(x*WalkingGridSize),
					float32(y*WalkingGridSize)-float32(WalkingGridSize/2),
					float32(WalkingGridSize),
					float32(WalkingGridSize),
					color.RGBA{0x80, 0x80, 0x80, 0xc0},
					true,
				)
			}
		}
	}
}

// Layout returns the logical screen size
func (g *WalkingGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WalkingScreenWidth, WalkingScreenHeight
}

// NewWalkingGame creates a new game instance
func NewWalkingGame() *WalkingGame {
	return &WalkingGame{
		adventurer: NewWalkingAdventurer(),
	}
}

// loadWalkingSprites loads the sprite assets
func loadWalkingSprites() (*ebiten.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(spritesheets.Adventurer_walk_png))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func main() {
	game := NewWalkingGame()
	runnerImg, err := loadWalkingSprites()
	if err != nil {
		log.Fatal("Failed to load sprites:", err)
	}
	game.runnerImg = runnerImg
	ebiten.SetWindowSize(WalkingScreenWidth*2, WalkingScreenHeight*2)
	ebiten.SetWindowTitle("Walking")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal("Failed to run game:", err)
	}
}
