#! /bin/bash

# This script updates the Auth0 parameters for a given environment. Only
# parameters provided in the list below are updated.

AUTH0_DOMAIN='';
AUTH0_CLIENT_ID='';

ENV='';
PROFILE='';

if [ -z "$ENV" ]; then
    echo "ERROR: missing environment"
    exit 1
fi

if [ -z "$PROFILE" ]; then
    echo "ERROR: missing profile"
    exit 1
fi

if [ -n "$AUTH0_DOMAIN" ]; then
    echo "updating Auth0 Domain: $AUTH0_DOMAIN"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-auth0-domain-$ENV" --description "Auth0 Domain" --value "$AUTH0_DOMAIN" --type "String" --overwrite
fi

if [ -n "$AUTH0_CLIENT_ID" ]; then
    echo "updating Auth0 Client ID: $AUTH0_CLIENT_ID"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-auth0-clientId-$ENV" --description "Auth0 Client ID" --value "$AUTH0_CLIENT_ID" --type "String" --overwrite
fi
