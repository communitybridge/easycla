#!/usr/bin/env bash

cd /srv/

# get current instance ip
INSTANCE_IP=`curl -s http://169.254.169.254/latest/meta-data/local-ipv4`

consul-template \
  -log-level warn \
  -consul-addr ${INSTANCE_IP}:8500 \
  -template "/srv/infra/cla_config.py.ctmpl:/srv/cla_config.py" \
  -once

NEW_RELIC_CONFIG_FILE=/srv/infra/newrelic/newrelic.ini newrelic-admin run-program hug -f /srv/cla/routes.py -p 5000