package gen

import (
	"fmt"
	"testing"
)

func TestRandomString(t *testing.T) {
	alphaChars := "abcdefghijklmnopqrstuvwxy"
	numChars := "1234567890"
	symbolChars := "!@#$%^&*()"

	tests := []struct {
		chars  string
		length int
		want   int
	}{
		{alphaChars, 4, 4},
		{numChars, 8, 8},
		{symbolChars, 8, 8},
		{alphaChars, 12, 12},
		{"a", 4, 4},
		{alphaChars, 12, 12},
		{alphaChars, 16, 16},
		{alphaChars, 16, 16},
		{alphaChars, 0, 0},
		{alphaChars, -2, 0},
		{"", -1, DefaultLength},
	}

	generatedStrings := make(map[string]struct{})

	for _, tc := range tests {
		str := RandomString(tc.chars, tc.length)
		got := len(str)
		if got != tc.want {
			t.Fatalf("got length: %d want: %d", got, tc.want)
		}

		if tc.length <= 0 && str == "" {
			// Skip tracking 0 (or less) lengths, they should return an empty string
			// + if we track them in generatedStrings, the tests will fail
			continue
		}

		if _, ok := generatedStrings[str]; ok {
			t.Fatalf("duplicate string detected")
		}
		generatedStrings[str] = struct{}{}
	}

	fmt.Printf("all the strings! %v\n", generatedStrings)
}
