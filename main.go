package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"github.com/derekmckinnon/test-saml-idp/idpserver"
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"os"
)

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func main() {
	log.Println("Initializing IdP Server...")

	key, err := loadKey()
	if err != nil {
		log.Fatalln("Could not load private key", err)
	}

	certificate, err := loadCertificate()
	if err != nil {
		log.Fatalln("Could not load certificate", err)
	}

	server := idpserver.New(idpserver.Options{
		BaseUrl:     url.URL{},
		Key:         key,
		Certificate: certificate,
	})

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")

	router.GET("/metadata", func(c *gin.Context) {
		metadata := server.GetMetadata()
		c.XML(200, metadata)
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.tmpl", gin.H{})
	})

	router.POST("/login", func(c *gin.Context) {
		var form LoginForm
		_ = c.Bind(&form)
		c.HTML(200, "login.tmpl", gin.H{
			"Username": form.Username,
			"Error":    "Unknown username or password",
		})
	})

	log.Println("Starting IdP Server...")
	err = router.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func loadKey() (crypto.PrivateKey, error) {
	data, err := os.ReadFile("./testdata/idp_key.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	return key, nil
}

func loadCertificate() (*x509.Certificate, error) {
	data, err := os.ReadFile("./testdata/idp_cert.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	certificate, _ := x509.ParseCertificate(block.Bytes)

	return certificate, nil
}
