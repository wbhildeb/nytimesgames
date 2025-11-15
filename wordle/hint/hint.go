package hint

import (
	"fmt"
	"strings"

	"github.com/enescakir/emoji"
)

type CharHint int

type Hint []CharHint

const (
	CharHintUnspecified CharHint = iota
	CharHintNone
	CharHintIncorrectPosition
	CharHintCorrect
)

// CalculateHint but handles non alphabetic characters, different sized guesses and targets
func CalculateHintSpecial(guess, target string) Hint {
	// Simplify guess and target
	targetSimple := ""
	for _, c := range target {
		if isAsciiAlpha(c) {
			targetSimple += strings.ToUpper(string(c))
		}
	}
	guessSimple := ""
	for _, c := range guess {
		if isAsciiAlpha(c) {
			guessSimple += strings.ToUpper(string(c))
		}
	}
	minLen := min(len(guessSimple), len(targetSimple))
	targetSimple = targetSimple[:minLen]
	guessSimple = guessSimple[:minLen]

	hintSimple := CalculateHint(guessSimple, targetSimple)

	fmt.Println(targetSimple, guessSimple, hintSimple)

	// now we need to recontextualize the hint to the target
	hint := make(Hint, len(target))
	hintSimpleIndex := 0
	for i, c := range target {
		if !isAsciiAlpha(c) { // leave unspecified
			continue
		}

		if hintSimpleIndex == len(hintSimple) {
			hint[i] = CharHintNone
			continue
		}

		hint[i] = hintSimple[hintSimpleIndex]
		hintSimpleIndex++
	}

	return hint
}

func CalculateHint(guess, target string) Hint {
	if len(guess) != len(target) {
		return nil
	}
	n := len(guess)

	guess = strings.ToUpper(guess)
	target = strings.ToUpper(target)

	targetLetterFreq := letterFrequency(target)

	hint := make([]CharHint, n)

	// Start with Correct
	for i := range n {
		g, t := guess[i], target[i]

		if g == t {
			hint[i] = CharHintCorrect
			targetLetterFreq[t]--
		}
	}

	// Then Incorrect Placement
	for i := range n {
		g := guess[i]

		if targetLetterFreq[g] > 0 && hint[i] == CharHintUnspecified {
			hint[i] = CharHintIncorrectPosition
			targetLetterFreq[g]--
		}
	}

	// Then No Such Chars
	for i := range n {
		if hint[i] == CharHintUnspecified {
			hint[i] = CharHintNone
		}
	}

	return hint
}

func FormattedHint(guess, target string) (formatted string) {
	hint := CalculateHintSpecial(guess, target)

	for i, charPosHint := range hint {
		nextStr := ""

		switch charPosHint {
		case CharHintCorrect:
			nextStr = string(emoji.GreenSquare)
		case CharHintIncorrectPosition:
			nextStr = string(emoji.YellowSquare)
		case CharHintNone:
			nextStr = string(emoji.BlackLargeSquare)
		default:
			if i < len(target) {
				nextStr = string(target[i])
			} else {
				nextStr = string(guess[i])
			}
		}

		formatted += nextStr
	}

	return
}

func letterFrequency(input string) map[byte]int {
	output := make(map[byte]int)
	for _, c := range input {
		output[byte(c)]++
	}

	return output
}

func isAsciiLowercase(c rune) bool {
	return c >= 'a' && c <= 'z'
}

func isAsciiUppercase(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func isAsciiAlpha(c rune) bool {
	return isAsciiLowercase(c) || isAsciiUppercase(c)
}
