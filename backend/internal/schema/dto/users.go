package dto

import "backend/internal/schema"

type LoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	ReCaptchaForm
}

type RegisterForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	ReCaptchaForm
}

type ReCaptchaForm struct {
	Recaptcha string `json:"recaptcha" binding:"required"`
}

type Token struct {
	Token   string `json:"token"`
	Expire  string `json:"expire" example:"2006-01-02T15:04:05Z07:00"`
	OrigIat string `json:"orig_iat" example:"2006-01-02T15:04:05Z07:00"`
}

type UserWithToken struct {
	User  *schema.User `json:"user"`
	Token *Token       `json:"token"`
}
