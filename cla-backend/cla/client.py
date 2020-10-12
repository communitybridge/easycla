# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Python Client that wraps boto3 client and exposes utility functions
"""

from boto3 import client
import cla

def get_ssm_key(region, key):
    """
    Fetches the specified SSM key value from the SSM key store
    """
    ssm_client = client('ssm', region_name=region)
    cla.log.debug(f'Fetching Key: {key}')
    response = ssm_client.get_parameter(Name=key, WithDecryption=True)
    cla.log.debug(f'Fetched Key: {key}, value: {response["Parameter"]["Value"]}')
    return response['Parameter']['Value']