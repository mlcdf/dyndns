# dyndns

Update Gandi LiveDNS based on the current (dynamic) ip.

## Highlights

- Supports both IPv4 and IPv6.
- Reports failures and successful updates on a Discord channel. /!\ This is not optional (by design).
- A faster IP finder is available if you have a Livebox v4 (may works on other models).

## Install

- From [GitHub releases](https://go.mlcdf.fr/dyndns/releases): download the binary corresponding to your OS and architecture.
- From source (make sure `$GOPATH/bin` is in your `$PATH`):

```sh
go get go.mlcdf.fr/dyndns
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
    --livebox            Query the Livebox (router) to find the IP instead of api.ipify.org
    --ttl                Time to live. Defaults to 3600
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

Checkout the example [contrib/deploy.sh](./contrib/deploy.sh) script.

## Development

Run the app

```sh
go run main.go
```

Run the tests

```sh
go test ./...
```
