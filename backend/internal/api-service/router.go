package api_service

import (
	_ "backend/docs"
	"backend/internal/api-service/endpoint"
	"backend/internal/api-service/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

// @title 		MeterScanPlus API
// @version		1.0
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @description Пример: `Bearer *token*`
func (s *Service) configureRouter() *gin.Engine {
	router := gin.Default()

	router.RedirectTrailingSlash = true
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))

	v1 := router.Group("/api/v1")
	endpointSever := (endpoint.Service)(*s)
	v1.POST("/users", endpointSever.RegisterHandler)
	v1.POST("/sessions", endpointSever.LoginHandler)
	v1.GET("/me", endpointSever.AuthHandler)

	middlewareServer := (middleware.Service)(*s)
	v1.GET("/refresh", middlewareServer.RefreshToken)

	authorized := v1.Group("/", middlewareServer.AuthMiddleware())
	{
		authorized.GET("media/:dir/*asset", endpointSever.MediaHandler)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
