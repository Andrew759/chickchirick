package helper

import (
	"regexp"
	"strings"
)

func ToSnakeCase(s string) string {
	// Разбиваем слова по границе между маленькой и большой буквой
	re := regexp.MustCompile("([a-z0-9])([A-Z])")

	return strings.ToLower(
		re.ReplaceAllString(s, "${1}_${2}"),
	)
}

func AddSingleSPostfix(s string) string {
	rs := []rune(s)
	sLen := len(rs)
	if sLen-1 >= 0 {
		lastCharter := string(rs[sLen-1:])
		if lastCharter == "d" &&
			sLen-2 >= 0 &&
			string(rs[sLen-2:]) == "ed" {
			return s
		}
		if lastCharter != "s" && lastCharter != "y" {
			s += "s"
		}
	}

	return s
}
