package dyndns

import (
	"fmt"
	"log"
	"net"

	"github.com/mlcdf/dyndns/internal/discord"
	"github.com/mlcdf/dyndns/internal/gandi"
	"github.com/mlcdf/dyndns/internal/ipfinder"
	"github.com/pkg/errors"
)

type DynDNS struct {
	domain string
	record string
	ttl    int
}

func New(domain string, record string, ttl int) *DynDNS {
	return &DynDNS{domain: domain, record: record, ttl: ttl}
}

func (dyndns *DynDNS) Run(ipf ipfinder.IpFinderFunc, g *gandi.Client, d *discord.Client) error {
	resolvedIPs, err := ipf()
	if err != nil {
		return err
	}
	log.Printf("Current dynamic IP(s): %s\n", resolvedIPs)

	dnsRecords, err := g.Get(dyndns.domain, dyndns.record)
	if err != nil {
		return err
	}

	ipsToUpdate := matchIPs(resolvedIPs, dnsRecords)

	if len(ipsToUpdate) == 0 {
		log.Println("IP address match - no further action")
		return nil
	}

	for _, ip := range ipsToUpdate {
		err := g.Post(dyndns.domain, dyndns.record, ip, dyndns.ttl)
		if err != nil {
			return err
		}
	}

	log.Printf("DNS record for %s.%s updated\n", dyndns.record, dyndns.domain)

	err = dyndns.notifyDiscord(d, resolvedIPs.Values())
	return err
}

func (dyndns *DynDNS) notifyDiscord(d *discord.Client, ips []*net.IP) error {
	fields := make([]discord.Field, 0, len(ips))
	for _, ip := range ips {
		field := &discord.Field{Inline: true, Value: ip.String()}

		if ip.To4() != nil {
			field.Name = "v4"
		} else {
			field.Name = "v6"
		}

		fields = append(fields, *field)
	}

	err := d.PostSuccess(&discord.Webhook{
		Embeds: []discord.Embed{
			{
				Title:       fmt.Sprintf("DNS record for %s.%s updated with the new IP adresses", dyndns.record, dyndns.domain),
				Description: fmt.Sprintf("See [Gandi Live DNS](https://admin.gandi.net/domain/%s/records)", dyndns.domain),
				Fields:      fields,
			},
		},
	})
	return errors.Wrap(err, "failed to post success message to Discord")
}

func matchIPs(resolvedIPs *ipfinder.IPAddrs, dnsRecords []*gandi.DomainRecord) []*net.IP {
	toUpdate := make([]*net.IP, 0, 2)

	isIpV4ToUpdate := false
	isIpV6ToUpdate := false

	for _, records := range dnsRecords {
		for _, rrsetValue := range records.RrsetValues {
			if resolvedIPs.V4 != nil && rrsetValue.Equal(*resolvedIPs.V4) {
				isIpV4ToUpdate = true
			}

			if resolvedIPs.V6 != nil && rrsetValue.Equal(*resolvedIPs.V6) {
				isIpV6ToUpdate = true
			}
		}
	}

	if resolvedIPs.V4 != nil && !isIpV4ToUpdate {
		toUpdate = append(toUpdate, resolvedIPs.V4)
	}

	if resolvedIPs.V6 != nil && !isIpV6ToUpdate {
		toUpdate = append(toUpdate, resolvedIPs.V6)
	}

	return toUpdate
}
