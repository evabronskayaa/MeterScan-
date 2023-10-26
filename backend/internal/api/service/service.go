package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	Router *gin.Engine
	DB     *gorm.DB
}
