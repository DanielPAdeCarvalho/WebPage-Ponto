package routes

import (
	"github.com/gin-gonic/gin"

	controllers "loginpage/controllers"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.GET("/login", controllers.LoginGetHandler())
	g.POST("/login", controllers.LoginPostHandler())
	g.GET("/", controllers.IndexGetHandler())

}

// Rotas q so pode usar depois do logon
func PrivateRoutes(g *gin.RouterGroup) {
	g.GET("/dashboard", controllers.DashboardGetHandler(false))
	g.GET("/logout", controllers.LogoutGetHandler())
	g.GET("/cadastro", controllers.CadastroGetHandler())
	g.POST("/signin", controllers.SigninGetHandler())
	g.POST("/baterponto", controllers.DashboardGetHandler(true))

}
