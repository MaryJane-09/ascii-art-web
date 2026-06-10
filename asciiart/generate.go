package asciiart

import (
	"strings"
)

func GenerateArt(input string, banner map[rune][]string) string {
	lines := SplitInput(input)

	var sb strings.Builder
	start := 0
	if len(lines) > 0 && lines[0] == "" {
		start = 1
	}
	for _, line := range lines[start:] {
		if line == "" {
			sb.WriteString("\n")
			continue
		}

		for row := 0; row < 8; row++ {
			for _, r := range line {
				art, ok := banner[r]
				if !ok {
					continue
				}
				sb.WriteString(art[row])
			}
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
