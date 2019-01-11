#! /bin/bash

# This script updates the Github parameters for a given environment. Only
# parameters provided in the list below are updated.

GH_APP_ID=''
GH_OAUTH_CLIENT_ID=''
GH_OAUTH_SECRET=''
GH_APP_PUBLIC_LINK=''
GH_APP_WEBHOOK_SECRET=''
GH_APP_PRIVATE_KEY_PATH='' # This is a filename

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

if [ -n "$GH_APP_ID" ]; then
    echo "updating app ID: $GH_APP_ID"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-gh-app-id-$ENV" --description "Github App ID" --value "$GH_APP_ID" --type "String" --overwrite
fi

if [ -n "$GH_OAUTH_CLIENT_ID" ]; then
    echo "updating oauth client ID: $GH_OAUTH_CLIENT_ID"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-gh-oauth-client-id-$ENV" --description "Github oauth client ID" --value "$GH_OAUTH_CLIENT_ID" --type "String" --overwrite
fi

if [ -n "$GH_APP_PUBLIC_LINK" ]; then
    echo "updating public link: $GH_APP_PUBLIC_LINK"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-gh-app-public-link-$ENV" --description "Github app public link" --value "$GH_APP_PUBLIC_LINK" --type "String" --overwrite
fi

if [ -n "$GH_OAUTH_SECRET" ]; then
    echo "updating oauth secret: $GH_OAUTH_SECRET"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-gh-oauth-secret-$ENV" --description "Github oauth secret" --value "$GH_OAUTH_SECRET" --type "String" --overwrite
fi

if [ -n "$GH_APP_WEBHOOK_SECRET" ]; then
    echo "updating webhook secret: $GH_APP_WEBHOOK_SECRET"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-gh-app-webhook-secret-$ENV" --description "Github webhook secret" --value "$GH_APP_WEBHOOK_SECRET" --type "String" --overwrite
fi

if [ -n "$GH_APP_PRIVATE_KEY_PATH" ]; then
    echo "updating private key: $GH_APP_PRIVATE_KEY_PATH"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-gh-app-private-key-$ENV" --description "Github private key" --value "file://$GH_APP_PRIVATE_KEY_PATH" --type "SecureString" --overwrite
fi
