#!/bin/bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

# This script copies SSM env vars from a source AWS account
# to a destination account.
#
# Set 'DRY_RUN=true' to print the copy commands, instead of
# running them directly.

SRC_PROFILE='';
DST_PROFILE='';

VAR_FILE='env_vars.txt';

for var_name in $(cat $VAR_FILE); do
    ssm_var=$(
        AWS_PROFILE=$SRC_PROFILE \
        aws ssm get-parameter \
            --name "$var_name" | jq '.Parameter'
    );

    name=$(echo $ssm_var | jq '.Name')
    value=$(echo $ssm_var | jq '.Value')
    type=$(echo $ssm_var | jq '.Type')

    command='AWS_PROFILE='$DST_PROFILE' aws ssm put-parameter --name '"$name"' --value '"$value"' --type '"$type"' --overwrite'

    if [ "$DRY_RUN" = true ]; then
        echo "$command"
    else
        echo "$name"
        eval $command
    fi
done
