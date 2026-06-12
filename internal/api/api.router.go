package api

import (
	"mist/internal/auth"
	"time"

	"github.com/gin-gonic/gin"
)

func registerApiRoutes(r *gin.RouterGroup) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now().Local().Format(time.DateTime),
		})
	})
	r.POST("/register", auth.RegisterHandler())
	r.POST("/login", auth.LoginHandler())
}
