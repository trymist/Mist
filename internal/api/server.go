package api

import (
	"context"
	"fmt"
	"mist/views"
	"mist/views/pages"
	"time"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()
	r.Static("/public", "./public")
	r.StaticFile("/favicon.ico", "./public/mist.png")
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
	r.GET("/hello", func(c *gin.Context) {
		component := views.Hello("bablu")
		component.Render(context.Background(), c.Writer)
	})
	r.GET("/register", func(c *gin.Context) {
		component := pages.Register()
		component.Render(context.Background(), c.Writer)
	})
	r.POST("/register", func(c *gin.Context) {
		fullName := c.PostForm("fullName")
		userName := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")
		confirmPassword := c.PostForm("passwordConfirm")
		fmt.Println(fullName, userName, email, password, confirmPassword)
		time.Sleep(5 * time.Second)
		c.String(200, "Failed to register user")
	})

	protected := r.Group("/api")
	protected.Use(authMiddleware())
	protected.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": "John Doe",
			"age":  30,
		})
	})

}
