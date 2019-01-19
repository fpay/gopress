package gopress

import (
	"bytes"
	"net/http"
)

var (
	// RequestHeaderProtocol header for identifying the protocol (HTTP or HTTPS) that a client used to connect to your proxy or load balancer.
	RequestHeaderProtocol = "X-Forwarded-Proto"

	// RequestHeaderForwarded header for identifying the originating IP address of a client connecting to a web server through an HTTP proxy or a load balancer.
	RequestHeaderForwarded = "X-Forwarded-For"

	// RequestHeaderRequestedWith header for JavaScript libraries sending requests from browser.
	RequestHeaderRequestedWith = "X-Requested-With"
)

const (
	urlSchemeHTTP  = "http"
	urlSchemeHTTPS = "https"
)

// RequestRawURL returns request original URL.
func RequestRawURL(r *http.Request) string {
	scheme := RequestScheme(r)
	host := r.Host
	path := r.RequestURI

	buf := new(bytes.Buffer)
	buf.WriteString(scheme)
	buf.WriteString("://")
	buf.WriteString(host)
	buf.WriteString(path)
	return buf.String()
}

// RequestScheme try to parses request scheme.
//
// If the web server is behind an HTTP proxy or a load balancer, it's hard to known if original request's scheme by
// http.Request.TLS property. Most HTTP proxies and load balancers will attach a header to tell upstream server
// the request's real scheme. So check if the header is set first.
func RequestScheme(r *http.Request) string {
	scheme := r.Header.Get(RequestHeaderProtocol)
	if scheme == "" {
		if r.TLS != nil {
			scheme = urlSchemeHTTPS
		} else {
			scheme = urlSchemeHTTP
		}
	}
	return scheme
}

// RequestRemoteAddr finds the real remote address from request.
//
// If the web server is behind an HTTP proxy or a load balancer, http.Request.RemoteAddr is IP of the proxy or load
// balancer. But most HTTP proxies and load balancers will attach a header to tell the web server request's real IP.
func RequestRemoteAddr(req *http.Request) string {
	s := req.Header.Get(RequestHeaderForwarded)
	if s == "" {
		return req.RemoteAddr
	}
	return s
}

// RequestIsAJAX check request header to see if the request is a XMLHttpRequest.
func RequestIsAJAX(req *http.Request) bool {
	return "XMLHttpRequest" == req.Header.Get(RequestHeaderRequestedWith)
}
