"""
Application configuration options.

These values should be tracked in version control.

Please put custom non-tracked configuration options (debug mode, keys, database
configuration, etc) in cla_config.py somewhere in your Python path.
"""

import logging

LOG_LEVEL = logging.INFO #: Logging level.
#: Logging format.
LOG_FORMAT = logging.Formatter('%(asctime)s %(levelname)-8s %(name)s: %(message)s')

DEBUG = False #: Debug off in production

BASE_URL = 'https://k0dcklbzoh.execute-api.us-east-1.amazonaws.com/dev-runze' #: Base URL used for callbacks and OAuth2 redirects.
SIGNED_CALLBACK_URL = BASE_URL + '/v1/signed' #: Default callback once signature is completed.
ALLOW_ORIGIN = '*' # Specify the CORS Access-Control-Allow-Origin response header value.

# Define the database we are working with.
DATABASE = 'DynamoDB' #: Database type ('SQLite', 'DynamoDB', etc).
DATABASE_HOST = 'http://localhost:8000' #: Database Host (':memory:', 'localhost', etc).

# Define the key-value we are working with.
KEYVALUE = 'DynamoDB' #: Key-value store type ('Memory', 'DynamoDB', etc).
KEYVALUE_HOST = '' #: Key-value store host - '' if type is 'Memory'.

# DynamoDB-specific configurations - this is applied to each table.
DYNAMO_REGION = 'us-east-1' #: DynamoDB AWS region.
DYNAMO_WRITE_UNITS = 1 #: DynamoDB table write units.
DYNAMO_READ_UNITS = 1 #: DynamoDB table read units.

# Endpoint where users end up to start the signing workflow.
CLA_CONSOLE_ENDPOINT = 'http://d37jq4fjnidrq1.cloudfront.net' # ICLA QA

# Define the signing service to use.
SIGNING_SERVICE = 'DocuSign' #: The signing service to use ('DocuSign', 'HelloSign', etc)
DOCUSIGN_ROOT_URL = 'https://demo.docusign.net/restapi/v2' #: DocuSign API root URL.
DOCUSIGN_USERNAME = 'c1b49625-7634-4f17-8886-cd5b78794350' #: DocuSign username or account UUID.
DOCUSIGN_PASSWORD = '8HEt4?J92JsecYf*zbfC' #: DocuSign password.
DOCUSIGN_INTEGRATOR_KEY = '9db886ef-221d-497e-b29c-7189b372f613' #: DocuSign integrator key.

# Repository settings.
AUTO_CREATE_REPOSITORY = True #: Create repository in database automatically on webhook.

# GitHub Repository Service.
GITHUB_OAUTH_CLIENT_ID = 'client_id' #: GitHub OAuth2 client ID.
GITHUB_OAUTH_SECRET = 'secret' #: GitHub OAuth2 secret.
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
GITHUB_APP_PRIVATE_KEY_PATH = 'path-to-file'
GITHUB_APP_ID = '0000'

# KeyCloak Authentication
KEYCLOAK_ENDPOINT = 'https://<keycloak-domain>'
KEYCLOAK_CLIENT_ID = '<client-id>'
KEYCLOAK_REALM = '<realm>'
KEYCLOAK_CLIENT_SECRET = 'secret'

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
S3_ACCESS_KEY = '' #: AWS Access Key ID for the S3 service.
S3_SECRET_KEY = '' #: AWS Secret Access Key for the S3 service.
S3_BUCKET = '' #: AWS S3 bucket used to store CLA files - must be unique.

# PDF Generation.
PDF_SERVICE = 'DocRaptor'

# DocRaptor Configuration.
DOCRAPTOR_API_KEY = '' #: DocRaptopr API Key.
DOCRAPTOR_DEBUG_MODE = False #: Whether or not to print debug information when generating PDF docs.
DOCRAPTOR_TEST_MODE = False #: Test mode creates documents with watermarks at no charge.
DOCRAPTOR_JAVASCRIPT = True #: Whether javascript processing is enabled for PDF generation.
