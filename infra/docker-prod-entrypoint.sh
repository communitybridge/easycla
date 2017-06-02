#!/bin/sh

# get current instance ip
INSTANCE_IP=`curl -s http://169.254.169.254/latest/meta-data/local-ipv4`

consul-template \
  -log-level warn \
  -consul-addr ${INSTANCE_IP}:8500 \
  -template "/srv/app/src/config/production.json.ctmpl:/srv/app/src/config/default.json" \
  -once

consul-template \
  -log-level warn \
  -consul-addr ${INSTANCE_IP}:8500 \
  -template "/srv/app/src/newrelic.js.ctmpl:/srv/app/src/newrelic.js" \
  -once

exec npm "$@"