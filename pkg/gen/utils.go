package gen

import (
	"fmt"
	"unicode"
)

func markerName(prefix, name string) string {
	return fmt.Sprintf("%s:%s", prefix, name)
}

func title(s string) string {
	r := []rune(s)

	return string(append([]rune{unicode.ToUpper(r[0])}, r[1:]...))
}