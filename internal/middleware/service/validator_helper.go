package service

import (
	"regexp"
	"strconv"
	"unicode/utf8"
)

func IsLatinSymbolOnly(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(s)
}

func IsPhoneNumber(s string) bool {
	return regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,20}$`).MatchString(s)
}

func IsEmail(s string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(s)
}

func IsHasCorrectLength(s string, length int) bool {
	return utf8.RuneCountInString(s) == length
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)

	return err == nil
}
