package middleware

import (
	"net/http"
	"strings"
)

func routePattern(r *http.Request) string {
	if r.Pattern == "" {
		return r.URL.Path
	}
	if _, path, found := strings.Cut(r.Pattern, " "); found {
		return path
	}
	return r.Pattern
}
