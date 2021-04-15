package ipfinder

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type IPAddrs struct {
	V4 *net.IP
	V6 *net.IP
}

type ResolveFunc func() (*IPAddrs, error)

var Resolvers []ResolveFunc

func Ipify() (*IPAddrs, error) {
	res, err := http.Get("https://api64.ipify.org")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ip := net.ParseIP(string(body))
	if ip == nil {
		err = fmt.Errorf("failed to parse ip")
		return nil, err
	}

	if ip.To16() == nil {
		// if ipv4 return here because there are not IPv6
		return &IPAddrs{V4: &ip}, nil
	}

	res, err = http.Get("https://api.ipify.org")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ip2 := net.ParseIP(string(body))
	if ip2 == nil {
		err = fmt.Errorf("failed to parse ip")
		return nil, err
	}

	return &IPAddrs{V6: &ip, V4: &ip2}, nil
}

type Content struct {
	Result Result `json:"result"`
}

type Result struct {
	Data Data `json:"data"`
}

type Data struct {
	IPv4 net.IP `json:"IPAddress"`
	IPv6 net.IP `json:"IPv6Address"`
}

func Livebox() (*IPAddrs, error) {
	payload := `{"service": "NMC", "method": "getWANStatus", "parameters": {}}"`
	res, err := http.Post("http://192.168.1.1/ws", "application/x-sah-ws-1-call+json", strings.NewReader(payload))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	content := Content{}
	err = json.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}

	return &IPAddrs{V4: &content.Result.Data.IPv4, V6: &content.Result.Data.IPv6}, nil
}

func init() {
	Resolvers = append(Resolvers, Livebox, Ipify)
}
