#!/usr/bin/env bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

set -e

usage() {
  echo "Usage : $0 -s <stage> -r <region of api> [-c](enable cloudfront)"
}

# Get STAGE and CLOUDFRONT configuration from command line.
CLOUDFRONT=false
while getopts ":s:r:c" opts; do
  case ${opts} in
  s) STAGE=${OPTARG} ;;
  r) REGION=${OPTARG} ;;
  c) CLOUDFRONT=true ;;
  *) break ;;
  esac
done
# Removes the parsed command line opts
shift $((OPTIND - 1))

if [ -z "${STAGE}" ]; then
  usage
  exit 1
fi

if [ -z "${REGION}" ]; then
  usage
  exit 1
fi

# Allow nvm to work within a script
. /usr/local/opt/nvm/nvm.sh

echo 'Building Distribution'
cd src
echo "Switching to node v8.17.0..."
nvm use v8.17.0
node --version
yarn build:${STAGE}
cd ../

echo 'Building Edge Function'
cd edge
echo "Switching to node v8.17.0..."
nvm use v8.17.0
node --version
yarn build
cd ../

echo "Switching to node v12.14.1..."
nvm use v12.14.1
node --version
echo 'Deploying Cloudfront and lambda@edge'
yarn sls deploy --stage="${STAGE}" --cloudfront="${CLOUDFRONT}"

echo 'Deploying Frontend Bucket'
yarn sls client deploy --stage="${STAGE}" --cloudfront="${CLOUDFRONT}" --no-confirm --no-policy-change --no-config-change

if [ ${CLOUDFRONT} = true ]; then
  echo 'Invalidating Cloudfront'
  yarn sls cloudfrontInvalidate --stage="${STAGE}" --cloudfront="${CLOUDFRONT}"
fi

exit 0
