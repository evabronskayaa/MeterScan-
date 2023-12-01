package endpoint

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

// MediaHandler godoc
//
//	@Summary		Достать файл
//	@Tags			files
//	@Param dir path string true "Директория"
//	@Param asset  path string true "Файл"
//	@Success		200 "file"
//	@Failure		404
//	@Router			/media/{dir}/{asset} [get]
//	@Security JWT
func (s *Service) MediaHandler(c *gin.Context) {
	dir := c.Param("dir")
	asset := c.Param("asset")
	if strings.TrimPrefix(asset, "/") == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	fullName := filepath.Join("./media", dir, filepath.FromSlash(path.Clean("/"+asset)))
	c.File(fullName)
}
