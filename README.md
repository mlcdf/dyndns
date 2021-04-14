# dyndns

Keeps Gandi LiveDNS records up to update with the current (dynamic) IP and reports failures and successful updates on a Discord channel.

## Install

- From [GitHub releases](https://github.com/mlcdf/dyndns/releases): download the binary corresponding to your OS and architecture.
- From source (make sure `$GOPATH/bin` is in your `$PATH`):
```sh
go get https://github.com/mlcdf/dyndns
```

## Usage

```
Usage:
    dyndns --domain [DOMAIN] --record [RECORD]

Options:
    --ttl                Time to live. Defaults to 3600.
    -v, --verbose        Print verbose output
    -V, --version        Print version

Examples:
    export DISCORD_WEBHOOK='https://discord.com/api/webhooks/xxx'
    export GANDI_TOKEN='foobar'
    dyndns --domain example.com --record "*.pi"
```

Setup as a `cron` jon
```bash
crontab -e
# Run every 10 minutes
*/10 * * * * /home/pi/dyndns --domain example.com --record "*.pi"
```

## Development

Run the app
```sh
go run main.go
```

Run the tests
```sh
go test ./...
```
