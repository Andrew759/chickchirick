package factory

import (
	"chickChirick/pkg/chirik_ast"
	"chickChirick/pkg/chirik_migrator/console/config"
	"chickChirick/pkg/chirik_migrator/migrator/dto"
	"fmt"
	"strconv"
)

// InitMigratorTag TODO: возможно стоит вынести это в DTO?
func InitMigratorTag(fields *chirik_ast.Fields) (dto.MigratorTag, error) {
	mTagDto := dto.MigratorTag{}
	tag := fields.Tag(config.MigratorTag)
	if tag == nil {
		return mTagDto, fmt.Errorf("fields has invalid tags: %v", fields)
	}
	mTagDto, err := FillMigratorTagDto(mTagDto, tag)
	if err != nil {
		return mTagDto, err
	}

	return mTagDto, nil
}

func FillMigratorTagDto(mTagDto dto.MigratorTag, tag *chirik_ast.Tag) (dto.MigratorTag, error) {
	tValue := tag.Values[0]

	switch tag.Key {
	case config.MigratorTag:
		switch tValue {
		case config.MigratorEnabled:
			mTagDto.MigratorEnabled = true
		case config.MigratorDisabled:
		default:
			mTagDto.MigratorEnabled = false
		}
	case config.MigratorOrmType:
		switch tValue {
		case strconv.Itoa(config.GORM):
			mTagDto.ORMType = config.GORM
		case strconv.Itoa(config.BUN):
			mTagDto.ORMType = config.BUN
		default:
			return mTagDto, fmt.Errorf("unknown ORM type: %s", tValue)
		}
	}

	return mTagDto, nil
}
