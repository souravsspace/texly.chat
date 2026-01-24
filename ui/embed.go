package ui

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed all:dist
var buildFS embed.FS

func RegisterRoutes(r *gin.Engine) error {
	// Get the build subdirectory
	build, err := fs.Sub(buildFS, "dist")
	if err != nil {
		return err
	}

	// Serve static files
	r.StaticFS("/", http.FS(build))

	return nil
}
