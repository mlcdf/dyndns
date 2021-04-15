# dyndns

Update Gandi LiveDNS based on the current (dynamic) ip.

## Highlights

- Supports both IPv4 and IPv6.
- Reports failures and successful updates on a Discord channel. /!\ This is not optional (by design).
- A faster IP finder is available if you have a Livebox v4 (may works on other models).

## Install

- From [GitHub releases](https://github.com/mlcdf/dyndns/releases): download the binary corresponding to your OS and architecture.
- From source (make sure `$GOPATH/bin` is in your `$PATH`):

```sh
go get https://github.com/mlcdf/dyndns
```

## Config

`dyndns` requires the following environment variables to be set:

| name                  | description                  | docs                                                                      |
| --------------------- | ---------------------------- | ------------------------------------------------------------------------- |
| `DISCORD_WEBHOOK_URL` | your Discord channel webhook | https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks |
| `GANDI_TOKEN`         | your Gandi API Key           | https://docs.gandi.net/en/domain_names/advanced_users/api.html            |

## Usage

```
Usage:
    dyndns --domain [DOMAIN] --record [RECORD]

Options:
    --livebox            Use the Livebox IP resolver instead of api.ipify.org
    --ttl                Time to live. Defaults to 3600
    -V, --version        Print version
```

Setup as a `cron` jon

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
