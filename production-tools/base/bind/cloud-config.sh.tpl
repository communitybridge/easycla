#!/bin/bash

# NewRelic Infrastructure Agent
echo "license_key: ${newrelic_key}" | sudo tee -a /etc/newrelic-infra.yml
printf "[newrelic-infra]\nname=New Relic Infrastructure\nbaseurl=http://download.newrelic.com/infrastructure_agent/linux/yum/el/6/x86_64\nenable=1\ngpgcheck=0" | sudo tee -a /etc/yum.repos.d/newrelic-infra.repo
yum -q makecache -y --disablerepo='*' --enablerepo='newrelic-infra'

# Install Bind & NewRelic Infrastructure
yum install bind bind-utils newrelic-infra -y

# Install consul-template
wget https://releases.hashicorp.com/consul-template/0.18.3/consul-template_0.18.3_linux_amd64.tgz
tar xvf consul-template_0.18.3_linux_amd64.tgz
mv ./consul-template /usr/bin/

# Adding bind config
echo "" > /etc/named.conf
cat > '/etc/named.conf' <<-EOF
options {
  listen-on port 53 { any; };
  directory 	"/var/named";
  dump-file 	"/var/named/data/cache_dump.db";
  statistics-file "/var/named/data/named_stats.txt";
  memstatistics-file "/var/named/data/named_mem_stats.txt";
  allow-query     { any; };
  recursion yes;

  dnssec-enable no;
  dnssec-validation no;

  /* Path to ISC DLV key */
  bindkeys-file "/etc/named.iscdlv.key";

  managed-keys-directory "/var/named/dynamic";

  forwarders {
    ${CIDR_PREFIX}2;

    8.8.8.8;
    8.8.4.4;
  };
};

zone "." IN {
  type hint;
  file "named.ca";
};

zone "internal" IN {
  type forward;
  forward only;
  forwarders {
    ${CIDR_PREFIX}2;
  };
};

include "/etc/named.rfc1912.zones";
include "/etc/named.root.key";
include "/etc/named/consul.conf";
EOF

# Adding template for the consul zone of Bind
touch /etc/named/consul.conf.ctmpl
cat > '/etc/named/consul.conf.ctmpl' <<-EOF
zone "consul" IN {
  type forward;
  forward only;
  forwarders {
    {{range service "consul" }}{{ if .NodeAddress | contains "${CIDR_PREFIX}" }}{{.NodeAddress}} port 8600;{{end}}
    {{end}}
  };
};
EOF

# Launching Bind Server through consul-template
consul-template \
  -log-level warn \
  -consul-addr consul.prod.engineering.internal:8500 \
  -template "/etc/named/consul.conf.ctmpl:/etc/named/consul.conf" \
  -retry 5s \
  -exec="named -f -g"