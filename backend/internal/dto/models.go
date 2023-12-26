package dto

type Token struct {
	Token   string `json:"token"`
	Expire  string `json:"expire" example:"2006-01-02T15:04:05Z07:00"`
	OrigIat string `json:"orig_iat" example:"2006-01-02T15:04:05Z07:00"`
}

type UserWithToken struct {
	User  User   `json:"user"`
	Token *Token `json:"token"`
}

type User struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	Validated bool   `json:"validated"`
}

type RecognitionResult struct {
	ID uint64 `json:"id"`
}

type Paging struct {
	Sort   string `json:"sort"`
	Page   string `json:"page"`
	Size   string `json:"size"`
	Order  string `json:"order"`
	Filter string `json:"filter"`
}
