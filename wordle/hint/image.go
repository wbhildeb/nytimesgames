package hint

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/wbhildeb/nytimesgames/imagegen"
)

func gridify(target string) []string {
	words := strings.Split(target, " ")
	if len(words) == 0 {
		return []string{}
	}

	// Find the longest word to set a minimum line width
	longestWord := 0
	for _, word := range words {
		if len(word) > longestWord {
			longestWord = len(word)
		}
	}

	// Try to pack words into lines, preferring to fit multiple words per line
	var lines []string
	currentLine := ""

	for i, word := range words {
		if currentLine == "" {
			// Start a new line with this word
			currentLine = word
		} else {
			// Try adding this word to the current line
			testLine := currentLine + " " + word

			// If adding this word would make the line too long, start a new line
			// Use longestWord * 1.5 as a reasonable max line width
			maxLineWidth := int(float64(longestWord))
			if len(testLine) > maxLineWidth {
				lines = append(lines, currentLine)
				currentLine = word
			} else {
				currentLine = testLine
			}
		}

		// Add the last line if we're at the end
		if i == len(words)-1 && currentLine != "" {
			lines = append(lines, currentLine)
		}
	}

	// Pad all lines to the same length
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	for i := range lines {
		lines[i] = fmt.Sprintf("%-*s", maxLen, lines[i])
	}

	return lines
}

// hintToColor converts a CharHint to a color.RGBA
func hintToColor(h CharHint) color.RGBA {
	switch h {
	case CharHintCorrect:
		return imagegen.ColorWordleGreen
	case CharHintIncorrectPosition:
		return imagegen.ColorWordleYellow
	case CharHintNone:
		return imagegen.ColorWordleGray
	default:
		return imagegen.ColorWhite
	}
}

// GenerateHintImage creates an image showing the guess with colored backgrounds
// matching the hint colors (green for correct, yellow for wrong position, gray for not in word)
func GenerateHintImage(guess, target string) (*image.RGBA, error) {
	lines := gridify(target)
	targetFlat := strings.Join(lines, "")
	hints := CalculateHintSpecial(guess, targetFlat)

	// Count total alphabetic characters in guess
	guessAlphaCount := 0
	for _, c := range guess {
		if isAsciiAlpha(c) {
			guessAlphaCount++
		}
	}

	// Build grid with guess characters and colors
	guessGrid := make([]string, len(lines))
	charColors := make([]color.RGBA, 0, len(targetFlat))
	charIndex := 0

	for lineNum, line := range lines {
		var guessLine strings.Builder
		for _, ch := range line {
			var displayChar rune
			var charColor color.RGBA

			if charIndex < len(hints) && hints[charIndex] != CharHintUnspecified {
				// Get the corresponding character from guess
				alphaCountInTarget := countAlphaUpTo(targetFlat, charIndex)

				// Check if the guess has enough alphabetic characters
				if alphaCountInTarget < guessAlphaCount {
					guessChar := ' '
					guessAlphaIndex := 0
					for _, gc := range guess {
						if isAsciiAlpha(gc) {
							if guessAlphaIndex == alphaCountInTarget {
								guessChar = gc
								break
							}
							guessAlphaIndex++
						}
					}
					displayChar = guessChar
				} else {
					// Guess is too short, show blank
					displayChar = ' '
				}

				charColor = hintToColor(hints[charIndex])
			} else {
				// For unspecified hints (non-alphabetic chars), hide alphabetic chars
				if isAsciiAlpha(ch) {
					displayChar = ' '
				} else {
					displayChar = ch
				}
				charColor = imagegen.ColorWhite
			}

			// Convert to uppercase
			guessLine.WriteRune([]rune(strings.ToUpper(string(displayChar)))[0])
			charColors = append(charColors, charColor)
			charIndex++
		}
		guessGrid[lineNum] = guessLine.String()
	}

	// Use imagegen to create the image
	return imagegen.GenerateGrid(guessGrid, charColors)
}

// countAlphaUpTo counts alphabetic characters in target up to position i
func countAlphaUpTo(s string, pos int) int {
	count := 0
	for i, c := range s {
		if i >= pos {
			break
		}
		if isAsciiAlpha(c) {
			count++
		}
	}
	return count
}
