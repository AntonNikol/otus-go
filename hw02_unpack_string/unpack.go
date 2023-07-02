package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

// Unpack Выполняет распаковку строки.
func Unpack(str string) (string, error) {
	var b strings.Builder
	runes := []rune(str)

	// Если 1 символ число
	if len(str) > 0 && unicode.IsDigit(runes[0]) {
		return "", ErrInvalidString
	}

	for i, v := range runes {
		switch {
		// Если следующий символ число
		case i+1 < len(runes) && unicode.IsDigit(runes[i+1]):
			if unicode.IsDigit(runes[i]) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(runes[i+1]))
			b.WriteString(strings.Repeat(string(v), count))
			continue
			// Если текущий символ число
		case unicode.IsDigit(v):
			continue
		default:
			b.WriteRune(v)
		}
	}
	return b.String(), nil
}
