package dto

type User struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Validated bool   `json:"validated"`
}

type Token struct {
	Token   string `json:"token"`
	Expire  string `json:"expire" example:"2006-01-02T15:04:05Z07:00"`
	OrigIat string `json:"orig_iat" example:"2006-01-02T15:04:05Z07:00"`
}

type UserWithToken struct {
	User  User   `json:"user"`
	Token *Token `json:"token"`
}
