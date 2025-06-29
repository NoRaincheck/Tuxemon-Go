// Copyright 2016 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

const (
	screenWidth  = 240
	screenHeight = 160

	gridSize  = 16
	gridXSize = screenWidth / gridSize
	gridYSize = screenHeight / gridSize

	bannerTarget = 3
	bannerStart  = bannerTarget - gridSize*3
)

var (
	fontSource  *text.GoTextFaceSource
	borderTiles *ebiten.Image
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(tuxemonFont.KenneyPixel_ttf))
	if err != nil {
		log.Fatal(err)
	}
	fontSource = s

	img, _, err := image.Decode(bytes.NewReader(borders.Borders_png))
	if err != nil {
		log.Fatal(err)
	}
	borderTiles = ebiten.NewImageFromImage(img)
}

type Game struct {
	count int

	bannerX      int
	bannerY      int
	bannerDeltaY int
	isAnimating  bool
	bannerTick   int
	bannerText   string
}

func (g *Game) Update() error {
	g.count++
	g.updateBannerPosition()
	return nil
}

func (g *Game) updateBannerPosition() {
	if !g.isAnimating {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// animate the banner from top to bottom, and then back up
			g.isAnimating = true

			g.bannerX = bannerTarget
			g.bannerY = bannerStart
			g.bannerDeltaY = 1
			if g.bannerText == `Paper Town` {
				g.bannerText = `City Park`
			} else {
				g.bannerText = `Paper Town`
			}
		}
	} else {
		// animate the banner from top to bottom, and then back up
		g.bannerY += g.bannerDeltaY
		if g.bannerY >= 3 && g.bannerDeltaY == 1 {
			g.bannerDeltaY = 0
			g.bannerTick = g.count
		} else if g.bannerY <= bannerStart {
			g.bannerDeltaY = 1
			g.isAnimating = false
		} else if g.count-g.bannerTick > 50 && g.bannerDeltaY == 0 {
			g.bannerDeltaY = -1
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawGrid(screen)
	g.drawBanner(screen)
}

func (g *Game) drawBanner(screen *ebiten.Image) {
	const (
		normalFontSize = 16
	)

	x, y := g.bannerX, g.bannerY
	bannerText := fmt.Sprintf(`  %s  `, g.bannerText)

	// get the size of the text
	face := &text.GoTextFace{
		Source: fontSource,
		Size:   normalFontSize,
	}
	w, h := text.Measure(bannerText, face, 0)

	// Draw the inside of the border
	vector.DrawFilledRect(screen, float32(x+borders.BorderSize), float32(y+borders.BorderSize), float32(w), float32(h), color.White, false)

	// Draw border using border tilemap
	g.drawBorder(screen, x, y, int(w), int(h))

	// Draw the sample text
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(float64(x+borders.BorderSize), float64(y+borders.BorderSize))
	textOp.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, bannerText, face, textOp)
}

func (g *Game) drawBorder(screen *ebiten.Image, x, y, w, h int) {
	// Border tile positions in the tileset
	const (
		topTileX, topTileY         = borders.BorderSize, 0
		bottomTileX, bottomTileY   = borders.BorderSize, borders.BorderSize * 2
		leftTileX, leftTileY       = 0, borders.BorderSize
		rightTileX, rightTileY     = borders.BorderSize * 2, borders.BorderSize
		topLeftX, topLeftY         = 0, 0
		topRightX, topRightY       = borders.BorderSize * 2, 0
		bottomLeftX, bottomLeftY   = 0, borders.BorderSize * 2
		bottomRightX, bottomRightY = borders.BorderSize * 2, borders.BorderSize * 2
	)

	// Helper function to draw a border tile
	drawTile := func(tileX, tileY, posX, posY int) {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(posX), float64(posY))
		rect := image.Rect(tileX, tileY, tileX+borders.BorderSize, tileY+borders.BorderSize)
		screen.DrawImage(borderTiles.SubImage(rect).(*ebiten.Image), op)
	}

	// Draw corners
	drawTile(topLeftX, topLeftY, x, y)                                                   // top-left
	drawTile(topRightX, topRightY, x+w+borders.BorderSize, y)                            // top-right
	drawTile(bottomLeftX, bottomLeftY, x, y+h+borders.BorderSize)                        // bottom-left
	drawTile(bottomRightX, bottomRightY, x+w+borders.BorderSize, y+h+borders.BorderSize) // bottom-right

	// Draw horizontal borders
	for i := borders.BorderSize; i <= w; i++ {
		drawTile(topTileX, topTileY, i+x, y)                            // top
		drawTile(bottomTileX, bottomTileY, i+x, y+h+borders.BorderSize) // bottom
	}

	// Draw vertical borders
	for j := borders.BorderSize; j <= h; j++ {
		drawTile(leftTileX, leftTileY, x, j+y)                        // left
		drawTile(rightTileX, rightTileY, x+w+borders.BorderSize, j+y) // right
	}
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

func main() {
	game := &Game{
		bannerX:     bannerTarget,
		bannerY:     bannerStart,
		isAnimating: false,
		bannerText:  `Paper Town`,
	}
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Display Banner")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
