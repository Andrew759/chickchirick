package user

type Properties struct {
	UserId int  `json:"user_id" gorm:"type:int"`
	User   User `json:"user" gorm:"references:UserId"`
	//TODO: Доработать таймзоны. Сейчас архитектура описана таким образом,
	// что всё время хранится без таймзон. Возможно потребуется отдельный сервис или
	// пакет для ресолва временной зоны пользователя
	Timezone int8   `json:"timezone" gorm:"type:smallint"`
	Email    string `json:"email" gorm:"type:varchar(256)"`
}
