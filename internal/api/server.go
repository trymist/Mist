package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()
	registerRoutes(r)
	r.Run(":7777")
}

func registerRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Local().Format(time.DateTime),
		})
	})

}
