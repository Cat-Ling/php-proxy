package middleware

import (
	"net/http"
	"strings"
)

import (
	"net/http"
	"net/url"
	"strings"
)

func BlockList(blockedURLs []string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			targetURL := r.URL.Query().Get("url")
			if targetURL != "" {
				parsedURL, err := url.Parse(targetURL)
				if err == nil {
					for _, blockedURL := range blockedURLs {
						if strings.Contains(parsedURL.Host, blockedURL) {
							http.Error(w, "Access to this URL is blocked", http.StatusForbidden)
							return
						}
					}
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}