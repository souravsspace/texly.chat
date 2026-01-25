package ui

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"

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

	// Serve static files via NoRoute to avoid Catch-All conflict with /api
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Do not handle API routes here
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		// Check if file exists in build fs
		// Trim leading slash for fs.Open
		filePath := strings.TrimPrefix(path, "/")
		f, err := build.Open(filePath)
		if err == nil {
			defer f.Close()
			stat, _ := f.Stat()
			if !stat.IsDir() {
				// Serve the specific file
				c.FileFromFS(path, http.FS(build))
				return
			}
		}

		// Fallback to index.html for SPA
		indexFile, err := build.Open("index.html")
		if err != nil {
			c.String(http.StatusNotFound, "index.html not found")
			return
		}
		defer indexFile.Close()

		indexContent, err := io.ReadAll(indexFile)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to read index.html")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", indexContent)
	})

	return nil
}
