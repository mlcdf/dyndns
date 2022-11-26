package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

type gandiClient struct {
	Token string
}

// domainRecord represents a DNS Record
type domainRecord struct {
	RrsetType   string    `json:"rrset_type,omitempty"`
	RrsetTTL    int       `json:"rrset_ttl,omitempty"`
	RrsetName   string    `json:"rrset_name,omitempty"`
	RrsetHref   string    `json:"rrset_href,omitempty"`
	RrsetValues []*net.IP `json:"rrset_values,omitempty"`
}

func (c *gandiClient) get(domain string, record string) ([]*domainRecord, error) {
	url := fmt.Sprintf("https://api.gandi.net/v5/livedns/domains/%s/records/%s", domain, record)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "ApiKey "+c.Token)
	req.Header.Set("Content-type", "application/json")

	res, err := defaultHTTP.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	records := make([]*domainRecord, 0)
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get %s/records/%s  response=%s", domain, record, body)
	}

	return records, nil
}

func rrsetType(ip *net.IP) string {
	if xx := ip.To4(); xx == nil {
		return "AAAA"
	}
	return "A"
}

func (c *gandiClient) put(domain string, name string, ips []*net.IP, ttl int) error {
	record := struct {
		Items []*domainRecord `json:"items"`
	}{Items: make([]*domainRecord, 0, 2)}

	for _, ip := range ips {
		item := &domainRecord{RrsetTTL: ttl, RrsetValues: []*net.IP{ip}, RrsetType: rrsetType(ip)}
		record.Items = append(record.Items, item)
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.gandi.net/v5/livedns/domains/%s/records/%s", domain, name)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "ApiKey "+c.Token)
	req.Header.Set("Content-type", "application/json")

	res, err := defaultHTTP.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("failed to perform PUT status=%d response=%s", res.StatusCode, body)
	}

	return nil
}
