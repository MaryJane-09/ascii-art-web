package asciiart

import (
	"errors"
	"strings"
)

func ValidateInput(input string) (rune, error) {
	for _, r := range input {
		if r == '\n' || r == '\r'{
			continue
		}
		if r < 32 || r > 126 {
			return r, errors.New("invalid Character")
		}
	}
	return 0, nil
}

func SplitInput(input string) []string {
	return strings.Split(input, "\n")
}