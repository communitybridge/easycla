FROM python:3

ENV PATH "$PATH:/usr/local/bin/"

RUN cd /srv/ && wget https://releases.hashicorp.com/consul-template/0.19.0/consul-template_0.19.0_linux_amd64.tgz
RUN tar -xvf /srv/consul-template_0.19.0_linux_amd64.tgz -C /usr/bin/

ADD . /srv/

RUN cd /srv/ && pip3 install -r requirements.txt

CMD ["/srv/scripts/run.sh"]