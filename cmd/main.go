package main

import (
	"github.com/gin-gonic/gin"
	"github.com/meirongdev/image-generator/internal/drawing"
)

func router() *gin.Engine {
	r := gin.Default()
	userRoute := r.Group("/image")
	{
		userRoute.GET("/circles", func(c *gin.Context) {
			file := drawing.DrawOne("circles")
			c.Header("Content-Type", "image/png")
			c.File(file)
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
