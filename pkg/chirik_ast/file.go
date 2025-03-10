package chirik_ast

import (
	"fmt"
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

func (f *File) Print() {
	fmt.Println("--------------------")
	fmt.Printf("package %s\n", f.Package.Name())
	fmt.Println("--------------------")

	fmt.Println("structs:")
	for i, s := range f.Structures.List() {
		if i > 0 {
			fmt.Printf("\n")
		}

		fmt.Printf("-Name: %s\n", s.Name())

		for _, field := range s.Fields().List() {
			name := field.Name()
			if name == "" {
				name = "*empty*"
			}

			t := field.Tags().String()
			if t == "" {
				t = "*empty*"
			}

			fmt.Printf("--%s | %s | %s\n", name, field.Type(), t)

		}
	}
}
