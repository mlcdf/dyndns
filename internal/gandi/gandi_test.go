package gandi

import (
	"log"
	"net"
	"os"
	"testing"
)

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("error: required environment variable %s is empty or missing", key)
	}
	return value
}

func TestGandiPut(t *testing.T) {
	token := mustEnv("GANDI_TOKEN")

	client := New(token)
	ip, _, err := net.ParseCIDR("192.168.1.1/24")
	if err != nil {
		t.Fatal(err)
	}

	err = client.Put("mlcdf.fr", "golang-test", []*net.IP{&ip}, 1200)
	if err != nil {
		t.Fatal(err)
	}

	err = client.Delete("mlcdf.fr", "golang-test", "A")
	if err != nil {
		t.Fatal(err)
	}
}
