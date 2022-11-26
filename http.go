package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var defaultHTTP = &http.Client{Timeout: 20 * time.Second}

type TestTransport struct{}

var smockerURL *url.URL

func init() {
	if isTest == "false" {
		return
	}

	endpoint := os.Getenv("SMOCKER_ENDPOINT")
	if endpoint == "" {
		log.Fatalf("missing required SMOCKER_ENDPOINT environment variable")
	}

	var err error
	smockerURL, err = url.Parse(endpoint)
	if err != nil {
		log.Fatalf("error parsing %s as an url.URL", endpoint)
	}

	defaultHTTP.Transport = &TestTransport{}
}

func (s *TestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Host = smockerURL.Host
	r.URL.Scheme = smockerURL.Scheme
	return http.DefaultTransport.RoundTrip(r)
}
