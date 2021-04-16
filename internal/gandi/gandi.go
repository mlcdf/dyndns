package gandi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	Token string
}

// DomainRecord represents a DNS Record
type DomainRecord struct {
	RrsetType   string    `json:"rrset_type,omitempty"`
	RrsetTTL    int       `json:"rrset_ttl,omitempty"`
	RrsetName   string    `json:"rrset_name,omitempty"`
	RrsetHref   string    `json:"rrset_href,omitempty"`
	RrsetValues []*net.IP `json:"rrset_values,omitempty"`
}

func (c *Client) Get(fqdn string, name string) ([]*DomainRecord, error) {
	url := fmt.Sprintf("https://api.gandi.net/v5/livedns/domains/%s/records/%s", fqdn, name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "ApiKey "+c.Token)

	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	records := make([]*DomainRecord, 0)
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get records response=%s", body)
	}

	return records, nil
}

func rrsetType(ip *net.IP) string {
	if xx := ip.To4(); xx == nil {
		return "AAAA"
	}
	return "A"
}

func (c *Client) Post(fqdn string, name string, ip *net.IP, ttl int) error {
	record := &DomainRecord{RrsetTTL: ttl, RrsetValues: []*net.IP{ip}, RrsetType: rrsetType(ip)}

	payload, err := json.Marshal(record)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.gandi.net/v5/livedns/domains/%s/records/%s", fqdn, name)

	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "ApiKey "+c.Token)
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("failed to perform update response=%s", body)
	}
	return nil
}
