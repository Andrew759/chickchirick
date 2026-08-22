package dto

type UserProfileResponse struct {
	Name     string  `json:"name"`
	Surname  string  `json:"surname"`
	Login    string  `json:"login"`
	Phone    string  `json:"phone"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}
