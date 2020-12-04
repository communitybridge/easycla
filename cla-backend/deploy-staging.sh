#!/usr/bin/env bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

# The golang lambda file list
declare -a lambdas=("backend-aws-lambda"
  "user-subscribe-lambda"
  "metrics-aws-lambda"
  "dynamo-events-lambda"
  "zipbuilder-scheduler-lambda"
  "zipbuilder-lambda")

echo "Installing dependencies..."
yarn install

missing_lambda=0
echo "Testing if the lambdas have been copied over..."
for i in "${lambdas[@]}"; do
  echo "Testing lambda file: ${i}..."
  if [[ ! -f "${i}" ]]; then
    echo "MISSING - lambda file: ${i}"
    missing_lambda=1
  else
    echo "PRESENT - lambda file: ${i}"
  fi
done

if [[ ${missing_lambda} -ne 0 ]]; then
  echo "Missing one or more lambda files - building golang binaries in 5 seconds..."
  sleep 5
  pushd "../cla-backend-go" || exit
  make all-linux
  popd || exit
  echo "Copying over files..."
  cp "${lambdas[@]}" .
else
  echo "All golang lambda files present."
fi

for i in "${lambdas[@]}"; do
  echo "Testing file: ${i}..."
  if ! diff -q "../cla-backend-go/${i}" "${i}" &>/dev/null; then
    echo "Golang file differs: ../cla-backend-go/${i} ${i}"
    exit 1
  fi
done

time yarn deploy:staging
