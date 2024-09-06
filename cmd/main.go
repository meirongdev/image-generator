package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/meirongdev/image-generator/internal/drawing"
	"github.com/meirongdev/image-generator/internal/pool"
)

func router(basePath string, taskPool *pool.TaskPool) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob(path.Join(basePath, "templates/*.html"))
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
	requestGroup := r.Group("/request")
	{
		type ReqTrace struct {
			ID  string `json:"id"`
			Url string `json:"url"`
		}

		getResultUrl := func(req *http.Request, id string) string {
			scheme := "http"
			if req.TLS != nil {
				scheme = "https"
			}
			return scheme + "://" + req.Host + "/request/res/" + id
		}

		requestGroup.GET("res/:id", func(c *gin.Context) {
			id := c.Param("id")
			p, done := taskPool.GetResult(id)
			if done {
				log.Printf("Found %s for id: %s\n", p, id)
				c.Header("Content-Type", "image/png")
				c.File(p)
			} else {
				log.Printf("Path not found for id: %s\n", id)
				c.Header("Content-Type", "image/jpg")
				c.Header("Cache-Control", "no-cache")
				c.File(path.Join(basePath, "static/loading.jpg"))
			}
		})

		requestGroup.GET(":name", func(c *gin.Context) {
			name := c.Param("name")
			id, err := gonanoid.New()
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
			}
			task := &pool.Task{
				ID:        id,
				ImageName: name,
			}
			taskPool.SubmitTask(task)
			url := getResultUrl(c.Request, id)
			res := ReqTrace{
				ID:  id,
				Url: url,
			}
			c.JSON(http.StatusOK, res)
		})
	}
	return r
}

func main() {
	basePath := "./"
	if envPath := os.Getenv("BASEPATH"); envPath != "" {
		basePath = envPath
	}
	taskPool := pool.NewTaskPool(2, 10)
	go taskPool.Start()

	go func() {
		err := router(basePath, taskPool).Run(":8080")
		if err != nil {
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	fmt.Println("Received stop signal, shutting down...")
	taskPool.Stop()
}
