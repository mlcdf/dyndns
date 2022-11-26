package main

import (
	"net/http"
	"time"
)

var defaultHTTP = &http.Client{Timeout: 20 * time.Second}

type TestTransport struct{}

func (s *TestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "http"
	r.URL.Host = "localhost:8080"
	return http.DefaultTransport.RoundTrip(r)
}

func init() {
	if isTest == "true" {
		defaultHTTP.Transport = &TestTransport{}
	}
}
