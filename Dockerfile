FROM nginx

ARG build_number
ARG git_hash

LABEL build_number $build_number
LABEL hash $git_hash
LABEL maintainer "engineering@linuxfoundation.org"

RUN apt-get update -y && \
    apt-get install -y wget curl gettext gnupg && \
    curl -sL https://deb.nodesource.com/setup_6.x | bash - && \
    apt-get install -y nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir /var/www && \
    chown nginx:nginx /var/www

RUN rm -f /etc/nginx/conf.d/default.conf
COPY infra/nginx/production/production.conf.tpl /etc/nginx/conf.d/production.conf.tpl
COPY infra/nginx/production/nginx.conf /etc/nginx/nginx.conf
COPY infra/docker-prod-entrypoint.sh /srv/entrypoint.sh
COPY src /srv/app/src/
COPY scripts/constants.ts /srv/app/src/ionic/services/constants.ts
COPY scripts/keycloak.json /srv/app/src/ionic/assets/keycloak.json
RUN rm -rf /srv/app/src/node_modules /srv/app/src/www

WORKDIR '/srv/app/src/www'

RUN npm install && \
    npm rebuild node-sass && \
    npm run build

RUN chown -R nginx:nginx /srv/app/src/www

ENTRYPOINT ["/srv/entrypoint.sh"]