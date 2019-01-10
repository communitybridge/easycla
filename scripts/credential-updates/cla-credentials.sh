#! /bin/bash

# This script updates additional vars needed for the CLA
# projects. Only parameters provided in the list below are updated.

CLA_API_URL='';
CORP_CONSOLE_LINK='';
CLA_BUCKET_LOGO_URL='';
CLA_SIGNATURE_FILES_BUCKET='';
SES_SENDER_EMAIL_ADDRESS='';
DOCRAPTOR_TEST_MODE='';

# LFID LDAP OAuth2 Credentials
LF_GROUP_CLIENT_ID='';
LF_GROUP_CLIENT_SECRET='';
LF_GROUP_REFRESH_TOKEN='';
LF_GROUP_CLIENT_URL='';

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

if [ -n "$DOCRAPTOR_TEST_MODE" ]; then
    echo "updating docraptor test mode: $DOCRAPTOR_TEST_MODE"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-docraptor-test-mode-$ENV" --description "Docraptor test mode" --value "$DOCRAPTOR_TEST_MODE" --type "String" --overwrite
fi

if [ -n "$LF_GROUP_CLIENT_ID" ]; then
    echo "updating lf group client id: $LF_GROUP_CLIENT_ID"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-lf-group-client-id-$ENV" --description "LFID Group Client ID" --value "$LF_GROUP_CLIENT_ID" --type "String" --overwrite
fi

if [ -n "$LF_GROUP_CLIENT_SECRET" ]; then
    echo "updating lf group client secret: $LF_GROUP_CLIENT_SECRET"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-lf-group-client-secret-$ENV" --description "LFID Group Client Secret" --value "$LF_GROUP_CLIENT_SECRET" --type "String" --overwrite
fi

if [ -n "$LF_GROUP_REFRESH_TOKEN" ]; then
    echo "updating lf group refresh token: $LF_GROUP_REFRESH_TOKEN"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-lf-group-refresh-token-$ENV" --description "LFID Group Refresh Token" --value "$LF_GROUP_REFRESH_TOKEN" --type "String" --overwrite
fi

if [ -n "$LF_GROUP_CLIENT_URL" ]; then
    echo "updating lf group client url: $LF_GROUP_CLIENT_URL"
    aws ssm put-parameter --profile $PROFILE --region us-east-1 --name "cla-lf-group-client-url-$ENV" --description "LFID Group Client URL" --value "$LF_GROUP_CLIENT_URL" --type "String" --overwrite
fi
