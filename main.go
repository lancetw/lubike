package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/lancetw/lubike/utils/ubikeutil"
)

func lubikeCommonEndpoint(c *gin.Context) {
	lat := c.Query("lat")
	lng := c.Query("lng")
	num := 2
	result, errno := ubikeutil.LoadNearbyUbikes(lat, lng, num)

	c.JSON(http.StatusOK, gin.H{
		"code":   errno,
		"result": result,
	})
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//var err error
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	store := persistence.NewInMemoryStore(time.Minute)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")
	{
		v1.GET("ubike-station/taipei", cache.CachePage(store, 10*time.Second, lubikeCommonEndpoint))
	}

	router.Run(":" + port)
}
