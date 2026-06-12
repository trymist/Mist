package auth

import (
	"mist/views/pages"

	"github.com/gin-gonic/gin"
)

func RegisterView() func(c *gin.Context) {
	return func(c *gin.Context) {
		isFirstUser, err := IsFirstUser()
		if err != nil {
			pages.Error(err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}
		if isFirstUser {
			pages.Register().Render(c.Request.Context(), c.Writer)
			return
		}
		c.Redirect(302, "/login")
	}
}

func LoginView() func(c *gin.Context) {
	return func(c *gin.Context) {
		isFirstUser, err := IsFirstUser()
		if err != nil {
			pages.Error(err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}
		if isFirstUser {
			c.Redirect(302, "/register")
			return
		}
		pages.Login().Render(c.Request.Context(), c.Writer)
	}
}
