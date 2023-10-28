package endpoint

import (
	"backend/internal/api-service/middleware"
	"backend/internal/api-service/service"
	"backend/internal/crud"
	"backend/internal/errors"
	"backend/internal/schema/dto"
	"backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

const ErrIncorrectPassword errors.SimpleError = "Неверный пароль"

type Service service.Service

// LoginHandler godoc
//
//	@Summary		Авторизация
//	@Tags			sessions
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.LoginForm	true	"User"
//	@Success		200		{object}	dto.UserWithToken
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/sessions [post]
func (s *Service) LoginHandler(c *gin.Context) {
	var loginForm dto.LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err := s.ReCaptcha.Verify(loginForm.Recaptcha, c.ClientIP()); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	email := loginForm.Email
	password := loginForm.Password

	user, err := crud.FindUser(s.DB, "email = ?", email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if !util.CheckPasswordHash(password, user.Password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrIncorrectPassword)
		return
	}

	obj, err := middleware.GenerateToken(s.JWTSecret, func(claims jwt.MapClaims) {
		claims["user_id"] = user.ID
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, dto.UserWithToken{
		User:  user,
		Token: obj,
	})
}

// RegisterHandler godoc
//
//	@Summary		Создание нового пользователя
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.RegisterForm	true	"Пользователь"
//	@Success		200		{object}	dto.UserWithToken
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/users [post]
func (s *Service) RegisterHandler(c *gin.Context) {
	var registerForm dto.RegisterForm
	if err := c.ShouldBind(&registerForm); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if err := s.ReCaptcha.Verify(registerForm.Recaptcha, c.ClientIP()); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if user, err := crud.CreateUser(s.DB, registerForm.Email, registerForm.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	} else {
		obj, err := middleware.GenerateToken(s.JWTSecret, func(claims jwt.MapClaims) {
			claims["user_id"] = user.ID
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, dto.UserWithToken{
			User:  user,
			Token: obj,
		})
	}
}

// AuthHandler godoc
//
//	@Summary		Пользователь по токену
//	@Tags			sessions
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	schema.User
//	@Failure		401		{object}	string
//	@Router			/me [get]
//	@Security JWT
func (s *Service) AuthHandler(c *gin.Context) {
	token, err := middleware.CheckIfTokenExpire(c, s.JWTSecret)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	claims := token.Claims.(jwt.MapClaims)

	var identity = uint(claims["user_id"].(float64))

	user, err := crud.FindUser(s.DB, "id = ?", identity)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
