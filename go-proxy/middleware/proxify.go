package middleware

import (
	"net/http"
)

import (
	"bytes"
	"go-proxy/util"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func Proxify() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body := &bytes.Buffer{}
			rw := &util.ResponseWriter{ResponseWriter: w, Body: body}
			h.ServeHTTP(rw, r)

			contentType := w.Header().Get("Content-Type")
			if !strings.HasPrefix(contentType, "text/html") {
				w.Write(body.Bytes())
				return
			}

			doc, err := html.Parse(body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var f func(*html.Node)
			f = func(n *html.Node) {
				if n.Type == html.ElementNode {
					for i, a := range n.Attr {
						if a.Key == "href" || a.Key == "src" || a.Key == "action" {
							targetURL, _ := url.Parse(r.URL.Query().Get("url"))
							n.Attr[i].Val = util.ProxifyURL(a.Val, targetURL)
						}
					}
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
			f(doc)

			buf := &bytes.Buffer{}
			if err := html.Render(buf, doc); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(buf.Bytes())
		})
	}
}