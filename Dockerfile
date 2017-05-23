FROM 433610389961.dkr.ecr.us-west-2.amazonaws.com/base:latest

MAINTAINER Linux Foundation <webmaster@linuxfoundation.org>

RUN yum install wget nodejs git gettext python python-pip -y

RUN useradd www-data

RUN usermod -u 1000 --shell /bin/bash www-data

RUN curl --silent --location https://rpm.nodesource.com/setup_6.x | bash -

RUN pip install awscli

USER www-data

COPY . /srv/app/

RUN chown -R www-data:www-data /srv/app

RUN cd /srv/app/src && npm install

RUN cd /srv/app/src && npm run build

WORKDIR '/srv/app/src'

CMD ["npm", "start"]