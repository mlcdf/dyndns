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

// DynDNS holds all the required dependencies
type DynDNS struct {
	finderFunc    ipfinder.IpFinderFunc
	gandiClient   *gandi.Client
	discordClient *discord.Client
}

// New creates a new DynDNS
func New(ipf ipfinder.IpFinderFunc, g *gandi.Client, d *discord.Client) *DynDNS {
	return &DynDNS{finderFunc: ipf, gandiClient: g, discordClient: d}
}

// Run check the current IPs, and the one defines in the DNS records.
// If necessary, it updates the DNS records and notify Discord.
func (dyndns *DynDNS) Run(domain string, record string, ttl int) error {
	resolvedIPs, err := dyndns.finderFunc()
	if err != nil {
		return err
	}
	log.Printf("Current dynamic IP(s): %s\n", resolvedIPs)

	dnsRecords, err := dyndns.gandiClient.Get(domain, record)
	if err != nil {
		return err
	}

	ipsToUpdate := matchIPs(resolvedIPs, dnsRecords)

	if len(ipsToUpdate) == 0 {
		log.Println("IP address match - no further action")
		return nil
	}

	for _, ip := range ipsToUpdate {
		err := dyndns.gandiClient.Post(domain, record, ip, ttl)
		if err != nil {
			return err
		}
	}

	log.Printf("DNS record for %s.%s updated\n", record, domain)

	err = dyndns.notifyDiscord(domain, record, resolvedIPs.Values())
	return err
}

func (dyndns *DynDNS) notifyDiscord(domain string, record string, ips []*net.IP) error {
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

	err := dyndns.discordClient.PostSuccess(&discord.Webhook{
		Embeds: []discord.Embed{
			{
				Title:       fmt.Sprintf("DNS record for %s.%s updated with the new IP adresses", record, domain),
				Description: fmt.Sprintf("See [Gandi Live DNS](https://admin.gandi.net/domain/%s/records)", domain),
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