package main

import (
	"crypto/x509"
	"encoding/pem"
	idp "github.com/derekmckinnon/test-saml-idp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	certificatePath = "saml.crt"
	privateKeyPath  = "saml.key"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func main() {
	cert, key, err := idp.GenerateDevelopmentCertificateAndKey()
	if err != nil {
		log.Fatal().Err(err).Msg("error generating certificate and key")
	}

	certBytes := cert.Raw
	outputPem(certificatePath, "CERTIFICATE", certBytes)

	keyByes, _ := x509.MarshalPKCS8PrivateKey(key)
	outputPem(privateKeyPath, "RSA PRIVATE KEY", keyByes)
}

func outputPem(filename, format string, bytes []byte) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating file")
	}
	defer file.Close()

	err = pem.Encode(file, &pem.Block{
		Type:  format,
		Bytes: bytes,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error writing pem")
	}
}
