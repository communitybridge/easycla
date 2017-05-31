#!/bin/bash
tee -a /etc/apt/sources.list.d/pritunl.list << EOF
deb http://repo.pritunl.com/stable/apt xenial main
EOF
apt-key adv --keyserver hkp://keyserver.ubuntu.com --recv 7568D9BB55FF9E5287D586017AE645C0CF8E292A
apt-get update
apt-get --assume-yes install pritunl-link
pritunl-link verify-off
pritunl-link provider aws
pritunl-link add pritunl://592f0fbfb8181a0a1cf5b91d:lZdDrSLfcQTRwxpj9lMGvQ0KQNWaQcJ0@vpn.engineering.tux.rocks