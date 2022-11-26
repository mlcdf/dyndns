# dyndns

[![test](https://github.com/mlcdf/dyndns/actions/workflows/test.yml/badge.svg)](https://github.com/mlcdf/dyndns/actions/workflows/test.yml)
[![coverage](https://raw.githubusercontent.com/mlcdf/dyndns/coverage/badge.svg)](https://mlcdf.github.io/dyndns/report.html)

Update Gandi LiveDNS based on the current (dynamic) ip.

## Why

Some (most ?) providers ([Orange](https://orange.fr/) for example) don't provide you with a fix IP address: they give you a dynamic IP that changes over time. This makes it hard to self-host software and share them publicly via a domain because you have to constantly
change the IP your domain is pointing to.

This program aims to solve that. It's intended to run on a always-on computer in your home (such as Raspberry Pi).

## Highlights

- Supports both IPv4 and IPv6.
- Ship as single binary with no OS dependencies.
- Reports failures and successful updates on a Discord channel. /!\ This is not optional (by design).

## Install

- From [GitHub releases](https://go.mlcdf.fr/dyndns/releases): download the binary corresponding to your OS and architecture.
- From source (make sure `$GOPATH/bin` is in your `$PATH`):

```sh
go install go.mlcdf.fr/dyndns
```

## Setup

`dyndns` requires the following environment variables to be set:

| name                  | description                                                                                               |
|-----------------------|-----------------------------------------------------------------------------------------------------------|
| `DISCORD_WEBHOOK_URL` | your [Discord channel webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks) |
| `GANDI_TOKEN`         | your [Gandi API Key](https://docs.gandi.net/en/domain_names/advanced_users/api.html)                      |

## Usage

```
Usage:
    dyndns --domain [DOMAIN] --record [RECORD]

Options:
    --ttl                Time to live in seconds. Defaults to 3600
    --always-notify      Always notify the Discord channel (even when nothing changes)
    -V, --version        Print version

Examples:
    export DISCORD_WEBHOOK_URL='https://discord.com/api/webhooks/xxx'
    export GANDI_TOKEN='foobar'
    dyndns --domain example.com --record "*.pi"
```

Setup as a `cron` job

```bash
crontab -e
# Run every 10 minutes
*/10 * * * * /path/to/dyndns --domain example.com --record "*.pi"
```

Check out the example [contrib/deploy.sh](./contrib/deploy.sh) script.

## Development

Requirements:
- go version > 1.19+
- docker-compose version > 2, installed as docker plugins (invoked via `docker compose`)

Run the app

```sh
go run main.go
```

Run the tests

```sh
go test ./...
```

Force `go test` to run all the tests (and don't kill the docker-compose containers so the following runs will be faster).
```sh
./scrits/test.sh
```

## License

[MIT](https://choosealicense.com/licenses/mit/)