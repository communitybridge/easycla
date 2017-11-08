#!/bin/sh

export EC2_LOCAL_IP=`curl -s http://169.254.169.254/latest/meta-data/local-ipv4`

envsubst '\\$$APP_DOMAINS \\$$EC2_LOCAL_IP' < /etc/nginx/conf.d/production.conf.tpl > /etc/nginx/conf.d/production.conf

exec /usr/sbin/nginx "$@"