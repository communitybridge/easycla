# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Application configuration options.

These values should be tracked in version control.

Please put custom non-tracked configuration options (debug mode, keys, database
configuration, etc) in cla_config.py somewhere in your Python path.
"""

import os

import logging

stage = os.environ.get('STAGE', '')

LOG_LEVEL = logging.DEBUG  #: Logging level.
#: Logging format.
LOG_FORMAT = logging.Formatter(fmt='%(asctime)s %(levelname)s %(name)s: %(message)s', datefmt='%Y-%m-%dT%H:%M:%S')

DEBUG = False  #: Debug off in production

# Base URL used for callbacks and OAuth2 redirects.
API_BASE_URL = os.environ.get('CLA_API_BASE', '')

# Contributor Console base URL
CONTRIBUTOR_BASE_URL = os.environ.get('CLA_CONTRIBUTOR_BASE', '')
CONTRIBUTOR_V2_BASE_URL = os.environ.get('CLA_CONTRIBUTOR_V2_BASE', '')

# Corporate Console base URL
CORPORATE_BASE_URL = os.environ.get('CLA_CORPORATE_BASE', '')
if CORPORATE_BASE_URL == 'corporate.prod.lfcla.com':
    CORPORATE_BASE_URL = 'corporate.lfcla.com'

# Landing Page
CLA_LANDING_PAGE = os.environ.get('CLA_LANDING_PAGE', '')

SIGNED_CALLBACK_URL = os.path.join(API_BASE_URL, 'v2/signed')  #: Default callback once signature is completed.
ALLOW_ORIGIN = '*'  # Specify the CORS Access-Control-Allow-Origin response header value.

# Define the database we are working with.
DATABASE = 'DynamoDB'  #: Database type ('SQLite', 'DynamoDB', etc).

# Define the key-value we are working with.
KEYVALUE = 'DynamoDB'  #: Key-value store type ('Memory', 'DynamoDB', etc).

# DynamoDB-specific configurations - this is applied to each table.
DYNAMO_REGION = 'us-east-1'  #: DynamoDB AWS region.
DYNAMO_WRITE_UNITS = 1  #: DynamoDB table write units.
DYNAMO_READ_UNITS = 1  #: DynamoDB table read units.

# Define the signing service to use.
SIGNING_SERVICE = 'DocuSign'  #: The signing service to use ('DocuSign', 'HelloSign', etc)

# Repository settings.
AUTO_CREATE_REPOSITORY = True  #: Create repository in database automatically on webhook.

# GitHub Repository Service.
#: GitHub OAuth2 Authorize URL.
GITHUB_OAUTH_AUTHORIZE_URL = 'https://github.com/login/oauth/authorize'
#: GitHub OAuth2 Callback URL.
GITHUB_OAUTH_CALLBACK_URL = os.path.join(API_BASE_URL, 'v2/github/installation')
#: GitHub OAuth2 Token URL.
GITHUB_OAUTH_TOKEN_URL = 'https://github.com/login/oauth/access_token'
#: How users get notified of CLA status in GitHub ('status', 'comment', or 'status+comment').
GITHUB_PR_NOTIFICATION = 'status+comment'

# GitHub Application Service.
GITHUB_APP_WEBHOOK_SECRET = 'webhook-secret'

# GitHub Oauth token used for authenticated GitHub API calls and testing
GITHUB_OAUTH_TOKEN = os.environ.get('GITHUB_OAUTH_TOKEN', '')

# Email Service.
EMAIL_SERVICE = 'SNS'  #: Email service to use for notification emails.
EMAIL_ON_SIGNATURE_APPROVED = True  #: Whether to email the user when signature has been approved.

# SMTP Configuration.
#: Sender email address for SMTP service (from address).
SMTP_SENDER_EMAIL_ADDRESS = 'test@cla.system'
SMTP_HOST = ''  #: Host of the SMTP service.
SMTP_PORT = '0'  #: Port of the SMTP service.

# Storage Service.
STORAGE_SERVICE = 'S3Storage'  #: The storage service to use for storing CLAs.

# LocalStorage Configuration.
LOCAL_STORAGE_FOLDER = '/tmp/cla'  #: Local folder when using the LocalStorage service.

# PDF Generation.
PDF_SERVICE = 'DocRaptor'
