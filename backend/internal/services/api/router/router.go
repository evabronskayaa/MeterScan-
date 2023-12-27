package router

import (
	_ "backend/docs"
	"backend/internal/services/api/endpoint"
	"backend/internal/services/api/middleware"
	"backend/internal/services/api/service"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"
)

// ConfigureRouter
//
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
func ConfigureRouter(s *service.Service, port int) *http.Server {
	router := gin.Default()

	router.RedirectTrailingSlash = true
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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
	v1.GET("/verify", endpointSever.VerifyHandler)

	authorized := v1.Group("/", middlewareServer.AuthMiddleware())
	{
		authorized.POST("/verify", endpointSever.RequestVerifyHandler)

		verified := authorized.Group("", middlewareServer.VerifyMiddleware())
		{
			verified.GET("media/:dir/*asset", endpointSever.MediaHandler)

			predictions := verified.Group("/predictions")
			{
				predictions.GET("", endpointSever.GetPredictionsHandler)
				predictions.POST("", endpointSever.PredictHandler)
				predictions.PUT("", endpointSever.UpdatePredictHandler)
				predictions.DELETE("", endpointSever.RemovePredictHandler)
			}

			settings := verified.Group("/settings")
			{
				settings.GET("", endpointSever.GetSettingsHandler)
				settings.PUT("/notification", endpointSever.SetNotificationHandler)
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: router,
	}
}
