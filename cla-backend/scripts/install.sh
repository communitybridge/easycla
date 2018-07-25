#!/usr/bin/env bash

pip3.6 install -r /srv/app/requirements.txt --user

python3.6 /srv/app/setup.py install --user

echo "Creating /home/www-data/.aws/credentials file"
mkdir -p /home/www-data/.aws
echo '[default]
aws_access_key_id=""
aws_secret_access_key=""' > /home/www-data/.aws/credentials

#python3.6 /srv/app/helpers/create_test_environment.py
python3.6 /srv/app/helpers/create_database.py
