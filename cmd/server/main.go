package main

import (
	idp "github.com/derekmckinnon/test-saml-idp"
	"log"
	"net/url"
)

func main() {
	log.SetPrefix("[IdP] ")

	port := getEnvOrDefault("PORT", "8080")
	baseUrlStr := getEnvOrDefault("BASE_URL", "http://localhost:"+port)

	baseUrl, err := url.Parse(baseUrlStr)
	if err != nil {
		log.Fatalf("cannot parse base URL: %v", err)
	}

	key, err := loadPrivateKey()
	if err != nil {
		log.Fatalf("cannot load private key: %v", err)
	}

	certificate, err := loadCertificate()
	if err != nil {
		log.Fatalf("cannot load certificate: %v", err)
	}

	config := idp.Config{}
	err = config.InitFromFile("./config.yml")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	log.Println("Initializing IdP Server...")
	server := idp.New(idp.ServerOptions{
		BaseUrl:     *baseUrl,
		Key:         key,
		Certificate: certificate,
	})

	err = server.LoadUsers(config.Users)
	if err != nil {
		log.Fatalf("cannot initialize users: %v", err)
	}

	err = server.LoadServices(config.Services)
	if err != nil {
		log.Fatalf("cannot initialize services: %v", err)
	}

	log.Println("Starting IdP Server...")
	if err = server.Run(); err != nil {
		log.Fatalf("cannot run server: %v", err)
	}
}
