package chirik_ast

import (
	"fmt"
	"go/ast"
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
	sParts := strings.Fields(trimmedTags)

	for _, part := range sParts {
		sSplitParts := strings.Split(part, ":")
		sSplitParts[0] = strings.Trim(sSplitParts[0], `"`)

		vStr := strings.Join(sSplitParts[1:], ":")
		vStrTrim := strings.Trim(vStr, `"`)
		vStrTrimWithoutEndSep := strings.Replace(vStrTrim, ";", "", 1)
		values := strings.Split(vStrTrimWithoutEndSep, ",")

		tags = append(tags, Tag{
			Key:    sSplitParts[0],
			Values: values,
		})
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

func TagsToString(tags []Tag) string {
	if len(tags) == 0 {
		return ""
	}

	tagStrings := make([]string, 0, len(tags))
	for _, t := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf(`%s:"%s"`, t.Key, strings.Join(t.Values, ",")))
	}

	return fmt.Sprintf("`%s`", strings.Join(tagStrings, " "))
}
