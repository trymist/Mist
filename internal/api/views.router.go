package api

import (
	"context"
	"mist/internal/auth"
	"mist/views/pages"

	"github.com/gin-gonic/gin"
)

func registerViewsRoutes(r *gin.RouterGroup) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/register")
	})
	r.GET("/register", auth.RegisterView())
	r.GET("/login", auth.LoginView())
	r.GET("/home", func(c *gin.Context) {
		pages.Home().Render(context.Background(), c.Writer)
	})
}
