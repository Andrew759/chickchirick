package chirik_ast

import (
	"chickChirick/pkg/chirik_migrator/config"
	"go/ast"
	"regexp"
	"strings"
)

type Tags struct {
	ast *ast.BasicLit
}

type Tag struct {
	Key    string
	Values []string
}

func (ts *Tags) List() []Tag {
	var tags []Tag

	trimmedTags := strings.Trim(ts.String(), "`")
	re := regexp.MustCompile(`(\w+):"([^"]+)"`)
	matches := re.FindAllStringSubmatch(trimmedTags, -1)

	if len(matches) > 0 {
		for _, match := range matches {
			key := strings.TrimLeft(strings.Trim(match[1], `"`), " ")
			values := strings.TrimLeft(strings.Trim(match[2], `"`), " ")

			valueList := strings.Split(values, ",")
			//TODO: нежелательная зависимость
			if key == config.MigratorGormTag {
				valueList = strings.Split(values, ";")
			}

			tags = append(tags, Tag{
				Key:    key,
				Values: valueList,
			})

		}
	}

	return tags
}

func (ts *Tags) Get(key string) (Tag, bool) {
	for _, t := range ts.List() {
		if t.Key == key {
			return t, true
		}
	}

	return Tag{}, false
}

func (ts *Tags) GetValue(key, value string) (string, bool) {
	if t, ok := ts.Get(key); ok {
		for _, v := range t.Values {
			if v == value {
				return v, true
			}
		}
	}

	return "", false
}

func (ts *Tags) HasValue(key, value string) bool {
	_, ok := ts.GetValue(key, value)

	return ok
}

func (ts *Tags) String() string {
	if ts.ast == nil {
		return ""
	}

	return ts.ast.Value
}
