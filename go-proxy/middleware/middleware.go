package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}
	return m[0](Chain(h, m[1:cap(m)]...))
}