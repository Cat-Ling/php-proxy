package middleware

import (
	"go-proxy/util"
	"net/http"
	"net/url"
	"path"
)

func HeaderRewrite() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Request headers
			r.Header.Set("Accept-Encoding", "identity")
			r.Header.Del("Referer")

			rw := &util.ResponseWriter{ResponseWriter: w}
			h.ServeHTTP(rw, r)

			// Response headers
			if location := w.Header().Get("Location"); location != "" {
				targetURL, _ := url.Parse(r.URL.Query().Get("url"))
				w.Header().Set("Location", util.ProxifyURL(location, targetURL))
			}

			// Forward specific headers
			forwardHeaders := []string{"Content-Type", "Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition", "Location", "Set-Cookie"}
			for key := range w.Header() {
				isForwardHeader := false
				for _, forwardHeader := range forwardHeaders {
					if key == forwardHeader {
						isForwardHeader = true
						break
					}
				}
				if !isForwardHeader {
					w.Header().Del(key)
				}
			}

			// Add Content-Disposition if not present
			if w.Header().Get("Content-Disposition") == "" {
				filename := path.Base(r.URL.Path)
				w.Header().Set("Content-Disposition", "filename=\""+filename+"\"")
			}

			// Set cache-control headers
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		})
	}
}