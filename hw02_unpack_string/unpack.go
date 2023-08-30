package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	rs := []rune(s)
	rsLen := len(rs)
	var b strings.Builder

	for i := 0; i < rsLen; i++ {
		if isDigit(rs[i]) {
			return "", ErrInvalidString
		}

		if rs[i] == '\\' {
			i++
		}

		if i == rsLen-1 || !isDigit(rs[i+1]) {
			b.WriteRune(rs[i])
			continue
		}

		unpacked := strings.Repeat(string(rs[i]), int(rs[i+1]-'0'))
		b.WriteString(unpacked)

		i++
	}

	return b.String(), nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
