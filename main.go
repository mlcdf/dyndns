package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

var (
	// Injected from linker flags like `go build -ldflags "-X main.version=$VERSION" -X ...`
	isTest = "false"
)

func main() {
	code := int(mainRun())
	os.Exit(code)
}

func mainRun() exitCode {
	log.SetFlags(0)
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	if len(os.Args) == 1 {
		flag.Usage()
		return exitOK
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
		fmt.Fprintln(os.Stdout, "dyndns "+build.String())
		return exitOK
	}

	webhook := os.Getenv("DISCORD_WEBHOOK_URL")
	if webhook == "" {
		log.Println("error: required environment variable DISCORD_WEBHOOK_URL is empty or missing")
		return exitError
	}

	discordClient := &discordClient{webhook}
	logErr := log.New(io.MultiWriter(os.Stderr, discordClient), "", 0)

	if domainFlag == "" {
		logErr.Println("error: required flag --domain is missing")
		return exitError
	}

	if recordFlag == "" {
		logErr.Println("error: required flag --record is missing")
		return exitError
	}

	token := os.Getenv("GANDI_TOKEN")
	if token == "" {
		log.Println("error: required environment variable GANDI_TOKEN is empty or missing")
		return exitError
	}

	gandiClient := &gandiClient{token}

	dyn := &DynDNS{
		gandiClient,
		discordClient,
	}

	err := dyn.execute(domainFlag, recordFlag, ttlFlag, alwaysNotifyFlag)
	if err != nil {
		logErr.Printf("error: %v", err)
		return exitError
	}

	return exitOK
}

// func exit(code int) {
// 	if isTest == "true" {
// 		bincover.ExitCode = code
// 		os.Exit(0)
// 	} else {
// 		os.Exit(code)
// 	}
// }
