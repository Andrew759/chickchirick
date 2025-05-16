package chirik_ast

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type File struct {
	Package    *Package
	Structures *Structures

	filepath string
	astFile  *ast.File
	fileSet  *token.FileSet
}

func ReadFile(filepath string) (*File, error) {
	fSet := token.NewFileSet()

	astFile, err := parser.ParseFile(fSet, filepath, nil, 0)
	if err != nil {
		return nil, err
	}

	return &File{
		Package:    &Package{ast: astFile.Name},
		Structures: &Structures{ast: astFile},
		filepath:   filepath,
		astFile:    astFile,
		fileSet:    fSet,
	}, nil
}
