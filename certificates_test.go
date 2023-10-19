package idp

import (
	"crypto/x509"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateDevelopmentCertificate(t *testing.T) {
	cert, err := GenerateDevelopmentCertificate()
	require.NoError(t, err)

	organization := cert.Subject.Organization
	require.Len(t, organization, 1)
	require.Equal(t, organization[0], "Test SAML IdP")
	require.Equal(t, cert.KeyUsage, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment)
}
