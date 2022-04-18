package main

import (
	"github.com/crewjam/saml/samlidp"
	"golang.org/x/crypto/bcrypt"
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

	config := Config{}
	err = config.InitFromFile("./config.yml")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	log.Println("Initializing IdP Server...")
	server := NewServer(ServerOptions{
		BaseUrl:     *baseUrl,
		Key:         key,
		Certificate: certificate,
	})

	err = initializeUsers(server.Store, config.Users)
	if err != nil {
		log.Fatalf("cannot initialize users: %v", err)
	}

	log.Println("Starting IdP Server...")
	if err = server.Run(); err != nil {
		log.Fatalf("cannot run server: %v", err)
	}
}

func initializeUsers(store *Store, users []User) error {
	for _, user := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

		err := store.AddUser(&samlidp.User{
			Name:           user.Username,
			Email:          user.Email,
			HashedPassword: hashedPassword,
			GivenName:      user.FirstName,
			Surname:        user.LastName,
		})

		if err != nil {
			return err
		}

		log.Printf("Initialized User: %s\n", user.Username)
	}

	return nil
}
