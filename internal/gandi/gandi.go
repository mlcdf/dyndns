package gandi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"go.mlcdf.fr/dyndns/internal/httpx"
)

type Client struct {
	httpClient *http.Client
	Token      string
}

// New creates a new Client with a default timeout
func New(token string) *Client {
	client := httpx.DefaultClient
	client.Timeout = time.Second * 10

	return &Client{client, token}
}

// DomainRecord represents a DNS Record
type DomainRecord struct {
	RrsetType   string    `json:"rrset_type,omitempty"`
	RrsetTTL    int       `json:"rrset_ttl,omitempty"`
	RrsetName   string    `json:"rrset_name,omitempty"`
	RrsetHref   string    `json:"rrset_href,omitempty"`
	RrsetValues []*net.IP `json:"rrset_values,omitempty"`
}

func (c *Client) Get(domain string, record string) ([]*DomainRecord, error) {
	url := fmt.Sprintf("https://api.gandi.net/v5/livedns/domains/%s/records/%s", domain, record)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "ApiKey "+c.Token)
	req.Header.Set("Content-type", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	records := make([]*DomainRecord, 0)
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

func (c *Client) Put(domain string, name string, ips []*net.IP, ttl int) error {
	record := struct {
		Items []*DomainRecord `json:"items"`
	}{Items: make([]*DomainRecord, 0, 2)}

	for _, ip := range ips {
		item := &DomainRecord{RrsetTTL: ttl, RrsetValues: []*net.IP{ip}, RrsetType: rrsetType(ip)}
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

	res, err := c.httpClient.Do(req)

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

func (c *Client) Delete(domain string, name string, rtype string) error {
	url := fmt.Sprintf("https://api.gandi.net/v5/livedns/domains/%s/records/%s/%s", domain, name, rtype)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "ApiKey "+c.Token)
	req.Header.Set("Content-type", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("failed to perform DELETE status=%d response=%s", res.StatusCode, body)
	}

	return nil
}
