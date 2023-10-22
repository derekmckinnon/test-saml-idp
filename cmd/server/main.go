package main

import (
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
		log.Error().Err(err).Msg("error loading configuration")
	}

	log.Info().Msg("Generating development certificate")
	cert, key, err := idp.GenerateDevelopmentCertificate()
	if err != nil {
		log.Error().Err(err).Msg("could not generate development certificate")
	}

	log.Info().Msg("Initializing IdP Server...")
	server := idp.New(idp.ServerOptions{
		Config:      config,
		Key:         key,
		Certificate: cert,
	})

	err = server.LoadUsers(config.Users)
	if err != nil {
		log.Error().Err(err).Msg("cannot initialize users")
	}

	err = server.LoadServices(config.Services)
	if err != nil {
		log.Error().Err(err).Msg("cannot initialize services")
	}

	log.Info().Msg("Starting IdP Server...")
	if err = server.Run(); err != nil {
		log.Error().Err(err).Msg("cannot run server")
	}
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
