#!/bin/bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

###############################################################################
# What to do if we die
###############################################################################
die() {
  echo "$@"
  exit 1
}

echo "Running yarn install in folder: $(pwd)"
yarn install

for dir in src edge; do
  if [[ -d ${dir} ]]; then
    pushd ${dir} || die "Unable to change to pushd to the ${dir} folder - exiting..."
    echo "Running yarn install in folder: $(pwd)"
    yarn install
    popd || die "Unable to change to popd from the ${dir} folder - exiting..."
  else
    echo "Missing folder: ${dir} - skipping yarn install"
  fi
done
