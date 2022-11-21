package dyndns

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/pkg/errors"
	"go.mlcdf.fr/dyndns/internal/discord"
	"go.mlcdf.fr/dyndns/internal/gandi"
	"go.mlcdf.fr/dyndns/internal/httpx"
)

// DynDNS holds all the required dependencies
type DynDNS struct {
	gandiClient   *gandi.Client
	discordClient *discord.Client
	alwaysNotify  bool
}

type IPAddrs struct {
	V4 *net.IP `json:"IPAddress"`
	V6 *net.IP `json:"IPv6Address"`
}

// New creates a new DynDNS
func New(g *gandi.Client, d *discord.Client, alwaysNotify bool) *DynDNS {
	return &DynDNS{gandiClient: g, discordClient: d, alwaysNotify: alwaysNotify}
}

func (ipAddrs *IPAddrs) String() string {
	str := "["
	if ipAddrs.V4 != nil {
		str += ipAddrs.V4.String()
	}
	if ipAddrs.V6 != nil {
		str += " " + ipAddrs.V6.String()
	}
	str += "]"
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

// resolveIPs finds the current IP(s) addresses pointu
func resolveIPs() (*IPAddrs, error) {
	res, err := httpx.DefaultClient.Get("https://api64.ipify.org")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
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

	res, err = httpx.DefaultClient.Get("https://api.ipify.org")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	ip2 := net.ParseIP(string(body))
	if ip2 == nil {
		return nil, fmt.Errorf("failed to parse ip")
	}

	return &IPAddrs{V6: &ip, V4: &ip2}, nil
}

// Run check the current IPs, and the one defines in the DNS records.
// If necessary, it updates the DNS records and notify Discord.
func (dyndns *DynDNS) Run(domain string, record string, ttl int) error {
	resolvedIPs, err := resolveIPs()
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
		log.Println("IP address(es) match - no further action")

		if dyndns.alwaysNotify {
			err := dyndns.discordClient.PostInfo(&discord.Webhook{
				Embeds: []discord.Embed{
					{
						Title:       fmt.Sprintf("IP address(es) match for record %s.%s - no further action", record, domain),
						Description: "To disable notifications when nothing happens, remove the `--always-notify` flag",
					},
				},
			})
			return errors.Wrap(err, "failed to send message to discord")
		}
		return nil
	}

	err = dyndns.gandiClient.Put(domain, record, ipsToUpdate, ttl)
	if err != nil {
		return err
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

func matchIPs(resolvedIPs *IPAddrs, dnsRecords []*gandi.DomainRecord) []*net.IP {
	toUpdate := make([]*net.IP, 0, 2)

	isIpV4ToUpdate := false
	isIpV6ToUpdate := false

	ipsFromDNS := make([]*net.IP, 0, 2)

	for _, records := range dnsRecords {
		for _, rrsetValue := range records.RrsetValues {
			ipsFromDNS = append(ipsFromDNS, rrsetValue)

			if resolvedIPs.V4 != nil && rrsetValue.Equal(*resolvedIPs.V4) {
				isIpV4ToUpdate = true
			}

			if resolvedIPs.V6 != nil && rrsetValue.Equal(*resolvedIPs.V6) {
				isIpV6ToUpdate = true
			}
		}
	}

	log.Printf("IP(s) from DNS:        %s", ipsFromDNS)

	if resolvedIPs.V4 != nil && !isIpV4ToUpdate {
		toUpdate = append(toUpdate, resolvedIPs.V4)
	}

	if resolvedIPs.V6 != nil && !isIpV6ToUpdate {
		toUpdate = append(toUpdate, resolvedIPs.V6)
	}

	return toUpdate
}
