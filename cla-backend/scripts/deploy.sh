#!/usr/bin/env bash
set -e

usage () { 
    echo "Usage : $0 -s <stage>"
}

while getopts ":s:c" opts; do
    case ${opts} in
        s) STAGE=${OPTARG} ;;
    esac
done

if [ -z "${STAGE}" ]; then
    usage
    exit 1
fi

./node_modules/.bin/serverless deploy --stage="${STAGE}" --region us-east-1

