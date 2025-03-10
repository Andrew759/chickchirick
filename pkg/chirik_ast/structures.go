package chirik_ast

import (
	"go/ast"
)

type Structures struct {
	ast *ast.File
}

func (s *Structures) Find(name string) (*Structure, bool) {
	list := s.List()
	for _, item := range list {
		if item.Name() == name {
			return item, true
		}
	}

	return nil, false
}

func (s *Structures) List() []*Structure {
	var result []*Structure

	ast.Inspect(s.ast, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		result = append(result, &Structure{
			astTypeSpec: typeSpec,
			astStruct:   structType,
		})

		return true
	})

	return result
}
