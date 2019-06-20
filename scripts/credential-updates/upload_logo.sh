#!/bin/bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

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

curl -v -XPUT \
    -H "Content-Type: image/png" \
    -T "$file_path" \
    "$signed_url";