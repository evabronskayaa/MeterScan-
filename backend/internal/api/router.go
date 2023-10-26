package api

import (
	"github.com/gin-gonic/gin"
)

func configureRouter() *gin.Engine {
	router := gin.Default()

	router.RedirectTrailingSlash = true

	return router
}
