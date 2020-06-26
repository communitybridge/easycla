#!/usr/bin/env bash

echo 'Running: flake8 --count --ignore=E501 --show-source --statistics *.py'
flake8 --count --ignore=E501 --show-source --statistics */**.py

echo 'Running: flake8 --ignore=E501 --count --exit-zero --max-complexity=10 --max-line-length=127 --statistics *.py'
flake8 --ignore=E501 --count --exit-zero --max-complexity=10 --max-line-length=127 --statistics */**.py
