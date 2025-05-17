package file

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/chirik_migrator/console/config"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"fmt"
	"github.com/spf13/viper"
	"io/fs"
	"path/filepath"
	"slices"
)

func ReadDir(entityNames []string) (map[string][]migratorDto.MigratorInfo, error) {
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

	fPathsLen := len(fPaths)
	migratorEntities := make(map[string][]migratorDto.MigratorInfo, fPathsLen)

	if err != nil {
		return migratorEntities, err
	}

	for _, path := range fPaths {
		//TODO: распараллелить? и подумать над более аккуратной обработкой ошибок
		migratorInfoList, err := readFile(path, entityNames)
		if err != nil {
			return migratorEntities, err
		}

		if migratorInfoList != nil {
			migratorEntities[path] = migratorInfoList
		}
	}

	return migratorEntities, err
}

func readFile(path string, entityNames []string) ([]migratorDto.MigratorInfo, error) {
	file, err := chirik_ast.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s", err)
	}
	if file == nil {
		return nil, fmt.Errorf("got invalid file %s", err)
	}

	hasNameRestriction := len(entityNames) > 0
	structureList := file.Structures.List()
	var mInfoList []migratorDto.MigratorInfo

	for _, structure := range structureList {
		mInfo := migratorDto.MigratorInfo{}

		if hasNameRestriction && !slices.Contains(entityNames, structure.Name()) {
			continue
		}

		mInfo.FillByEntity(*structure)

		mInfoList = append(mInfoList, mInfo)
	}

	return mInfoList, nil
}
