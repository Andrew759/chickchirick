package service

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/migration/console/config"
	"chickChirick/pkg/migration/console/factory"
	"fmt"
	"github.com/spf13/viper"
	"io/fs"
	"path/filepath"
	"slices"
)

func ReadDir(entityNames []string) error {
	var fPaths []string
	err := filepath.WalkDir(viper.GetString(config.EntityPath),
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				fPaths = append(fPaths, path)
			}
			return nil
		})

	for _, path := range fPaths {
		//TODO: распараллелить?
		err = readFile(path, entityNames)
	}

	if err != nil {
		return fmt.Errorf("error while reading files %s", err)
	}

	return nil
}

func readFile(path string, entityNames []string) error {
	file, err := chirik_ast.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error while reading file %s", err)
	}
	structureList := file.Structures.List()
	for _, structure := range structureList {
		if len(entityNames) > 0 && !slices.Contains(entityNames, structure.Name()) {
			return nil
		}
		tag, err := factory.InitMigratorTag(structure.Fields())
		//TODO: временная строка
		println(tag, err)
	}

	return nil
}
