package idp

import (
	"crypto/x509"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateDevelopmentCertificate(t *testing.T) {
	cert, key, err := GenerateDevelopmentCertificate()
	require.NoError(t, err)
	require.NotNil(t, cert)
	require.NotNil(t, key)

	if assert.Len(t, cert.Subject.Organization, 1) {
		assert.Equal(t, cert.Subject.Organization[0], "Test SAML IdP")
	}
	assert.Equal(t, cert.KeyUsage, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment)
}
