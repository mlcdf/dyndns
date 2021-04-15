package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"

	"github.com/mlcdf/dyndns/internal/discord"
	"github.com/mlcdf/dyndns/internal/gandi"
	"github.com/mlcdf/dyndns/internal/ipfinder"
)

const usage = `Usage:
    dyndns --domain [DOMAIN] --record [RECORD]

Options:
    --livebox            Query the Livebox (router) to find the IP instead of api.ipify.org
    --ttl                Time to live. Defaults to 3600
    -V, --version        Print version

Examples:
    export DISCORD_WEBHOOK_URL='https://discord.com/api/webhooks/xxx'
    export GANDI_TOKEN='foobar'
    dyndns --domain example.com --record "*.pi"

How to create a Discord webhook: https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
How to generate your Gandi token: https://docs.gandi.net/en/domain_names/advanced_users/api.html
`

// Version can be set at link time to override debug.BuildInfo.Main.Version,
// which is "(devel)" when building from within the module. See
// golang.org/issue/29814 and golang.org/issue/29228.
var Version string

func main() {
	log.SetFlags(0)
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	var (
		versionFlag bool
		domainFlag  string
		recordFlag  string
		ttlFlag     int = 3600
		liveboxFlag bool
	)

	flag.StringVar(&domainFlag, "domain", domainFlag, "")
	flag.StringVar(&recordFlag, "record", recordFlag, "")

	flag.BoolVar(&liveboxFlag, "livebox", liveboxFlag, "Use the Livebox IP resolver instead of api.ipify.org")
	flag.IntVar(&ttlFlag, "ttl", ttlFlag, "Time to live. Defaults to 3600.")

	flag.BoolVar(&versionFlag, "version", versionFlag, "print the version")
	flag.BoolVar(&versionFlag, "V", versionFlag, "print the version")

	flag.Parse()

	if versionFlag {
		if Version != "" {
			fmt.Println(Version)
			return
		}
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			fmt.Println(buildInfo.Main.Version)
			return
		}
		fmt.Println("(unknown)")
		return
	}

	disc := &discord.Client{WebhookURL: mustEnv("DISCORD_WEBHOOK_URL")}
	log.SetOutput(io.MultiWriter(os.Stderr, disc))

	if domainFlag == "" {
		log.Fatal("error: required flag --domain is missing")
	}

	if recordFlag == "" {
		log.Fatal("error: required flag --record is missing")
	}

	resolvedIPs := &ipfinder.IPAddrs{}
	var err error

	if liveboxFlag {
		resolvedIPs, err = ipfinder.Livebox()
	} else {
		resolvedIPs, err = ipfinder.Ipify()
	}

	if err != nil {
		log.Fatal(err)
	}

	gandi := &gandi.Client{Token: mustEnv("GANDI_TOKEN")}

	dnsRecords, err := gandi.Get(domainFlag, recordFlag)
	if err != nil {
		log.Fatal(err)
	}

	var isIPv4UpToDate bool
	var isIPv6UpToDate bool

	for _, records := range dnsRecords {
		for _, rrsetValue := range records.RrsetValues {
			if rrsetValue.Equal(*resolvedIPs.V4) {
				isIPv4UpToDate = true
			}

			if rrsetValue.Equal(*resolvedIPs.V6) {
				isIPv6UpToDate = true
			}
		}
	}

	if isIPv4UpToDate && isIPv6UpToDate {
		fmt.Fprintln(os.Stderr, "success: nothing to do")
		os.Exit(0)
	}

	if !isIPv4UpToDate {
		err := gandi.Post(domainFlag, recordFlag, resolvedIPs.V4, ttlFlag)
		if err != nil {
			log.Fatal(err)
		}
	}

	if !isIPv6UpToDate {
		err := gandi.Post(domainFlag, recordFlag, resolvedIPs.V6, ttlFlag)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Fprintf(os.Stderr, "success: DNS record for %s.%s updated with the new IP adresses: %s, %s", recordFlag, domainFlag, resolvedIPs.V4, resolvedIPs.V6)
	err = disc.PostSuccess(&discord.Webhook{
		Embeds: []discord.Embed{
			{
				Title:       fmt.Sprintf("DNS record for %s.%s updated with the new IP adresses", recordFlag, domainFlag),
				Description: fmt.Sprintf("See [Gandi Live DNS](https://admin.gandi.net/domain/%s/records)", domainFlag),
				Fields: []discord.Field{
					{
						Name:   "v4",
						Value:  resolvedIPs.V4.String(),
						Inline: true,
					},
					{
						Name:   "v6",
						Value:  resolvedIPs.V6.String(),
						Inline: true,
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("error: required environment variable %s is missing", key)
	}
	return value
}
