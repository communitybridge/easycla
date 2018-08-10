#!/usr/bin/env bash
set -e

[ -z "$DOMAIN" ] && echo "Need to set DOMAIN environment variable" && exit 1;
[ -z "$EMAIL"] && echo "Need to set EMAIL environment variable" && exit 1;

# Time to run certbot, to get a TLS cert for our domain
sudo certbot --nginx --redirect --domain "${DOMAIN}" --non-interactive --agree-tos --email "${EMAIL}"
sudo systemctl restart nginx
sudo systemctl restart drone
