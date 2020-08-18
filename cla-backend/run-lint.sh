#!/usr/bin/env bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

echo 'Running: flake8 --count --ignore=E501 --show-source --statistics *.py'
flake8 --count --ignore=E501 --show-source --statistics */**.py

echo 'Running: flake8 --ignore=E501 --count --exit-zero --max-complexity=10 --max-line-length=127 --statistics *.py'
flake8 --ignore=E501 --count --exit-zero --max-complexity=10 --max-line-length=127 --statistics */**.py
