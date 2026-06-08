package api

import "github.com/gin-gonic/gin"

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("mist-token")
		if err != nil {
			c.AbortWithStatus(401)
			return
		}
		if token == "" {
			c.AbortWithStatus(401)
			return
		}
	}
}
