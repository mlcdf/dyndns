#!/bin/bash
set -eu -o pipefail

docker compose up -d

go build -tags smocker -o dist/dyndns.test .

export SMOCKER_ENDPOINT="http://localhost:8080"
export VENOM_VERBOSE=1 \
export VENOM_VAR_owh=$(pwd)/owh-test \
export VENOM_VAR_owh_consumer_key=${OWH_CONSUMER_KEY}

go run github.com/ovh/venom/cmd/venom@v1.0.1 run $(find tests/venom ! -path '*vars*' -type f -name '*.yml' | sort)