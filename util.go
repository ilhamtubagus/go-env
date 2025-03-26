package go_env

import (
	"strings"
	"unicode"
)

func snakeToCamelCase(input string) string {
	words := strings.Split(strings.ToLower(input), "_")
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			r := []rune(words[i])
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, "")
}
