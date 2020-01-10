#!/bin/bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

yarn install
cd src
npm install
cd ..
cd edge
yarn install
cd ..