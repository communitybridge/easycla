#!/usr/bin/env bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier=MIT

docker run \
  --rm \
  -it \
  --name easycla-python-bash \
  --entrypoint /bin/bash \
  726224182707.dkr.ecr.us-east-1.amazonaws.com/lfx-easycla-dev:latest


