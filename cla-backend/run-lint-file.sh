#!/usr/bin/env bash

if [[ $# -eq 0 ]]; then
  echo "Expecting file path as input"
  echo "USAGE   : ${0} <input_file.py>"
  echo "EXAMPLE : ${0} cla/controlers/user.py"
  exit 1
fi

echo "Running: flake8 --count --ignore=E501 --show-source --statistics ${1}"
flake8 --count --ignore=E501 --show-source --statistics ${1}

echo "Running: flake8 --ignore=E501 --count --exit-zero --max-complexity=10 --max-line-length=127 --statistics ${1}"
flake8 --ignore=E501 --count --exit-zero --max-complexity=10 --max-line-length=127 --statistics ${1}
