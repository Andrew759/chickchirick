package file

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/chirik_migrator/config"
	cConfig "chickChirick/pkg/chirik_migrator/console/config"
	migratorDto "chickChirick/pkg/chirik_migrator/migrator/provider"
	"fmt"
	"github.com/spf13/viper"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"
)

func ReadEntityDir(entityNames []string) (map[string][]migratorDto.MigratorInfo, error) {
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
		migratorInfoList, err := ReadEntityFile(path, entityNames)
		if err != nil {
			return migratorEntities, err
		}

		if migratorInfoList != nil {
			migratorEntities[path] = migratorInfoList
		}
	}

	return migratorEntities, err
}

func ReadEntityFile(path string, entityNames []string) ([]migratorDto.MigratorInfo, error) {
	file, err := chirik_ast.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s", err)
	}
	if file == nil {
		return nil, fmt.Errorf("got invalid file %s", err)
	}

	eCount := len(entityNames)
	hasNameRestriction := eCount > 0
	if eCount == 1 && entityNames[0] == cConfig.AllKey {
		hasNameRestriction = false
	}

	var mInfoList []migratorDto.MigratorInfo

	entityNamespace := file.Package.Name()
	for _, structure := range file.Structures.List() {
		if hasNameRestriction && !slices.Contains(entityNames, structure.Name()) {
			continue
		}
		mInfo := migratorDto.MigratorInfo{}
		mInfo.EntityNamespace = entityNamespace

		mInfo.FillByEntity(*structure)

		mInfoList = append(mInfoList, mInfo)
	}

	return mInfoList, nil
}

func WriteSQLToFile(sql string, filePostfix string, migrationPath string) error {
	fileName := fmt.Sprintf("m%s_%s.sql",
		time.Now().Format("20060102_150405"),
		filePostfix,
	)
	fullFN := migrationPath + "/" + fileName

	return os.WriteFile(fullFN, []byte(sql), 0644)
}
