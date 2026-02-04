package widget

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var widgetFS embed.FS

/*
 * RegisterRoutes registers the widget routes
 * Serves the compiled widget JavaScript at /widget/*
 */
func RegisterRoutes(r *gin.Engine) error {
	// Get the dist subdirectory
	dist, err := fs.Sub(widgetFS, "dist")
	if err != nil {
		return err
	}

	// Serve widget files at /widget/
	r.StaticFS("/widget", http.FS(dist))

	return nil
}
