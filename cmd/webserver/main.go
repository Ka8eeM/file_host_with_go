package main

import (
	"os"

	db "db/files"
	dbMemory "db/memory"
	storage "storage/files"

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
	storageService, err := storage.NewService(loggerService, "./uploads")

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
