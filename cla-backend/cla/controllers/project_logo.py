import os
from urllib.parse import urlsplit, urljoin

import boto3

import cla
from cla.models.s3_storage import S3Storage
from cla.auth import AuthUser, admin_list

cla_logo_url = os.environ.get('CLA_BUCKET_LOGO_URL', '')
logo_bucket_parts = urlsplit(cla_logo_url)
logo_bucket = logo_bucket_parts.path.replace('/', '')

s3_client = S3Storage()
s3_client.initialize(None, bucket_name=logo_bucket)

def upload_logo(auth_user: AuthUser, project_sfdc_id: str, data):
    if auth_user.username not in admin_list:
        return {'error': 'unauthorized'}

    cla.log.info('logo ({}) uploaded by {}'.format(project_sfdc_id, auth_user.username))

    file_path = '{}.png'.format(project_sfdc_id)
    try:
        s3_client.store(file_path, data)
    except Exception as err:
        return {'error': err}
