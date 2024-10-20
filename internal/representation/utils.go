package representation

import (
	"strings"
	"unicode"
)

func escape(s string) string {
	return strings.Replace(s, "/", "", -1)
}

func capitalize(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
