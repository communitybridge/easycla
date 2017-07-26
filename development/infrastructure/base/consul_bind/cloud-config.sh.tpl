#!/bin/bash

# Launching Bind Server through consul-template
consul-template \
  -log-level warn \
  -consul-addr 127.0.0.1:8500 \
  -template "/etc/named/consul.conf.ctmpl:/etc/named/consul.conf" \
  -retry 5s \
  -exec="named -f -g"