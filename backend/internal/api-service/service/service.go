package service

import (
	"backend/internal/proto"
	"backend/internal/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	Router                 *gin.Engine
	Port                   int
	DB                     *gorm.DB
	ReCaptcha              util.ReCaptcha
	JWTSecret              []byte
	ImageProcessingService proto.ImageProcessingServiceClient
}
