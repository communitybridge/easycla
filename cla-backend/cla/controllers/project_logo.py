import os
from urllib.parse import urlsplit, urljoin

import boto3

import cla
from cla.auth import AuthUser, admin_list

cla_logo_url = os.environ.get('CLA_BUCKET_LOGO_URL', '')
logo_bucket_parts = urlsplit(cla_logo_url)
logo_bucket = logo_bucket_parts.path.replace('/', '')

s3_client = boto3.client('s3')

def create_signed_logo_url(auth_user: AuthUser, project_sfdc_id: str):
    if auth_user.username not in admin_list:
        return {'error': 'unauthorized'}

    cla.log.info('signed url ({}) created by {}'.format(project_sfdc_id, auth_user.username))

    file_path = '{}.png'.format(project_sfdc_id)

    params = {
        'Bucket': logo_bucket,
        'Key': file_path
    }

    try:
        signed_url = s3_client.generate_presigned_url('put_object', Params=params, ExpiresIn=100)
        return {'signed_url': signed_url}
    except Exception as err:
        return {'error': err}
