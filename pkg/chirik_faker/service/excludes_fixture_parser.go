package service

import "strings"

func ParseExcludeString(s string) map[string]struct{} {
	fieldsResult := make(map[string]struct{})
	fields := strings.Split(s, ",")
	for i := range fields {
		fieldsResult[strings.TrimSpace(fields[i])] = struct{}{}
	}

	return fieldsResult
}
