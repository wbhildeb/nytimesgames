package imagegen

import (
	"image"
	"image/color"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	ColorWordleGreen  = color.RGBA{106, 170, 100, 255} // Wordle green
	ColorWordleYellow = color.RGBA{201, 180, 88, 255}  // Wordle yellow
	ColorWordleGray   = color.RGBA{120, 124, 126, 255} // Wordle gray
	ColorWhite        = color.RGBA{255, 255, 255, 255}
)

// GenerateGrid creates an image from a grid of characters and their corresponding hints.
// The grid parameter is a slice of strings where each string represents a row.
// All strings must be the same length.
// The hints parameter must have the same dimensions as the grid (flattened row-major order).
func GenerateGrid(grid []string, charColors []color.RGBA) (*image.RGBA, error) {
	if len(grid) == 0 {
		return nil, nil
	}

	// Image dimensions
	charWidth := 60
	charHeight := 70
	padding := 6
	fontSize := 32.0

	gridHeight := len(grid)
	gridWidth := len(grid[0])

	// Calculate image size based on grid dimensions
	width := gridWidth*(charWidth+padding) + padding
	height := gridHeight*(charHeight+padding) + padding

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background with white
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, ColorWhite)
		}
	}

	// Load font
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	// Create freetype context
	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(ttfFont)
	ctx.SetFontSize(fontSize)
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)

	// Draw each character with its hint color background
	charIndex := 0
	for lineNum, line := range grid {
		x := padding
		y := padding + lineNum*(charHeight+padding)

		for _, ch := range line {
			// Determine background color based on hint
			bgColor := ColorWordleGray
			if charIndex < len(charColors) {
				bgColor = charColors[charIndex]
			}

			// Draw background rectangle
			for dy := 0; dy < charHeight; dy++ {
				for dx := 0; dx < charWidth; dx++ {
					img.Set(x+dx, y+dy, bgColor)
				}
			}

			// Draw character text - use white on colored backgrounds, black on white
			textColor := ColorWhite
			if bgColor == ColorWhite {
				textColor = color.RGBA{0, 0, 0, 255}
			}

			ctx.SetSrc(image.NewUniform(textColor))
			pt := freetype.Pt(x+charWidth/2-int(fontSize/2.5), y+charHeight/2+int(fontSize/2.5))
			_, err := ctx.DrawString(string(ch), pt)
			if err != nil {
				return nil, err
			}

			x += charWidth + padding
			charIndex++
		}
	}

	return img, nil
}
