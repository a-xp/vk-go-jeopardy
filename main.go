package main

import (
	"github.com/gin-gonic/gin"
	"goj/configuration"
	"goj/domain"
	"goj/game"
	"goj/rating"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var AppConfig *configuration.Configuration

func main() {
	rand.Seed(time.Now().UnixNano())
	AppConfig = configuration.LoadConfigFile()
	domain.InitEngine(AppConfig)
	defer domain.StopEngine()
	gin.SetMode(AppConfig.Http.Mode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/api/reload", func(c *gin.Context) {
		domain.ReloadGames()
		c.Status(http.StatusAccepted)
	})
	r.POST("/api/callback", game.HandleVKEvent)
	rating.ConfigureAPI(r)
	if err := r.Run(AppConfig.Http.ListenAddr); err != nil {
		log.Fatal(err)
	}
}
