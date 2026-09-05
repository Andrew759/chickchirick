package dto

type UserBrief struct {
	ID       int    `json:"id"`
	UserUuid string `json:"userUuid"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Login    string `json:"login"`
}
