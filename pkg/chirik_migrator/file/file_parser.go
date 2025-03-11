package file

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/file/dto"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/dto"
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

func readFile(path string, entityNames []string) migratorDto.MigratorInfo {
	file, err := chirik_ast.ReadFile(path)
	mInfo := migratorDto.MigratorInfo{}
	if err != nil {
		mInfo.Err = fmt.Errorf("error while reading file %s", err)
	}
	structureList := file.Structures.List()
	for _, structure := range structureList {
		if len(entityNames) > 0 && !slices.Contains(entityNames, structure.Name()) {
			continue
		}

		fileInfo := dto.FileInfo{}
		fileInfo.File = *file
		fileInfo.Struct = *structure

		mInfo := migratorDto.MigratorInfo{}
		err := mInfo.FillByFileInfo(fileInfo)
		if err != nil {
			mInfo.Err = err
		}
	}

	return nil
}
