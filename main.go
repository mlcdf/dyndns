package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	dyndns "go.mlcdf.fr/dyndns/internal"
	"go.mlcdf.fr/dyndns/internal/discord"
	"go.mlcdf.fr/dyndns/internal/gandi"
	"go.mlcdf.fr/sally/build"
)

const usage = `Usage:
    dyndns --domain [DOMAIN] --record [RECORD]

Options:
    --ttl                Time to live in seconds. Defaults to 3600
    --always-notify      Always notify the Discord channel (even when nothing changes)
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
		versionFlag      bool
		domainFlag       string
		recordFlag       string
		ttlFlag          int = 3600
		alwaysNotifyFlag bool
	)

	flag.StringVar(&domainFlag, "domain", domainFlag, "")
	flag.StringVar(&recordFlag, "record", recordFlag, "")

	flag.IntVar(&ttlFlag, "ttl", ttlFlag, "Time to live. Defaults to 3600.")

	flag.BoolVar(&versionFlag, "version", versionFlag, "print the version")
	flag.BoolVar(&versionFlag, "V", versionFlag, "print the version")

	flag.BoolVar(&alwaysNotifyFlag, "always-notify", alwaysNotifyFlag, "")

	flag.Parse()

	if versionFlag {
		fmt.Println("dyndns " + build.String())
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

	gandiClient := gandi.New(mustEnv("GANDI_TOKEN"))
	dyn := dyndns.New(gandiClient, discordClient, alwaysNotifyFlag)

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
