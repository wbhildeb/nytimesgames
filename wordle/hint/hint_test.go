package hint_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wbhildeb/nytimesgames/wordle/hint"
)

func TestCalculateHint(t *testing.T) {
	testCases := []struct {
		guess    string
		target   string
		expected hint.Hint
	}{
		{
			guess:    "PIPPLE",
			target:   "APPLES",
			expected: hint.Hint{hint.CharHintIncorrectPosition, hint.CharHintNone, hint.CharHintCorrect, hint.CharHintNone, hint.CharHintIncorrectPosition, hint.CharHintIncorrectPosition},
		},
	}

	for _, tc := range testCases {
		// t.Log(hint.FormattedHint(tc.guess, tc.target))
		hint := hint.CalculateHint(tc.guess, tc.target)
		assert.Equal(t, tc.expected, hint)
	}
}

func TestFormatHint(t *testing.T) {
	testCases := []struct {
		guess    string
		target   string
		expected string
	}{
		{
			guess:    "Truuu I forgot about cracker lmao",
			target:   "Suck my nuts, Winnipegger",
			expected: "🟨🟨🟨🟨 ⬛🟨 ⬛⬛⬛🟨, ⬛⬛⬛⬛⬛⬛⬛🟨⬛⬛⬛",
		},
	}

	for _, tc := range testCases {
		// t.Log(hint.FormattedHint(tc.guess, tc.target))
		hint := hint.FormattedHint(tc.guess, tc.target)
		assert.Equal(t, tc.expected, hint)
	}
}
