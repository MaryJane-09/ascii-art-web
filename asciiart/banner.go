package asciiart

import (
	"os"
	"strings"
)

func BannerCheck(filename string) (map[rune][]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(content, "\n")

	banner := make(map[rune][]string)

	start := 1

	for i := 32; i <= 126; i++ {
		art := lines[start : start+8]
		banner[rune(i)] = art
		start += 9
	}
	return banner, nil
}
