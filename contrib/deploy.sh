#! /usr/bin/env bash
# Setup dyndns on a Raspberry Pi 2 Model B

set -e

# download latest release of dyndns
LOCATION=$(curl -s https://api.github.com/repos/mlcdf/dyndns/releases/latest | grep browser_download_url | grep linux-arm | cut -d '"' -f 4)
echo $LOCATION
curl -L "${LOCATION}" -o dyndns

# copy to host
scp dyndns pi:/home/pi/

# add to crontab
# will run every 10 minutes
cmd='echo "*/10 * * * *" /home/pi/dyndns --domain example.com --record "*.pi"'
ssh pi "sudo sh -c '${cmd} > /etc/cron.d/dyndns'"

rm dyndns