#! /bin/bash

# This script uploads a user and accompanying permissions to
# the DynamoDB cla-{env}-user-permissions table

# This script assumes

USER='{
    "user_id": { "S": "" },
    "projects": { "SS": ["PROJECT_ID"] },
    "companies": { "SS": ["COMPANY_ID"]}
}'

ENV='';

if [ -z "$ENV" ]; then
    echo "ERROR: missing environment"
    exit 1
fi

if [ "$ENV" == "prod" ]; then
    echo "Are you sure you want to update a production user?"
    exit 1
fi

aws dynamodb put-item \
    --table-name "cla-$ENV-user-permissions" \
    --item "$USER" \
    --profile lf-cla --region us-east-1