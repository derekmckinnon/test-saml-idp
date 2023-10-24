package main

import (
	"crypto/rsa"
	"crypto/x509"
	idp "github.com/derekmckinnon/test-saml-idp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("/etc/test-saml-idp/")
	viper.AddConfigPath(".")

	viper.SetDefault("Host", "http://localhost:8080")
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading configuration")
	}

	log.Info().Msg("Loading certificate and key")
	cert, key := loadCertificateAndKey(config)

	server := idp.New(idp.ServerOptions{
		Config:      config,
		Key:         key,
		Certificate: cert,
	})

	log.Info().Msg("Loading users")
	err = server.LoadUsers(config.Users)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading users")
	}

	log.Info().Msg("Loading services")
	err = server.LoadServices(config.Services)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading services")
	}

	log.Info().Msg("Starting server")
	if err = server.Run(); err != nil {
		log.Fatal().Err(err).Msg("error running server")
	}
}

func loadCertificateAndKey(config *idp.Config) (*x509.Certificate, *rsa.PrivateKey) {
	certPath, keyPath := config.CertificatePath, config.KeyPath

	if certPath == "" || keyPath == "" {
		log.Info().Msg("Generating development certificate")
		cert, key, err := idp.GenerateDevelopmentCertificateAndKey()
		if err != nil {
			log.Fatal().Err(err).Msg("could not generate development certificate")
		}

		return cert, key
	}

	log.Info().Str("path", config.CertificatePath).Msg("Loading certificate from the filesystem")
	cert, err := idp.LoadCertificatePem(config.CertificatePath)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading certificate")
	}

	log.Info().Str("path", config.KeyPath).Msg("Loading private key from the filesystem")
	key, err := idp.LoadPrivateKeyPem(config.KeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading private key")
	}

	return cert, key
}

func loadConfig() (*idp.Config, error) {
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
