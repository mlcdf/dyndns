package smockertest

import (
	"fmt"
	"net/http"
	"os"
)

// RedirectTransport implement a Roundtrip that redirects requests to a running smocker instance
type RedirectTransport struct{}

var _ http.RoundTripper = &RedirectTransport{}

func (s *RedirectTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "http"
	r.URL.Host = fmt.Sprintf("localhost:%d", port)
	return http.DefaultTransport.RoundTrip(r)
}

// PushMock send the mockfile to the smocker server
func PushMock(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:%d/mocks?reset=true", adminPort)
	res, err := http.Post(url, "content-type: application/x-yaml", f)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("error %d while performing POST %s", res.StatusCode, url)
	}

	return err
}
