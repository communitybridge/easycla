FROM python:3

RUN cd /srv/ && wget https://releases.hashicorp.com/consul-template/0.19.0/consul-template_0.19.0_linux_amd64.tgz
RUN tar -xvf /srv/consul-template_0.19.0_linux_amd64.tgz -C /usr/bin/

ADD . /srv/

RUN cd /srv/ && python3 setup.py install

CMD ["/srv/scripts/run.sh"]