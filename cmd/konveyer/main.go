package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
)

// Configuration - Application Runtime Configuration
type Configuration struct {
	Source string `envconfig:"source" required:"true"`
	Debug  bool   `envconfig:"debug" default:"false"`
}

// SetupRouter - Locations for REST Server
func SetupRouter() *gin.Engine {
	var s Configuration
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	router := gin.Default()

	v1 := router.Group("api/v1")
	{
		v1.GET("/", index)
		v1.GET("/status", s.status)
		v1.GET("/deploy/:namespace/:deployment/:newImage", s.updateImage)
	}

	return router
}

func index(c *gin.Context) {
	c.JSON(200, gin.H{"ok": "GET /"})
}

func main() {
	router := SetupRouter()
	router.Run(":5000")
}
