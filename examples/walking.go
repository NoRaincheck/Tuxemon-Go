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

const (
	screenWidth  = 240
	screenHeight = 160

	frameOX    = 0
	frameOY    = 32
	gridSize   = 16
	frameCount = 3
	gridXSize  = screenWidth / gridSize
	gridYSize  = screenHeight / gridSize

	spriteWidth  = 16
	spriteHeight = 32

	frameSpeed = 100
)

const (
	adventureTickUpdateInterval = 10
)

var (
	runnerImage *ebiten.Image
)

type Game struct {
	count int

	// keeping track of the adventurer's state
	adventurerX         int
	adventurerY         int
	deltaX              int
	deltaY              int
	deltaXFrame         int
	deltaYFrame         int
	isMoving            bool
	adventureTickOffset int
	prevKeyPressed      ebiten.Key
}

type AdventurerInterimState struct {
	deltaX      int
	deltaY      int
	deltaXFrame int
	deltaYFrame int
	isMoving    bool
}

func getAdventurerInterimState(key ebiten.Key) AdventurerInterimState {
	switch key {
	case ebiten.KeyLeft:
		return AdventurerInterimState{
			deltaX:      -1,
			deltaY:      0,
			deltaXFrame: 0,
			deltaYFrame: 2,
			isMoving:    true,
		}
	case ebiten.KeyRight:
		return AdventurerInterimState{
			deltaX:      1,
			deltaY:      0,
			deltaXFrame: 0,
			deltaYFrame: 3,
			isMoving:    true,
		}
	case ebiten.KeyUp:
		return AdventurerInterimState{
			deltaX:      0,
			deltaY:      -1,
			deltaXFrame: 0,
			deltaYFrame: 1,
			isMoving:    true,
		}
	case ebiten.KeyDown:
		return AdventurerInterimState{
			deltaX:      0,
			deltaY:      1,
			deltaXFrame: 0,
			deltaYFrame: 0,
			isMoving:    true,
		}
	}
	return AdventurerInterimState{}
}

func (g *Game) updateAdventurerState() {
	if !g.isMoving {
		// check if a direction key is pressed for the character
		var keys []ebiten.Key
		var interimState AdventurerInterimState
		keys = inpututil.AppendPressedKeys(keys)

		if slices.Contains(keys, g.prevKeyPressed) {
			interimState = getAdventurerInterimState(g.prevKeyPressed)
		} else {
			for _, key := range keys {
				interimState = getAdventurerInterimState(key)
				if interimState.isMoving {
					break
				}
			}
		}

		// ensure never out of bounds
		if g.adventurerX+interimState.deltaX < 0 {
			interimState = AdventurerInterimState{}
		}
		if g.adventurerX+interimState.deltaX >= gridXSize {
			interimState = AdventurerInterimState{}
		}
		if g.adventurerY+interimState.deltaY < 0 {
			interimState = AdventurerInterimState{}
		}
		if g.adventurerY+interimState.deltaY >= (gridYSize - 1) {
			interimState = AdventurerInterimState{}
		}

		if interimState.isMoving {
			g.isMoving = true
			g.deltaX = interimState.deltaX
			g.deltaY = interimState.deltaY
			g.adventureTickOffset = g.count
			g.deltaXFrame = 0
			g.deltaYFrame = interimState.deltaYFrame
		}
	} else {
		// now update state to determine how to actually animate the adventurer
		if g.count-g.adventureTickOffset < adventureTickUpdateInterval {
			g.deltaXFrame = 1
		} else if g.count-g.adventureTickOffset < adventureTickUpdateInterval*2 {
			g.deltaXFrame = 2
		} else {
			g.deltaXFrame = 0
			g.isMoving = false
			g.adventurerX += g.deltaX
			g.adventurerY += g.deltaY
			g.deltaX = 0
			g.deltaY = 0
		}
	}
}

func (g *Game) Update() error {
	g.count++
	g.updateAdventurerState()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawGrid(screen)
	g.drawAdventurer(screen)
}

func (g *Game) drawAdventurer(screen *ebiten.Image) {
	offsetX := gridSize/2 + g.deltaXFrame*(gridSize/3)*g.deltaX
	offsetY := gridSize/2 + g.deltaXFrame*(gridSize/3)*g.deltaY
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(spriteWidth)/2, -float64(spriteHeight)/2)
	op.GeoM.Translate(float64(g.adventurerX*gridSize+offsetX), float64(g.adventurerY*gridSize+offsetY))
	sx, sy := frameOX+g.deltaXFrame*spriteWidth, frameOY*g.deltaYFrame
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+spriteWidth, sy+spriteHeight)).(*ebiten.Image), op)
}

func (g *Game) drawGrid(screen *ebiten.Image) {
	for x := 0; x < gridXSize; x++ {
		for y := 0; y <= gridYSize; y++ {
			if x%2 == y%2 {
				// note that the y grid starts with an offset due to the sprite size
				vector.DrawFilledRect(screen, float32(x*gridSize), float32(y*gridSize)-float32(gridSize/2), float32(gridSize), float32(gridSize), color.RGBA{0x80, 0x80, 0x80, 0xc0}, true)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) NewGame() {
	g.adventurerX = gridXSize / 2
	g.adventurerY = gridYSize / 2
}

func main() {
	game := &Game{}
	game.NewGame()
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(spritesheets.Adventurer_walk_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Walking")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
