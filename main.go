package main

import (
	"log"
	"net/url"
)

func main() {
	baseUrl, err := url.Parse("http://localhost:8080")
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

	log.Println("Initializing IdP Server...")
	server := InitServer(ServerOptions{
		BaseUrl:     *baseUrl,
		Key:         key,
		Certificate: certificate,
	})

	log.Println("Starting IdP Server...")
	if err = server.Run(); err != nil {
		log.Fatalf("cannot run server: %v", err)
	}
}
