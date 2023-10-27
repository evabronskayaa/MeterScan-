package service

import (
	"backend/internal/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	Router    *gin.Engine
	DB        *gorm.DB
	ReCaptcha util.ReCaptcha
	JWTSecret []byte
}
