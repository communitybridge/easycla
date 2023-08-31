#!/usr/bin/env bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier=MIT

podman run \
  --rm \
  -it \
  --name easycla-python-bash \
  --entrypoint /bin/bash \
  easycla-python:latest


