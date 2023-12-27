package middleware

import (
	"backend/internal/dto"
	"backend/internal/errors"
	"backend/internal/proto"
	"backend/internal/services/api/service"
	errors2 "errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valyala/fasthttp"
	"net/http"
	"strings"
	"time"
)

const (
	ErrEmptyAuthHeader         errors.SimpleError = "Отсутствует токен авторизации"
	ErrInvalidAuthHeader       errors.SimpleError = "Некорректный токен авторизации"
	ErrInvalidSigningAlgorithm errors.SimpleError = "Неверный алгоритм подписи"
	ErrFailedTokenCreation     errors.SimpleError = "Ошибка при создании токена"
	ErrExpiredToken            errors.SimpleError = "Токен просрочен"
)

const (
	SigningAlgorithm = "HS256"
	timeout          = time.Hour * 1000
	authHeader       = "Authorization"
	tokenHeadName    = "Bearer"
)

type Service service.Service

func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := CheckIfTokenExpire(c, s.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		identity := uint64(claims["user_id"].(float64))

		user, err := s.DatabaseService.GetUser(c, &proto.UserRequest{Id: &identity})

		if err == nil {
			c.Set("user", user)
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, err)
			return
		}

		c.Next()
	}
}

func parseToken(c *gin.Context, key []byte) (*jwt.Token, error) {
	authHeader := c.Request.Header.Get(authHeader)

	if authHeader == "" {
		return nil, ErrEmptyAuthHeader
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == tokenHeadName) {
		return nil, ErrInvalidAuthHeader
	}

	token := parts[1]

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}

		c.Set("JWT_TOKEN", token)

		return key, nil
	})
}

// RefreshToken godoc
//
//	@Summary		Обновление токена
//	@Tags			sessions
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dto.Token
//	@Failure		401		{object}	string
//	@Router			/refresh [get]
//	@Security JWT
func (s *Service) RefreshToken(c *gin.Context) {
	token, err := CheckIfTokenExpire(c, s.JWTSecret)
	if err != nil {
		c.AbortWithStatusJSON(fasthttp.StatusUnauthorized, err)
		return
	}

	obj, err := GenerateToken(s.JWTSecret, func(newClaims jwt.MapClaims) {
		claims := token.Claims.(jwt.MapClaims)
		for key := range claims {
			newClaims[key] = claims[key]
		}
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, obj)
}

func CheckIfTokenExpire(c *gin.Context, key []byte) (*jwt.Token, error) {
	token, err := parseToken(c, key)
	if err != nil {
		var validationErr *jwt.ValidationError
		ok := errors2.As(err, &validationErr)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}

	claims := token.Claims.(jwt.MapClaims)

	if int64(claims["exp"].(float64)) < time.Now().Unix() {
		return nil, ErrExpiredToken
	}

	return token, nil
}

func GenerateToken(key []byte, genClaims func(claims jwt.MapClaims)) (*dto.Token, error) {
	newToken := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	genClaims(newClaims)

	expire := time.Now().Add(timeout)
	origIat := time.Now()
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = origIat.Unix()
	token, err := newToken.SignedString(key)
	if err != nil {
		return nil, ErrFailedTokenCreation
	}

	return &dto.Token{
		Token:   token,
		Expire:  expire.Format(time.RFC3339),
		OrigIat: origIat.Format(time.RFC3339),
	}, nil
}
