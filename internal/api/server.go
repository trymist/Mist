package api

import (
	"github.com/gin-gonic/gin"
)

func StartServer() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//TODO: embed static files in binary for final distribution
	r.Static("/public", "./public")
	r.StaticFile("/favicon.ico", "./public/mist.png")
	registerRoutes(r)
	r.Run(":7777")
}

func registerRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api", gin.Logger())
	viewsGroup := r.Group("/")
	registerApiRoutes(apiGroup)
	registerViewsRoutes(viewsGroup)
}
