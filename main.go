package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"go.mlcdf.fr/dyndns/internal/build"
	"go.mlcdf.fr/dyndns/internal/discord"
	"go.mlcdf.fr/dyndns/internal/dyndns"
	"go.mlcdf.fr/dyndns/internal/gandi"
	"go.mlcdf.fr/dyndns/internal/ipfinder"
)

const usage = `Usage:
    dyndns --domain [DOMAIN] --record [RECORD]

Options:
    --livebox            Query the Livebox (router) to find the IP instead of api.ipify.org
    --ttl                Time to live in seconds. Defaults to 3600
    -V, --version        Print version

Examples:
    export DISCORD_WEBHOOK_URL='https://discord.com/api/webhooks/xxx'
    export GANDI_TOKEN='foobar'
    dyndns --domain example.com --record "*.pi"

How to create a Discord webhook: https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
How to generate your Gandi token: https://docs.gandi.net/en/domain_names/advanced_users/api.html
`

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
		fmt.Println(build.Long())
		return
	}

	discordClient := &discord.Client{WebhookURL: mustEnv("DISCORD_WEBHOOK_URL")}
	logErr := log.New(io.MultiWriter(os.Stderr, discordClient), "", 0)

	if domainFlag == "" {
		logErr.Fatal("error: required flag --domain is missing")
	}

	if recordFlag == "" {
		logErr.Fatal("error: required flag --record is missing")
	}

	var ipf ipfinder.IpFinderFunc
	if liveboxFlag {
		ipf = ipfinder.Livebox
	} else {
		ipf = ipfinder.Ipify
	}

	gandiClient := gandi.New(mustEnv("GANDI_TOKEN"))
	dyn := dyndns.New(ipf, gandiClient, discordClient)

	err := dyn.Run(domainFlag, recordFlag, ttlFlag)
	if err != nil {
		logErr.Fatalf("error: %v", err)
	}
}

func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("error: required environment variable %s is empty or missing", key)
	}
	return value
}
