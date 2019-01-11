"""
Storage service that stores files in AWS S3 buckets.
"""

import io
import os
import boto3
import botocore
import cla
from cla.models import storage_service_interface

stage = os.environ.get('STAGE', '')
signature_files_bucket = os.environ.get('CLA_SIGNATURE_FILES_BUCKET', '')

class S3Storage(storage_service_interface.StorageService):
    """
    Store documents in AWS S3.
    """
    def __init__(self):
        self.bucket = None

    def initialize(self, config):
        self.bucket = signature_files_bucket

    def _get_client(self):
        """Mockable method to get the S3 client."""
        if stage == 'local':
            return boto3.client('s3', endpoint_url='http://localhost:8001')

        return boto3.client('s3')

    def store(self, filename, data):
        cla.log.info('Storing filename content in S3 bucket %s: %s', self.bucket, filename)
        try:
            obj = io.BytesIO(data)
            client = self._get_client()
            client.upload_fileobj(obj, self.bucket, filename)
        except Exception as err:
            cla.log.error('Could not save filename %s in S3: %s', filename, str(err))
            raise Exception('*** Upload file failed. See details from stack traceback ^^^ ***')

    def retrieve(self, filename):
        cla.log.info('Retrieving filename content from S3: %s', filename)
        data = io.BytesIO()
        try:
            client = self._get_client()
            client.download_fileobj(self.bucket, filename, data)
            data.seek(0)
        except botocore.exceptions.ClientError as err:
            cla.log.error('Client error while retrieving file from S3 %s: %s', filename, str(err))
        except Exception as err:
            cla.log.error('Unknown error while retrieving file from S3 %s: %s', filename, str(err))
        return data.read()

    def delete(self, filename):
        cla.log.info('Deleting from S3 storage: %s', filename)
        try:
            client = self._get_client()
            client.delete_object(Bucket=self.bucket, Key=filename)
        except Exception as err:
            cla.log.error('Error while deleting filename %s in S3: %s', filename, str(err))

class MockS3Storage(S3Storage):
    """Mock AWS S3 storage model."""
    def _get_client(self):
        return MockS3StorageClient()

    def _create_bucket(self, client=None):
        pass

class MockS3StorageClient(object):
    """Mock AWS S3 storage client."""
    def __init__(self, buckets=None):
        if buckets is None:
            self.buckets = {'Buckets': [{'Name': 'Test Bucket'}]}
        else:
            self.buckets = buckets

    def list_buckets(self):
        """Mock method for listing S3 bucket information."""
        return self.buckets

    def download_fileobj(self, bucket, filename, data): # pylint: disable=unused-argument,no-self-use
        """Mock method for downloading S3 file object data."""
        with open('resources/test.pdf', 'rb') as fhandle:
            data.write(fhandle.read())
        return data
