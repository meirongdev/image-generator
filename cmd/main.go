package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meirongdev/image-generator/internal/drawing"
)

func router() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
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
	err := router().Run(":8080")
	if err != nil {
		panic(err)
	}
}
