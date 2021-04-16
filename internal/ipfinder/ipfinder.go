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

type IpFinderFunc func() (*IPAddrs, error)

type IPAddrs struct {
	V4 *net.IP `json:"IPAddress"`
	V6 *net.IP `json:"IPv6Address"`
}

func (ipAddrs *IPAddrs) String() string {
	str := ""
	if ipAddrs.V4 != nil {
		str += ipAddrs.V4.String()
	}
	if ipAddrs.V6 != nil {
		str += ", " + ipAddrs.V6.String()
	}
	return str
}

func (ipAddrs *IPAddrs) Values() []*net.IP {
	values := make([]*net.IP, 0, 2)
	if ipAddrs.V4 != nil {
		values = append(values, ipAddrs.V4)
	}
	if ipAddrs.V6 != nil {
		values = append(values, ipAddrs.V6)
	}
	return values
}

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
		return nil, fmt.Errorf("failed to parse ip")
	}

	if ip.To4() != nil {
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
		return nil, fmt.Errorf("failed to parse ip")
	}

	return &IPAddrs{V6: &ip, V4: &ip2}, nil
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

	type Content struct {
		Result struct {
			Data IPAddrs `json:"data"`
		} `json:"result"`
	}

	content := Content{}
	err = json.Unmarshal(data, &content)
	if err != nil {
		return nil, err
	}

	return &content.Result.Data, nil
}
