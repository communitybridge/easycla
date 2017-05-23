#!/bin/sh

consul-template \
  -log-level warn \
  -consul-addr us-west-2.logstash.service.consul:8500 \
  -template "/srv/app/src/config/production.json.ctmpl:/srv/app/src/config/default.json" \
  -once

consul-template \
  -log-level warn \
  -consul-addr us-west-2.logstash.service.consul:8500 \
  -template "/srv/app/src/newrelic.js.ctmpl:/srv/app/src/newrelic.js" \
  -once

exec npm "$@"