package gopress

import (
	"crypto/tls"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestRawURL(t *testing.T) {
	cases := []struct {
		host, uri, scheme string
		tls               *tls.ConnectionState
		expect            string
	}{
		{"google.com", "/ncr", "", nil, "http://google.com/ncr"},
		{"stackoverflow.com", "/q?id=1", "https", nil, "https://stackoverflow.com/q?id=1"},
		{"github.com", "/fpay", "", &tls.ConnectionState{}, "https://github.com/fpay"},
	}

	for _, c := range cases {
		req, _ := http.NewRequest("GET", "/", nil)
		req.Host = c.host
		req.RequestURI = c.uri
		req.TLS = c.tls
		req.Header.Set(RequestHeaderProtocol, c.scheme)

		assert.Equal(t, c.expect, RequestRawURL(req), "raw url should match")
	}
}

func TestRequestRemoteAddr(t *testing.T) {
	cases := []struct {
		uri, addr string
		forwarded string
		expect    string
	}{
		{"/ncr", "127.0.0.1", "", "127.0.0.1"},
		{"/q?id=1", "127.0.0.1", "", "127.0.0.1"},
		{"/fpay", "127.0.0.1", "10.0.0.1", "10.0.0.1"},
	}

	for _, c := range cases {
		req, _ := http.NewRequest("GET", "/", nil)
		req.RemoteAddr = c.addr
		req.Header.Set(RequestHeaderForwarded, c.forwarded)

		assert.Equal(t, c.expect, RequestRemoteAddr(req), "remote addr should match")
	}
}
