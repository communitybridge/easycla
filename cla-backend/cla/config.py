"""
Application configuration options.

These values should be tracked in version control.

Please put custom non-tracked configuration options (debug mode, keys, database
configuration, etc) in cla_config.py somewhere in your Python path.
"""

import logging
import os

stage = os.environ.get('STAGE', '')

LOG_LEVEL = logging.INFO #: Logging level.
#: Logging format.
LOG_FORMAT = logging.Formatter('%(asctime)s %(levelname)-8s %(name)s: %(message)s')

DEBUG = False #: Debug off in production

# Base URL used for callbacks and OAuth2 redirects.
if stage == 'prod':
    BASE_URL = 'https://ckr858t6zb.execute-api.us-east-1.amazonaws.com/prod'
elif stage == 'staging':
    BASE_URL = 'https://wbkv5r3eyf.execute-api.us-east-1.amazonaws.com/staging'
elif stage == 'qa':
    BASE_URL = 'https://pd9alkfok2.execute-api.us-east-1.amazonaws.com/qa'
else:
    BASE_URL = '' #ADD YOUR DEV STAGE STAGE HERE
SIGNED_CALLBACK_URL = BASE_URL + '/v1/signed' #: Default callback once signature is completed.
ALLOW_ORIGIN = '*' # Specify the CORS Access-Control-Allow-Origin response header value.

# Define the database we are working with.
DATABASE = 'DynamoDB' #: Database type ('SQLite', 'DynamoDB', etc).

# Define the key-value we are working with.
KEYVALUE = 'DynamoDB' #: Key-value store type ('Memory', 'DynamoDB', etc).

# Endpoint where users end up to start the signing workflow.
if stage == 'prod':
    CLA_CONSOLE_ENDPOINT = 'https://d1fivluqxpmxmf.cloudfront.net'
elif stage == 'staging':
    CLA_CONSOLE_ENDPOINT = 'https://d7pvqqazh4kg5.cloudfront.net'
elif stage == 'qa':
    CLA_CONSOLE_ENDPOINT = 'https://d37jq4fjnidrq1.cloudfront.net'
else:
    CLA_CONSOLE_ENDPOINT = 'http://localhost:8100' #MODIFY HERE IF DEPLOYING TO DEV STAGE

# Define the signing service to use.
SIGNING_SERVICE = 'DocuSign' #: The signing service to use ('DocuSign', 'HelloSign', etc)

# Repository settings.
AUTO_CREATE_REPOSITORY = True #: Create repository in database automatically on webhook.

# GitHub Repository Service.
#: GitHub OAuth2 Authorize URL.
GITHUB_OAUTH_AUTHORIZE_URL = 'https://github.com/login/oauth/authorize'
#: GitHub OAuth2 Callback URL.
GITHUB_OAUTH_CALLBACK_URL = BASE_URL + '/v2/github/installation'
#: GitHub OAuth2 Token URL.
GITHUB_OAUTH_TOKEN_URL = 'https://github.com/login/oauth/access_token'
#: How users get notified of CLA status in GitHub ('status', 'comment', or 'status+comment').
GITHUB_PR_NOTIFICATION = 'status+comment'

# GitHub Application Service.
GITHUB_APP_WEBHOOK_SECRET = 'webhook-secret'

# GitLab Repository Service.
GITLAB_DOMAIN = 'https://<gitlab-domain>' #: URL to GitLab instance.
GITLAB_TOKEN = 'token' #: GitLab personal access token for the CLA system user.
GITLAB_CLIENT_ID = 'client_id' #: GitLab OAuth2 client ID.
GITLAB_SECRET = 'secret' #: GitLab OAuth2 secret.
#: GitLab OAuth2 Authorize URL.
GITLAB_OAUTH_AUTHORIZE_URL = 'https://<gitlab-domain>/oauth/authorize'
#: GitLab OAuth2 Token URL.
GITLAB_OAUTH_TOKEN_URL = 'https://<gitlab-domain>/oauth/token'
#: How users get notified of CLA status in GitLab ('status', 'comment', or 'status+comment').
GITLAB_MR_NOTIFICATION = 'status+comment'

# Email Service.
EMAIL_SERVICE = 'SMTP' #: Email service to use for notification emails.
EMAIL_ON_SIGNATURE_APPROVED = True #: Whether to email the user when signature has been approved.

# SMTP Configuration.
#: Sender email address for SMTP service (from address).
SMTP_SENDER_EMAIL_ADDRESS = 'test@cla.system'
SMTP_HOST = '' #: Host of the SMTP service.
SMTP_PORT = '0' #: Port of the SMTP service.

# AWS SES Configuration.
SES_SENDER_EMAIL_ADDRESS = 'test@cla.system' #: SES sender email address - must be verified in SES.
SES_REGION = 'us-east-1' #: The AWS region out of which the emails will be sent.
SES_ACCESS_KEY = None #: AWS Access Key ID for the SES service.
SES_SECRET_KEY = None #: AWS Secret Access Key for the SES service.

# Storage Service.
STORAGE_SERVICE = 'S3Storage' #: The storage service to use for storing CLAs.

# LocalStorage Configuration.
LOCAL_STORAGE_FOLDER = '/tmp/cla' #: Local folder when using the LocalStorage service.

# S3Storage Configuration.
# ADD KEYS IF DEPLOYING TO DEV STAGE
S3_ACCESS_KEY = '' #: AWS Access Key ID for the S3 service.
S3_SECRET_KEY = '' #: AWS Secret Access Key for the S3 service.

# PDF Generation.
PDF_SERVICE = 'DocRaptor'
