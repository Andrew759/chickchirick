package dto

import "chickChirick/pkg/chirik_ast"

type FileInfo struct {
	Path   string
	File   chirik_ast.File
	Struct chirik_ast.Structure
}
