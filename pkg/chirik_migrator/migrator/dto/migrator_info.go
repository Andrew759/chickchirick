package dto

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/file/dto"
	"fmt"
	"strconv"
)

type MigratorInfo struct {
	MigratorEnabled bool
	ORMType         int
	FileInfo        dto.FileInfo
	Err             error
}

func (mInfo *MigratorInfo) FillByFileInfo(fileInfo dto.FileInfo) {
	fields := fileInfo.Struct.Fields()
	tag := fields.Tag(config.MigratorTag)
	if tag == nil {
		mInfo.Err = fmt.Errorf("fields has invalid tags: %v", fields)
		return
	}
	mInfo.fillByTag(tag)
	mInfo.FileInfo = fileInfo
}

func (mInfo *MigratorInfo) fillByTag(tag *chirik_ast.Tag) {
	tValue := tag.Values[0]

	switch tag.Key {
	case config.MigratorTag:
		switch tValue {
		case config.MigratorEnabled:
			mInfo.MigratorEnabled = true
		case config.MigratorDisabled:
		default:
			mInfo.MigratorEnabled = false
		}
	case config.MigratorOrmType:
		switch tValue {
		case strconv.Itoa(config.GORM):
			mInfo.ORMType = config.GORM
		case strconv.Itoa(config.BUN):
			mInfo.ORMType = config.BUN
		default:
			mInfo.Err = fmt.Errorf("unknown ORM type: %s", tValue)
		}
	}

	return
}

func (mInfo *MigratorInfo) HasError() bool {
	return mInfo.Err != nil
}
