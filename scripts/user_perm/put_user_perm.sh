#! /bin/bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

# This script uploads a user and accompanying permissions to
# the DynamoDB cla-{env}-user-permissions table

set -e

help () {
  echo "Usage : $0 -s <stage> -r <region> -u '<auth0-user-id>' -p '<\"project1\",\"project2\">' -c <\"company1\",\"company2\">"
}

# Get stage, user, projects and companies
while [ "$#" -gt 0 ]; do
    case $1 in
        -h|--help) help; exit 1 ;;
        -s|--stage) STAGE=$2 ; shift; shift;;
        -r|--region) REGION=$2 ; shift; shift ;;
        -u|--user) USERNAME=$2 ; shift; shift ;;
        -p|--projects) PROJECTS=$2 ; shift; shift ;;
        -c|--companies) COMPANIES=$2 ; shift; shift ;;
        *) printf "invalid parameter: $1\n" >&2 ; exit 1 ;;
    esac
done

if [ -x "$STAGE" ]; then
    echo "ERROR: missing stage"
fi

if [ -x "$USERNAME" ]; then
    echo "ERROR: missing username"
fi

if [ -z "$REGION" ]; then
    echo "ERROR: missing region"
    exit 1
fi

if [ -z "$PROJECTS" ] && [ -z "$COMPANIES" ]; then
    echo "ERROR: projects and companies cannot both be empty"
    exit 1
fi

USER="{ \"username\": { \"S\": \"$USERNAME\" } ";
if [ -n "$PROJECTS" ]; then
    USER="$USER, \"projects\": { \"SS\": [$PROJECTS] } "
fi
if [ -n "$COMPANIES" ]; then
    USER="$USER ,\"companies\": { \"SS\": [$COMPANIES]} "
fi
USER="$USER }"

aws dynamodb put-item \
    --table-name "cla-$STAGE-user-permissions" \
    --item "$USER" \
    --region "$REGION" 
    #--endpoint-url http://localhost:8000 #local dev env
