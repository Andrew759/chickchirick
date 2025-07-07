package message

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model            `c_migrator:"enabled"`
	RuleInstallerId       int               `json:"rule_installer_id" gorm:"type:int"`
	RuleInstallerRelation UserRelation      `json:"read_user_relation" gorm:"references:RuleInstallerId"`
	RuleUserId            int               `json:"rule_user_id" gorm:"type:int"`
	RuleUserRelation      UserRelation      `json:"rule_user_relation" gorm:"references:RuleUserId"`
	Rule                  pgtype.JSONBCodec `json:"rule" gorm:"type:jsonb;default:'[]';not null"`
}
