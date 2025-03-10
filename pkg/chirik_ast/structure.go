package chirik_ast

import (
	"go/ast"
)

type Structure struct {
	astTypeSpec *ast.TypeSpec
	astStruct   *ast.StructType
}

func (s *Structure) Name() string {
	return s.astTypeSpec.Name.Name
}

func (s *Structure) SetName(name string) {
	s.astTypeSpec.Name.Name = name
}

func (s *Structure) Fields() *Fields {
	return &Fields{
		astStruct: s.astStruct,
	}
}
