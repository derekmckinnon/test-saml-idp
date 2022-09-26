package idp

import (
	"github.com/gomarkdown/markdown"
	"html/template"
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

func renderMarkdown(md string) template.HTML {
	raw := []byte(md)
	output := markdown.ToHTML(raw, nil, nil)

	return template.HTML(output)
}
