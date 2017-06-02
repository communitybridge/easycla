#!/bin/bash
tee -a /etc/apt/sources.list.d/pritunl.list << EOF
deb http://repo.pritunl.com/stable/apt xenial main
EOF
apt-key adv --keyserver hkp://keyserver.ubuntu.com --recv 7568D9BB55FF9E5287D586017AE645C0CF8E292A
apt-get update
apt-get --assume-yes install pritunl-link
pritunl-link verify-off
pritunl-link provider aws
pritunl-link add pritunl://592ef6a7b8181a0a1cf53601:K4q0MCLtyLeM7DkJ50uPB5bjcvDK5z5a@vpn.engineering.tux.rocks