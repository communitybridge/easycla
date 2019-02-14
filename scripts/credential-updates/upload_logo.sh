#!/bin/bash

# file_path
# project_sfdc_id
# bearer_token
# base_url

echo "file_path $file_path"
echo "project_sfdc_id $project_sfdc_id"
# echo "bearer_token $bearer_token"
echo "base_url $base_url"

signed_url_resp=$(curl -s -XGET \
    -H "Authorization: Bearer $bearer_token" \
    "$base_url/v1/project/logo/$project_sfdc_id");

signed_url=$(echo $signed_url_resp | jq -r '.signed_url');

echo "Signed url: $signed_url";

curl -XPUT \
    -H "Content-Type: image/png" \
    --data @$file_path \
    "$signed_url";