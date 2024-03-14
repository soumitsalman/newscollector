package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	datautils "github.com/soumitsalman/data-utils"
	"github.com/soumitsalman/document-loader/loaders"
	"golang.org/x/time/rate"
)

const _SITEMAPS_CSV = "./sitemaps.csv"

func readSitemapsCSV() [][]string {
	file, _ := os.Open(_SITEMAPS_CSV)
	defer file.Close()
	sitemaps, _ := csv.NewReader(file).ReadAll()
	// ignore the header
	return sitemaps[1:]
}

func createSiteLoaders() []*loaders.WebLoader {
	site_loaders := datautils.Transform(readSitemapsCSV(), func(item *[]string) *loaders.WebLoader {
		return loaders.NewDefaultNewsSitemapLoader(1, (*item)[0])
	})
	return append(site_loaders,
		// this is a specialied loader
		loaders.NewYCHackerNewsSiteLoader(),
	)
}

func collectHandler(ctx *gin.Context) {
	site_loaders := createSiteLoaders()
	go datautils.ForEach(site_loaders, func(loader **loaders.WebLoader) {
		beans := (**loader).LoadSite()
		log.Println(len(beans), "new beans found")
		storeNewBeans(beans)
	})
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
	newServer(2, 5).Run()
}
