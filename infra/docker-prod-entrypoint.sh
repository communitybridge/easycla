#!/bin/sh

# get current instance ip
INSTANCE_IP=`curl -s http://169.254.169.254/latest/meta-data/local-ipv4`

rm -f /srv/newrelic/newrelic.ini

consul-template \
  -log-level warn \
  -consul-addr ${INSTANCE_IP}:8500 \
  -template "/srv/app/newrelic/newrelic.ini.ctmpl:/srv/app/newrelic/newrelic.ini" \
  -once

consul-template \
  -log-level warn \
  -consul-addr ${INSTANCE_IP}:8500 \
  -template "/srv/app/src/config/cla_config.py.ctmpl:/srv/app/src/cla_config.py" \
  -once

export NEW_RELIC_CONFIG_FILE=/srv/app/newrelic/newrelic.ini

exec ~/.local/bin/newrelic-admin run-program "$@"