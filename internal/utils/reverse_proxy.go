package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// forwards full request to target url after trimming prefix
// keeps all parameters, headers, payload as is
func ReverseProxy(target, prefix string) gin.HandlerFunc {
	// convert target into URL object
	targetURL, err := url.Parse(target)
	if err != nil {
		panic(err)
	}

	// create reverse proxy object that forwards incoming requests to target url
	// automatically manage headers, req and res bodies
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Director rewrites req before it is sent to the backend
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req) // set basic URL and headers (default logic)

		// ensure req goes to correct backend host
		req.URL.Scheme = targetURL.Scheme // Scheme -> protocol (http/https/ws etc)
		req.URL.Host = targetURL.Host     // Host -> Domain name and port number

		// strip prefix from url path
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)

		// ensure path always starts with "/"
		if !strings.HasPrefix(req.URL.Path, "/") {
			req.URL.Path = "/" + req.URL.Path
		}

		// pass original query params, headers, etc
		req.Host = targetURL.Host

	}

	// convert standard proxy to a gin compatible handler
	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
