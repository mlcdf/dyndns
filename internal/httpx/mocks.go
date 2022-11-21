//go:build smocker

package httpx

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

var DefaultClient = &http.Client{Transport: &TestTransport{}}

type TestTransport struct{}

var smockerURL *url.URL

func init() {
	endpoint := os.Getenv("SMOCKER_ENDPOINT")
	if endpoint == "" {
		log.Fatalf("missing required SMOCKER_ENDPOINT environment variable")
	}

	var err error
	smockerURL, err = url.Parse(endpoint)
	if err != nil {
		log.Fatalf("error parsing %s as an url.URL", endpoint)
	}
}

func (s *TestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Host = smockerURL.Host
	r.URL.Scheme = smockerURL.Scheme
	return http.DefaultTransport.RoundTrip(r)
}
