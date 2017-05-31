FROM nginx

ARG build_number
ARG git_hash

LABEL build_number $build_number
LABEL hash $git_hash
LABEL maintainer "engineering@linuxfoundation.org"

RUN apt-get update -y && apt-get install curl gettext -y && rm -rf /var/lib/apt/lists/*
RUN rm -f /etc/nginx/conf.d/default.conf
COPY production.conf.tpl /etc/nginx/conf.d/production.conf.tpl
ADD https://releases.hashicorp.com/consul-template/0.18.3/consul-template_0.18.3_linux_amd64.tgz /srv/
RUN tar -xvf /srv/consul-template_0.18.3_linux_amd64.tgz -C /usr/bin/
COPY entrypoint.sh /srv/

ENTRYPOINT ["/srv/entrypoint.sh"]

CMD [ "-log-level", "debug", \
      "-consul-addr", "consul.service.consul:8500", \
      "-template", "/etc/nginx/conf.d/production.conf.ctmpl:/etc/nginx/conf.d/production.conf", \
      "-retry", "5s", \
      "-exec", "/usr/sbin/nginx -g 'daemon off;'" \
]