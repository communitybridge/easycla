# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import sys
import os
import boto3

PROMPT_CHAR = '    > '

def get_private_key_str(key_path):
    """
    Read private pem file and return the key string
    """
    return open(key_path).read()

def get_ssm_client(aws_profile, aws_region):
    """
    Init a ssm client with boto3
    """
    session = boto3.Session(profile_name=aws_profile)
    return session.client('ssm', region_name=aws_region)

def main():
    print("This script will help you provision GH app private key file to AWS parameter store.")
    print("First, enter the aws profile name it should use otherwise default will be used:")
    aws_profile = input(PROMPT_CHAR)
    print("Second, enter the aws region name it should look up:")
    aws_region = input(PROMPT_CHAR)
    print("Third, enter the stage name it should provision to:")
    stage_name = input(PROMPT_CHAR)
    print("Last, enter the relative path of private key file:")
    key_file_path = input(PROMPT_CHAR)

    ssm_client = get_ssm_client(aws_profile, aws_region)
    private_key_value = get_private_key_str(key_file_path)

    response = ssm_client.put_parameter(
        Name='cla-gh-app-private-key-{}'.format(stage_name),
        Description='A gh app private key used to sign jwt',
        Value=private_key_value,
        Overwrite=True,
        Type='SecureString')

    if response['ResponseMetadata']['HTTPStatusCode'] is 200:
        print ('Provision Completed see response below:')
    else:
        print ('Error when provision check response below:')
    print (response)

if __name__ == '__main__':
    main()
