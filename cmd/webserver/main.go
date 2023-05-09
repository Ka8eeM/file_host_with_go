package main

import (
	"os"

	itemService "github.com/Ka8eeM/file_host_with_go"
	db "github.com/Ka8eeM/file_host_with_go/db/files"
	dbMemory "github.com/Ka8eeM/file_host_with_go/db/memory"

	storage "github.com/Ka8eeM/file_host_with_go/storage/files"
	"github.com/gin-gonic/gin"
	kitlog "github.com/go-kit/log"
)

// RedisHost env var name
const RedisHost = "REDIS_HOST"

func main() {

	_ = dbMemory.Repo{}
	_ = db.Repo{}

	// init service
	loggerService := kitlog.With(kitlog.NewJSONLogger(os.Stdout), "ts", kitlog.DefaultTimestampUTC)
	storageService, err := storage.New(loggerService, "./uploads")

	if err != nil {
		loggerService.Log("message", "could not init storage service", "error", err)
		return
	}
	dbService, err := db.New(loggerService, "./dbFiles")
	if err != nil {
		loggerService.Log("message", "could not init db service", "error", err)
	}

	s := itemService.New(loggerService, storageService, dbService)

	r := gin.Default()
	r.Use(setItemService(s))

	r.GET("/ping", pong)
	r.GET("/i/:id", getItem)
	r.POST("/upload", upload)

	r.Run(":8080")
}

func setItemService(s *itemService.Srvs) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("itemService", s)
	}
}
