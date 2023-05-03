package main

import (
	globals "loginpage/globals"
	middleware "loginpage/middleware"
	routes "loginpage/routes"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set the router as the default one shipped with Gin
	gin.SetMode(gin.DebugMode)
	router := gin.New()

	router.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	public := router.Group("/")
	routes.PublicRoutes(public)

	private := router.Group("/")

	// Additional configs for the gin sessions middleware and the AuthRequired middleware
	private.Use(middleware.AuthRequired)
	routes.PrivateRoutes(private)
	router.SetTrustedProxies(nil)

	router.Run(":8080")
}
