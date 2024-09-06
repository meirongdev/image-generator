package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/meirongdev/image-generator/internal/drawing"
)

func router(templatePath string) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob(templatePath)
	imageGroup := r.Group("/image")
	{
		imageGroup.GET("/:name", func(c *gin.Context) {
			name := c.Param("name")
			file := drawing.DrawOne(name)
			c.Header("Content-Type", "image/png")
			c.File(file)
		})
	}
	listGroup := r.Group("/list")
	{
		var keys []string
		for k := range drawing.DRAWINGS {
			keys = append(keys, k)
		}
		listGroup.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "bootstrap5.html", gin.H{
				"keys": keys,
			})
		})
	}
	return r
}

func main() {
	templatePath := "templates/*.html"
	if envPath := os.Getenv("TEMPLATE_PATH"); envPath != "" {
		templatePath = envPath
	}
	err := router(templatePath).Run(":8080")
	if err != nil {
		panic(err)
	}
}
