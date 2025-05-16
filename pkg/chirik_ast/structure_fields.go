package chirik_ast

import (
	"go/ast"
)

type Fields struct {
	astStruct *ast.StructType
}

func (fs *Fields) List() []*Field {
	out := make([]*Field, 0, len(fs.astStruct.Fields.List))

	for _, f := range fs.astStruct.Fields.List {
		if len(f.Names) < 0 {
			continue
		}

		out = append(out, &Field{
			astStruct: fs.astStruct,
			astField:  f,
		})
	}

	return out
}

func (fs *Fields) Find(name string) *Field {
	for _, f := range fs.List() {
		if f.Name() == name {
			return f
		}
	}

	return nil
}

func (fs *Fields) IsHasTag(tag string) bool {
	for _, f := range fs.List() {
		_, status := f.Tags().Get(tag)
		if status {
			return true
		}
	}
	return false
}

func (fs *Fields) Tag(tag string) *Tag {
	for _, f := range fs.List() {
		tag, status := f.Tags().Get(tag)
		if status {
			return &tag
		}
	}
	return nil
}
