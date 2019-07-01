# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import os
from urllib.parse import urlsplit, urljoin

import boto3

import cla
from cla.auth import AuthUser, admin_list

stage = os.environ.get('STAGE', '')

cla_logo_url = os.environ.get('CLA_BUCKET_LOGO_URL', '')
logo_bucket_parts = urlsplit(cla_logo_url)
logo_bucket = logo_bucket_parts.path.replace('/', '')

endpoint_url = None
if stage == 'local':
    endpoint_url = 'http://localhost:8001'

s3_client = boto3.client('s3', endpoint_url=endpoint_url)

def create_signed_logo_url(auth_user: AuthUser, project_sfdc_id: str):
    if auth_user.username not in admin_list:
        return {'error': 'unauthorized'}

    cla.log.info('signed url ({}) created by {}'.format(project_sfdc_id, auth_user.username))

    file_path = '{}.png'.format(project_sfdc_id)

    params = {
        'Bucket': logo_bucket,
        'Key': file_path,
        'ContentType': 'image/png'
    }

    try:
        signed_url = s3_client.generate_presigned_url('put_object', Params=params, ExpiresIn=300, HttpMethod='PUT')
        return {'signed_url': signed_url}
    except Exception as err:
        return {'error': err}
