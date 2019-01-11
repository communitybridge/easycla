#! /bin/bash

# This script retrieves all credentials from the aws parameter store.

aws ssm describe-parameters --region us-east-1 | jq '."Parameters" [] ."Name"' | tr '"' ' ' | \
while read param1; do
    read param2
    read param3
    read param4
    read param5
    read param6
    read param7
    read param8
    read param9
    read param10

    aws ssm get-parameters --region us-east-1 --name $param1 $param2 $param3 $param4 $param5 $param6 $param7 $param8 $param9 $param10 | jq '."Parameters" [] | "\(.Name),\(.Value),\(.Type)"' | tr -d '"'
done
