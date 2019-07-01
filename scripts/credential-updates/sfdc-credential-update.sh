#! /bin/bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

# This script updates the SalesForce parameters for a given environment. Only
# parameters provided in the list below are updated.

INSTANCE_URL=''
USERNAME=''
PASSWORD=''
SECURITY_TOKEN=''
CONSUMER_KEY=''
CONSUMER_SECRET=''

ENV=''
PROFILE=''

if [ -z "$ENV" ]; then
    echo "ERROR: missing environment"
    exit 1
fi

if [ -z "$PROFILE" ]; then
    echo "ERROR: missing profile"
    exit 1
fi

if [ -n "$INSTANCE_URL" ]; then
    echo "updating instance url: $INSTANCE_URL"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-sf-instance-url-$ENV" --description "SalesForce instance URL" --value "$INSTANCE_URL" --type "String" --overwrite
fi

if [ -n "$USERNAME" ]; then
    echo "updating username: $USERNAME"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-sf-username-$ENV" --description "SalesForce user name" --value "$USERNAME" --type "String" --overwrite
fi

# The SalesForce API password is the user password concatenated with a security token.
if [ -n "$PASSWORD$SECURITY_TOKEN" ]; then
    echo "updating password: $PASSWORD$SECURITY_TOKEN"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-sf-password-$ENV" --description "SalesForce password. Combined user password and secret token" --value "$PASSWORD$SECURITY_TOKEN" --type "String" --overwrite
fi

if [ -n "$CONSUMER_KEY" ]; then
    echo "updating consumer key: $CONSUMER_KEY"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-sf-consumer-key-$ENV" --description "SalesForce Connected App Consumer Key" --value "$CONSUMER_KEY" --type "String" --overwrite
fi

if [ -n "$CONSUMER_SECRET" ]; then
    echo "updating consumer secret: $CONSUMER_SECRET"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-sf-consumer-secret-$ENV" --description "SalesForce Connected App Consumer Secret" --value "$CONSUMER_SECRET" --type "String" --overwrite
fi
