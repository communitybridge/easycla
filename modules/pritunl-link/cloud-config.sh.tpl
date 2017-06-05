#!/bin/bash

# Making sure we are always using Amazon DNS on pritunl-link instances
echo "supersede domain-name-servers ${DNS_SERVER};" >> /etc/dhcp/dhclient.conf
service networking restart

# Pritunl-Link Installation
tee -a /etc/apt/sources.list.d/pritunl.list << EOF
deb http://repo.pritunl.com/stable/apt xenial main
EOF
apt-key adv --keyserver hkp://keyserver.ubuntu.com --recv 7568D9BB55FF9E5287D586017AE645C0CF8E292A
apt-get update
apt-get --assume-yes upgrade
apt-get --assume-yes --allow-unauthenticated install pritunl-link
pritunl-link verify-off
pritunl-link provider aws
pritunl-link add ${PRITUNL_LINK}