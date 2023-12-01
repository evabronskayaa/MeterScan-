package endpoint

import (
	"backend/internal/dto"
	"backend/internal/errors"
	"backend/internal/proto"
	"backend/internal/services/api/middleware"
	"backend/internal/services/api/service"
	"backend/internal/util"
	errors2 "errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"time"
)

const (
	ErrIncorrectPassword errors.SimpleError = "Неверный пароль"
)

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

	if err := loginForm.Validate(dto.ValidateArgs{Ctx: c, ReCaptcha: s.ReCaptcha}); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	email := loginForm.Email
	password := loginForm.Password

	user, err := s.DatabaseService.GetUser(c, &proto.UserRequest{Email: &email})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if !util.CheckPasswordHash(user.Password, password) {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrIncorrectPassword)
		return
	}

	obj, err := middleware.GenerateToken(s.JWTSecret, func(claims jwt.MapClaims) {
		claims["user_id"] = user.Id
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
	} else if err := registerForm.Validate(dto.ValidateArgs{Ctx: c, ReCaptcha: s.ReCaptcha}); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	if user, err := s.DatabaseService.CreateUser(c, &proto.UserRequest{Email: &registerForm.Email, Password: &registerForm.Password}); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	} else {
		if err := s.sendRequestVerification(user); err != nil {
			log.Printf("Request verification with err: %v", err)
		}

		obj, err := middleware.GenerateToken(s.JWTSecret, func(claims jwt.MapClaims) {
			claims["user_id"] = user.Id
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
//	@Success		200		{object}	proto.UserResponse
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

	var identity = uint64(claims["user_id"].(float64))

	user, err := s.DatabaseService.GetUser(c, &proto.UserRequest{Id: &identity})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// VerifyHandler godoc
//
//	@Summary		Подтвердить аккаунт (сообщение на почте)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param token query dto.VerifyTokenForm true "token"
//	@Success		200
//	@Failure		400		{object}	string
//	@Failure		401		{object}	string
//	@Failure		500		{object}	string
//	@Router			/verify [get]
func (s *Service) VerifyHandler(c *gin.Context) {
	var vToken dto.VerifyTokenForm
	if err := c.ShouldBindQuery(&vToken); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	token, err := jwt.Parse(vToken.Token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(middleware.SigningAlgorithm) != t.Method {
			return nil, middleware.ErrInvalidSigningAlgorithm
		}

		return s.JWTSecret, nil
	})
	if err != nil {
		var validationErr *jwt.ValidationError
		ok := errors2.As(err, &validationErr)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	_, ok := claims["email"]
	if !ok || int64(claims["exp"].(float64)) < time.Now().Unix() {
		c.AbortWithStatusJSON(http.StatusUnauthorized, middleware.ErrExpiredToken)
		return
	}

	identity := uint64(claims["user_id"].(float64))

	if _, err := s.DatabaseService.VerifyUser(c, &proto.UserRequest{Id: &identity}); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	} else {
		c.Status(http.StatusOK)
	}
}

// RequestVerifyHandler godoc
//
//	@Summary		Попросить выслать новое письмо для подтверждения
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		405		{object}	string
//	@Failure		500		{object}	string
//	@Router			/verify [post]
//	@Security JWT
func (s *Service) RequestVerifyHandler(c *gin.Context) {
	user := c.MustGet("user").(*proto.UserResponse)

	if user.Verified {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	if err := s.sendRequestVerification(user); err != nil {
		log.Printf("Request verification with err: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
	} else {
		c.Status(http.StatusOK)
	}
}

func (s *Service) sendRequestVerification(user *proto.UserResponse) error {
	if _, err := middleware.GenerateToken(s.JWTSecret, func(claims jwt.MapClaims) {
		claims["user_id"] = user.Id
		claims["email"] = user.Email
	}); err != nil {
		return err
	} else {
		return nil
		//return user.RequestVerification(s.MailClient, generateToken.Token)
	}
}
