#!/bin/bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

yarn install
cd src
npm install
cd ..
cd edge
yarn install
cd ..