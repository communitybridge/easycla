FROM 433610389961.dkr.ecr.us-west-2.amazonaws.com/base:latest

MAINTAINER Linux Foundation <webmaster@linuxfoundation.org>

RUN curl --silent --location https://rpm.nodesource.com/setup_6.x | bash - && yum install -y nodejs gcc-c++ make

ADD https://releases.hashicorp.com/consul-template/0.18.3/consul-template_0.18.3_linux_amd64.tgz /srv/
RUN tar -xvf /srv/consul-template_0.18.3_linux_amd64.tgz -C /usr/bin/ && \
    rm -f /srv/consul-template_0.18.3_linux_amd64.tgz

RUN useradd www-data && \
    usermod -u 1000 --shell /bin/bash www-data

COPY infra/docker-prod-entrypoint.sh /srv/entrypoint.sh

COPY . /srv/app/

RUN chown -R www-data:www-data /srv/app

USER www-data

WORKDIR '/srv/app/src'

RUN npm install && \
    npm rebuild node-sass && \
    npm run build && \
    rm -rf /srv/app/src/config/default.json && \
    rm -rf /srv/app/src/newrelic.js

ENTRYPOINT ["/srv/entrypoint.sh"]

CMD ["start"]