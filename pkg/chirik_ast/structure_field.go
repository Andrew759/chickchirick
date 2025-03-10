package chirik_ast

import (
	"go/ast"
)

type Field struct {
	astStruct *ast.StructType
	astField  *ast.Field
}

func (f *Field) Name() string {
	if f.astField.Names != nil && len(f.astField.Names[0].Name) > 0 {
		return f.astField.Names[0].Name
	}

	return ""
}

func (f *Field) Type() *Type {
	return &Type{f.astField.Type}
}

func (f *Field) Tags() *Tags {
	return &Tags{f.astField.Tag}
}
