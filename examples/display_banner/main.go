package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	tuxemonFont "tuxemon/mods/tuxemon/font"
	"tuxemon/mods/tuxemon/gfx/borders"
)

// Common game configuration constants
const (
	ScreenWidth  = 240
	ScreenHeight = 160
	GridSize     = 16
	GridXSize    = ScreenWidth / GridSize
	GridYSize    = ScreenHeight / GridSize
)

const (
	bannerTarget = 3
	bannerStart  = bannerTarget - GridSize*3
	bannerPause  = 50 // frames to pause at bottom
)

// Banner represents an animated text banner
type Banner struct {
	X, Y      int
	deltaY    int
	animating bool
	tick      int
	text      string
	texts     []string
	textIndex int
}

// Game represents the main game state
type Game struct {
	count  int
	banner *Banner
	font   *text.GoTextFaceSource
	border *ebiten.Image
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

// NewBanner creates a new banner instance
func NewBanner() *Banner {
	return &Banner{
		X:         bannerTarget,
		Y:         bannerStart,
		text:      "Paper Town",
		texts:     []string{"Paper Town", "City Park"},
		textIndex: 0,
	}
}

// Update handles banner animation
func (b *Banner) Update(count int) {
	if !b.animating {
		return
	}

	b.Y += b.deltaY

	switch {
	case b.Y >= 3 && b.deltaY == 1:
		b.deltaY = 0
		b.tick = count
	case b.Y <= bannerStart:
		b.deltaY = 1
		b.animating = false
	case count-b.tick > bannerPause && b.deltaY == 0:
		b.deltaY = -1
	}
}

// StartAnimation begins the banner animation
func (b *Banner) StartAnimation() {
	b.animating = true
	b.X = bannerTarget
	b.Y = bannerStart
	b.deltaY = 1

	// Cycle through texts
	b.textIndex = (b.textIndex + 1) % len(b.texts)
	b.text = b.texts[b.textIndex]
}

func (g *Game) Update() error {
	g.count++

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.banner.StartAnimation()
	}

	g.banner.Update(g.count)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawGrid(screen, GridSize, GridXSize, GridYSize)
	g.drawBanner(screen)
}

func (g *Game) drawBanner(screen *ebiten.Image) {
	const fontSize = 16

	x, y := g.banner.X, g.banner.Y
	bannerText := fmt.Sprintf("  %s  ", g.banner.text)

	// Measure text dimensions
	face := &text.GoTextFace{
		Source: g.font,
		Size:   fontSize,
	}
	w, h := text.Measure(bannerText, face, 0)

	// Draw background
	vector.DrawFilledRect(
		screen,
		float32(x+borders.BorderSize),
		float32(y+borders.BorderSize),
		float32(w),
		float32(h),
		color.White,
		false,
	)

	// Draw border
	g.drawBorder(screen, x, y, int(w), int(h))

	// Draw text
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(x+borders.BorderSize), float64(y+borders.BorderSize))
	textOp.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, bannerText, face, textOp)
}

func (g *Game) drawBorder(screen *ebiten.Image, x, y, w, h int) {
	// Border tile positions in the tileset
	tiles := map[string]struct{ tileX, tileY int }{
		"top":         {borders.BorderSize, 0},
		"bottom":      {borders.BorderSize, borders.BorderSize * 2},
		"left":        {0, borders.BorderSize},
		"right":       {borders.BorderSize * 2, borders.BorderSize},
		"topLeft":     {0, 0},
		"topRight":    {borders.BorderSize * 2, 0},
		"bottomLeft":  {0, borders.BorderSize * 2},
		"bottomRight": {borders.BorderSize * 2, borders.BorderSize * 2},
	}

	// Helper function to draw a border tile
	drawTile := func(tileX, tileY, posX, posY int) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(posX), float64(posY))
		rect := image.Rect(tileX, tileY, tileX+borders.BorderSize, tileY+borders.BorderSize)
		screen.DrawImage(g.border.SubImage(rect).(*ebiten.Image), op)
	}

	// Draw corners
	drawTile(tiles["topLeft"].tileX, tiles["topLeft"].tileY, x, y)
	drawTile(tiles["topRight"].tileX, tiles["topRight"].tileY, x+w+borders.BorderSize, y)
	drawTile(tiles["bottomLeft"].tileX, tiles["bottomLeft"].tileY, x, y+h+borders.BorderSize)
	drawTile(tiles["bottomRight"].tileX, tiles["bottomRight"].tileY, x+w+borders.BorderSize, y+h+borders.BorderSize)

	// Draw horizontal borders
	for i := borders.BorderSize; i <= w; i++ {
		drawTile(tiles["top"].tileX, tiles["top"].tileY, i+x, y)
		drawTile(tiles["bottom"].tileX, tiles["bottom"].tileY, i+x, y+h+borders.BorderSize)
	}

	// Draw vertical borders
	for j := borders.BorderSize; j <= h; j++ {
		drawTile(tiles["left"].tileX, tiles["left"].tileY, x, j+y)
		drawTile(tiles["right"].tileX, tiles["right"].tileY, x+w+borders.BorderSize, j+y)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// NewGame creates a new game instance
func NewGame() *Game {
	// Load font
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(tuxemonFont.KenneyPixel_ttf))
	if err != nil {
		log.Fatal(err)
	}

	// Load border tiles
	img, _, err := image.Decode(bytes.NewReader(borders.Borders_png))
	if err != nil {
		log.Fatal(err)
	}
	borderTiles := ebiten.NewImageFromImage(img)

	return &Game{
		banner: NewBanner(),
		font:   fontSource,
		border: borderTiles,
	}
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(ScreenWidth*2, ScreenHeight*2)
	ebiten.SetWindowTitle("Display Banner")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
