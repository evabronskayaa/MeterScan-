package dto

import (
	"backend/internal/errors"
	"backend/internal/util"
	"github.com/gin-gonic/gin"
	"net/mail"
)

const (
	ErrInvalidEmail  errors.SimpleError = "Некорректная почта"
	ErrShortPassword errors.SimpleError = "Пароль должен содержать не менее 6 символов"
)

type ValidatableForm interface {
	Validate(args ValidateArgs) error
}

type ValidateArgs struct {
	Ctx       *gin.Context
	ReCaptcha util.ReCaptcha
}

func validateForms(args ValidateArgs, forms ...ValidatableForm) error {
	for _, form := range forms {
		if err := form.Validate(args); err != nil {
			return err
		}
	}
	return nil
}

type EmailPasswordForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (f EmailPasswordForm) Validate(_ ValidateArgs) error {
	if _, err := mail.ParseAddress(f.Email); err != nil {
		return ErrInvalidEmail
	} else if len(f.Password) < 6 {
		return ErrShortPassword
	}
	return nil
}

type ReCaptchaForm struct {
	Recaptcha string `json:"recaptcha" binding:"required"`
}

func (f ReCaptchaForm) Validate(args ValidateArgs) error {
	return args.ReCaptcha.Verify(f.Recaptcha, args.Ctx.ClientIP())
}

type LoginForm struct {
	EmailPasswordForm
	ReCaptchaForm
}

func (f LoginForm) Validate(args ValidateArgs) error {
	return validateForms(args, &f.EmailPasswordForm, &f.ReCaptchaForm)
}

type RegisterForm struct {
	EmailPasswordForm
	ReCaptchaForm
}

func (f RegisterForm) Validate(args ValidateArgs) error {
	return validateForms(args, &f.EmailPasswordForm, &f.ReCaptchaForm)
}

type VerifyTokenForm struct {
	Token string `form:"token"`
}

type UpdatePredictionForm struct {
	ID            uint64 `json:"id" binding:"required"`
	MeterReadings string `json:"meter_readings" binding:"required"`
}
