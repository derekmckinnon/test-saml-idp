package idp

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func Test_GetBasePath(t *testing.T) {
	testCases := []struct {
		baseUrl  string
		basePath string
	}{
		{"https://idp.test.com", "/"},
		{"https://idp.test.com/", "/"},
		{"https://idp.test.com:8080", "/"},
		{"https://idp.test.com:8080/", "/"},
		{"https://idp.test.com/foobar", "/foobar"},
		{"https://idp.test.com/foobar/", "/foobar"},
		{"https://idp.test.com:8080/foobar", "/foobar"},
		{"https://idp.test.com:8080/foobar/", "/foobar"},
	}

	for _, testCase := range testCases {
		baseUrl, _ := url.Parse(testCase.baseUrl)

		expected := testCase.basePath
		actual := getBasePath(*baseUrl)

		require.Equal(t, expected, actual)
	}
}
