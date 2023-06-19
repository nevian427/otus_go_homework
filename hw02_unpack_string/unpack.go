package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	const (
		stateNone = iota
		stateEsc
		stateSymbol
	)

	var (
		res   strings.Builder
		runes []rune = []rune(str)
		state int
	)

	for i := 0; i < len(runes); i++ {
		switch state {
		case stateNone:
			switch {
			case runes[i] >= '0' && runes[i] <= '9':
				return "", ErrInvalidString
			case runes[i] == '\\':
				state = stateEsc
			default:
				state = stateSymbol
			}
		case stateEsc:
			switch {
			case runes[i] == '\\', runes[i] >= '0' && runes[i] <= '9':
				state = stateSymbol
			default:
				return "", ErrInvalidString
			}
		case stateSymbol:
			switch {
			case runes[i] >= '0' && runes[i] <= '9':
				count, err := strconv.Atoi(string(runes[i]))
				if err != nil {
					return "", ErrInvalidString
				} else {
					res.WriteString(strings.Repeat(string(runes[i-1]), count))
					state = stateNone
				}
			case runes[i] == '\\':
				state = stateEsc
				res.WriteRune(runes[i-1])
			default:
				res.WriteRune(runes[i-1])
			}
		}
	}

	switch state {
	case stateEsc:
		return "", ErrInvalidString
	case stateSymbol:
		res.WriteRune(runes[len(runes)-1])
	}

	return res.String(), nil
}
