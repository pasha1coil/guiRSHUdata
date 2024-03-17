package utils

import (
	"strings"
	"unicode/utf8"
)

func GetInitials(s string) string {
	words := strings.Fields(s)
	initials := ""
	for _, word := range words {
		r, _ := utf8.DecodeRuneInString(word)
		initials += string(r)
	}
	return initials
}
