package middleware

import (
	"net/http"
)

import (
	"go-proxy/util"
	"net/http"
	"net/url"
	"strings"
)

func Cookie() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Rewrite request cookies
			if cookieHeader := r.Header.Get("Cookie"); cookieHeader != "" {
				r.Header.Del("Cookie")
				cookies := strings.Split(cookieHeader, "; ")
				for _, cookie := range cookies {
					parts := strings.SplitN(cookie, "=", 2)
					if len(parts) == 2 {
						name := parts[0]
						value := parts[1]
						if strings.HasPrefix(name, "pc_") {
							nameParts := strings.SplitN(name, "__", 2)
							if len(nameParts) == 2 {
								domain := strings.ReplaceAll(nameParts[0][3:], "_", ".")
								cookieName := nameParts[1]
								targetURL, _ := url.Parse(r.URL.Query().Get("url"))
								if strings.HasSuffix(targetURL.Host, domain) {
									r.AddCookie(&http.Cookie{Name: cookieName, Value: value})
								}
							}
						}
					}
				}
			}

			// Capture response cookies
			rw := &util.ResponseWriter{ResponseWriter: w}
			h.ServeHTTP(rw, r)

			// Rewrite response cookies
			for _, cookie := range w.Header()["Set-Cookie"] {
				parts := strings.SplitN(cookie, ";", 2)
				cookieParts := strings.SplitN(parts[0], "=", 2)
				if len(cookieParts) == 2 {
					name := cookieParts[0]
					value := cookieParts[1]
					targetURL, _ := url.Parse(r.URL.Query().Get("url"))
					domain := targetURL.Host
					cookieName := "pc_" + strings.ReplaceAll(domain, ".", "_") + "__" + name
					http.SetCookie(w, &http.Cookie{Name: cookieName, Value: value})
				}
			}
		})
	}
}