package main

import (
	"net/http"
	"time"

	"go.mlcdf.fr/dyndns/tests/smockertest"
)

var defaultHTTP = &http.Client{Timeout: 20 * time.Second}

func init() {
	if isTest == "true" {
		defaultHTTP.Transport = &smockertest.RedirectTransport{}
	}
}
