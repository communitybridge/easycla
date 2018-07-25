#!/usr/bin/env bash

echo "Creating ~/.aws/credentials file"
mkdir -p ~/.aws
echo '[default]
aws_access_key_id=""
aws_secret_access_key=""' > ~/.aws/credentials