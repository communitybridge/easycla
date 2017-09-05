#!/usr/bin/env bash

yum install libxslt-devel libxml2-devel gcc python36u python36u-pip python36u-devel python36u-dev -y

pip3.6 install -r /srv/app/requirements.txt

python3.6 /srv/app/setup.py install

python3.6 /srv/app/helpers/create_database.py