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
	if len(str) == 0 {
		return "", nil
	}
	if unicode.IsDigit(rune(str[0])) {
		return "", ErrInvalidString
	}

	runes := []rune(str)
	var result strings.Builder

	for i, v := range runes {
		switch {
		// Если следующий символ число
		case i+1 < len(runes) && unicode.IsDigit(runes[i+1]):
			if unicode.IsDigit(runes[i]) {
				return "", ErrInvalidString
			}
			count, _ := strconv.Atoi(string(runes[i+1]))
			result.WriteString(strings.Repeat(string(v), count))
			continue
			// Если текущий символ число
		case unicode.IsDigit(v):
			continue
		default:
			result.WriteRune(v)
		}
	}
	return result.String(), nil
}
