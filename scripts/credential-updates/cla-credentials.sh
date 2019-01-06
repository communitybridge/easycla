#! /bin/bash

# This script updates additional vars needed for the CLA
# projects. Only parameters provided in the list below are updated.

CLA_API_URL='';
CORP_CONSOLE_LINK='';
CLA_BUCKET_LOGO_URL='';
CLA_SIGNATURE_FILES_BUCKET='';
SES_SENDER_EMAIL_ADDRESS='';

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

if [ -n "$CLA_API_URL" ]; then
    echo "updating API URL: $CLA_API_URL"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-cla-api-url-$ENV" --description "CLA API URL" --value "$CLA_API_URL" --type "String" --overwrite
fi

if [ -n "$CORP_CONSOLE_LINK" ]; then
    echo "updating corp console link: $CORP_CONSOLE_LINK"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-corp-console-link-$ENV" --description "Corporation Console Link" --value "$CORP_CONSOLE_LINK" --type "String" --overwrite
fi

if [ -n "$CLA_BUCKET_LOGO_URL" ]; then
    echo "updating cla logo url: $CLA_BUCKET_LOGO_URL"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-cla-logo-s3-url-$ENV" --description "CLA Logo S3 URL" --value "$CLA_BUCKET_LOGO_URL" --type "String" --overwrite
fi

if [ -n "$CLA_SIGNATURE_FILES_BUCKET" ]; then
    echo "updating cla signature files bucket: $CLA_SIGNATURE_FILES_BUCKET"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-signature-files-bucket-$ENV" --description "CLA signature files bucket" --value "$CLA_SIGNATURE_FILES_BUCKET" --type "String" --overwrite
fi

if [ -n "$SES_SENDER_EMAIL_ADDRESS" ]; then
    echo "updating ses sender email address: $SES_SENDER_EMAIL_ADDRESS"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-ses-sender-email-address-$ENV" --description "SES Sender Email Address" --value "$SES_SENDER_EMAIL_ADDRESS" --type "String" --overwrite
fi
