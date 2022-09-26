package idp

import (
	"net/url"
	"strings"
)

func getBasePath(u url.URL) string {
	basePath := u.Path

	if basePath == "" || basePath == "/" {
		return "/"
	}

	return strings.TrimSuffix(basePath, "/")
}
