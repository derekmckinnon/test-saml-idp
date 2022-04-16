package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.tmpl", gin.H{})
	})

	r.POST("/login", func(c *gin.Context) {
		var form LoginForm
		c.Bind(&form)
		c.HTML(200, "login.tmpl", gin.H{
			"Username": form.Username,
			"Error":    "Unknown username or password",
		})
	})

	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
