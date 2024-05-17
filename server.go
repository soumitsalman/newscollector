package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var collector NewsSiteCollector

func collectHandler(ctx *gin.Context) {
	go collector.Collect()
	ctx.JSON(http.StatusOK, "collection started")
}

func createRateLimitHandler(r rate.Limit, b int) gin.HandlerFunc {
	rate_limiter := rate.NewLimiter(r, b)
	return func(ctx *gin.Context) {
		if rate_limiter.Allow() {
			ctx.Next()
		} else {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
		}
	}
}

func newServer(r rate.Limit, b int) *gin.Engine {
	router := gin.Default()
	auth_group := router.Group("/")
	auth_group.Use(createRateLimitHandler(r, b)) // authn and ratelimit middleware
	auth_group.POST("/collect", collectHandler)  // routes
	return router
}

func main() {
	godotenv.Load()
	switch os.Getenv("INSTANCE_MODE") {
	case "DEBUG":
		start_time := time.Now()
		// initialize to save locally
		collector = NewCollector(localFileStore)
		collector.Collect()
		log.Println("Collection took", time.Since(start_time))

	default:
		// initialize collector to save remotely
		collector = NewCollector(remoteStoreBeans)
		// initialize and run server
		newServer(2, 5).Run()
	}

}
