package main

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	idp "github.com/derekmckinnon/test-saml-idp"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
)

func main() {
	log.SetPrefix("[IdP] ")

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	baseUrl, err := url.Parse(config.Host)
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

func loadPemFile(path string) (*pem.Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)

	return block, nil
}

func loadPrivateKey() (crypto.PrivateKey, error) {
	block, err := loadPemFile("./idp_key.pem")
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS8PrivateKey(block.Bytes)
}

func loadCertificate() (*x509.Certificate, error) {
	block, err := loadPemFile("./idp_cert.pem")
	if err != nil {
		return nil, err
	}

	return x509.ParseCertificate(block.Bytes)
}
