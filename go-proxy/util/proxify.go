package util

import (
	"net/url"
	"strings"
)

func ProxifyURL(targetURL string, baseURL *url.URL) string {
	if strings.HasPrefix(targetURL, "data:") {
		return targetURL
	}

	absURL, err := baseURL.Parse(targetURL)
	if err != nil {
		return targetURL
	}

	return "/?url=" + url.QueryEscape(absURL.String())
}