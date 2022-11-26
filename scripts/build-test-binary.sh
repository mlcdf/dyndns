#!/usr/bin/env bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd .. && go test . -tags testbincover -coverpkg="./..." -c -o "${SCRIPT_DIR}/../dist/dyndns.test" -ldflags "-X go.mlcdf.fr/dyndns.isTest=true"