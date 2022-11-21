//go:build !smocker

package httpx

import "net/http"

var DefaultClient = http.DefaultClient
