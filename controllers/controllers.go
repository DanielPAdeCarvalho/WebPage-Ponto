package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	globals "loginpage/globals"
	helpers "loginpage/helpers"
)

func LoginGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html",
				gin.H{
					"content": "Please logout first",
					"user":    user,
				})
			return
		}
		c.HTML(http.StatusOK, "login.html", gin.H{
			"content": "",
			"user":    user,
		})
	}
}

func LoginPostHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if user != nil {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Please logout first"})
			return
		}

		username := c.PostForm("username")
		password := c.PostForm("password")

		if helpers.EmptyUserPass(username, password) {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"content": "Parameters can't be empty"})
			return
		}

		if !helpers.CheckUserPass(username, password) {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"content": "Incorrect username or password"})
			return
		}

		session.Set(globals.Userkey, username)
		if err := session.Save(); err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"content": "Failed to save session"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, "/dashboard")
	}
}

func LogoutGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		log.Println("logging out user:", user)
		if user == nil {
			log.Println("Invalid session token")
			return
		}
		session.Clear()
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}

func IndexGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		c.HTML(http.StatusOK, "index.html", gin.H{
			"content": "This is an index page...",
			"user":    user,
		})
	}
}

// Se o batePonto for true ele vai bater o ponto do funcionario e depois atualizar a pagina
func DashboardGetHandler(batePonto bool) gin.HandlerFunc {
	texto := "Ultimos Pontos:"
	responseBatePonto := ""

	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		if batePonto {
			helpers.BatePonto(user.(string))
			responseBatePonto = "Ponto Batido com Sucesso."
		}
		pontos := helpers.UltimosPontos()
		t, _ := time.Parse("2006-01-02_15:04:05", pontos[0].Data)
		bianca := "Bianca: " + t.Format("02/01/2006 - 15:04")
		t, _ = time.Parse("2006-01-02_15:04:05", pontos[1].Data)
		danilo := "Danilo: " + t.Format("02/01/2006 - 15:04")
		t, _ = time.Parse("2006-01-02_15:04:05", pontos[2].Data)
		paty := "Patrícia: " + t.Format("02/01/2006 - 15:04")
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"texto":             texto,
			"danilo":            danilo,
			"paty":              paty,
			"bianca":            bianca,
			"user":              user,
			"responsebateponto": responseBatePonto,
		})
	}
}

func CadastroGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)
		c.HTML(http.StatusOK, "cadastro.html", gin.H{
			"user": user,
		})
	}
}

// Cadastrar um usuario novo
func SigninGetHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get(globals.Userkey)

		//Pegar as coisas do uusario pra gravar
		username := c.PostForm("username")
		cpf := c.PostForm("cpf")
		datanascimento := c.PostForm("datanascimento")
		nomecompleto := c.PostForm("nomecompleto")
		password := c.PostForm("password")
		helpers.Cadastro(username, cpf, datanascimento, nomecompleto, password)

		c.HTML(http.StatusOK, "cadastro.html", gin.H{
			"user":       user,
			"cadastrado": "Usuário Cadastrado com sucesso.",
		})
		//TODO ajustar o mapeamento para cadastrar os itens com o mapeamento json correto
	}
}
