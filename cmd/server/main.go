package main

import (
	idp "github.com/derekmckinnon/test-saml-idp"
	"github.com/spf13/viper"
	"log"
)

func main() {
	log.SetPrefix("[IdP] ")

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	log.Println("Generating development certificate")
	cert, key, err := idp.GenerateDevelopmentCertificate()
	if err != nil {
		log.Fatalf("could not generate development certificate: %v", err)
	}

	log.Println("Initializing IdP Server...")
	server := idp.New(idp.ServerOptions{
		Config:      config,
		Key:         key,
		Certificate: cert,
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

func loadConfig() (*idp.Config, error) {
	viper.SetDefault("Host", "http://localhost:8080")

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/etc/test-saml-idp/")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.BindEnv("Host", "HOST")
	if err != nil {
		return nil, err
	}

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &idp.Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
