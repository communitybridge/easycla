# Infrastructure DevOps

This folder contains a number of infrastructure support routines for
deploying and managing EasyCLA with various AWS environments.

By default, we established a Pulumi account with a Community Bridge
organization.  This is where the state files are backed up/stored.

This directory contains a set of deployment routines which were migrated away
from the normal CI/CD serverless deployment via CircleCI due to complicated
update refresh process involving DynamoDB schema updates. In our case we
needed to update/add additional DynamoDB columns and indices over time. By
rule, you can't update multiple Global Secondary Indexes at the same time. By
extracting the DynamoDB table management and handling it separately via
Pulumi, we were able to iteratively update the the table definitions from the
command line to slowly converge on the desired state (rather than via CI/CD
which involves tagging a release for STAGING and PROD, then CI fails,
re-tag/re-run/re-tag, CI partially works, then fails, etc. etc.).

We chose Pulumi since it was simple, easy, and written in a language that
most of the team members understood (rather than learning a new DSL).

## Prerequisites

- AWS Account Details and Credentials
- AWS [Command Line Tools](https://aws.amazon.com/cli/)
- [Pulumi](https://www.pulumi.com/) installed, see the
  [getting started guide](https://www.pulumi.com/docs/get-started/).
- Access to the Pulumi Community Bridge organization (where the state files
  are stored)

## Setup

To setup Pulumi, run:

```bash
# Install the libraries
yarn install

# Log into the Pulumi servier
pulumi login
```

## Scripts

A few command line helper scripts are provided for assisting with S3 and
DyanamoDB backup and restore functions. In general, run the scripts without any
arguments to see the script usage.

| Script                    | Description |
|:--------------------------|:------------|
| backup-dynamodb-tables.sh | Backups up the DyanamoDB tables (creates a snapshot) |
| backup-s3-buckets.sh      | Backup S3 bucket contents to the local folder |
| restore-s3-buckets.sh     | Restores the S3 bucket contents to the S3 bucket |
| update-s3-permissions.sh  | Updates the S3 bucket permissions - some items are exptected to be public-read |
| logger.sh                 | Helper utility script for output log date/time stamp |
| colors.sh                 | Helper utility script for colorized output |
| utils.sh                  | Helper utility script with a couple convenience routines |


## Pulumi Stack List

We have a stack for each environment. You can see this by running
`pulumi stack ls` once you have logged into the Pulumi website via
`pulumi login`.

```bash
pulumi stack ls
NAME                     LAST UPDATE  RESOURCE COUNT  URL
communitybridge/dev      2 hours ago  18              https://app.pulumi.com/communitybridge/easycla/dev
communitybridge/prod*    in progress  18              https://app.pulumi.com/communitybridge/easycla/prod
communitybridge/staging  2 hours ago  18              https://app.pulumi.com/communitybridge/easycla/staging
```

## Importing from AWS

If you have an infrastructure already deployed and managed by other means
(manual, terraform, cloud formation, serverless, or other pulumi setup), you
can simply import the existing resources into your owen Pulumi Stack
configuration. This is common when you switch from one management
tool/deployment to another.

In order to "build up your stack" from existing resources, we recommend the
following approach - which is what we did to manage the existing DynamoDB
tables and indices from an existing deployment:

1. Comment out all the resource creation in the `index.ts` file.
1. Selectively enable each resource one at a time with the `import` clause
   flag enabled
1. If there is an error importing, review the "details" which will show you
   the delta and adjust the settings to match the previously provisioned
   resource. They need to match to get a successful import.
1. Once the resource is imported, disable the import flag and adjust the
   resource attributes to match your desired state.  Repeat this step
   until the resource is exactly like you want it.  Sometime this will be
   required multiple times for DynamoDB tables/indices since AWS only
   allows one Global Secondary Index update to run at a time.
1. Repeat for the other resources one at a time. Once a resource is
   loaded into the stack you shouldn't need to touch it again unless
   you need to make subsequent changes.

## Pulumi Deploy

```bash
pulumi up

# with sync/refresh
pulumi up -r
```

## Pulumi Destroy / Remove

Make sure you back up your data before doing this!!

```bash
pulumi destroy
```

## Backup Files from S3 to Local

Occasionally, there is a need to backup and restore data from a S3 bucket.
This section describes the backup procedure going from S3 to a local folder.

The AWS CLI makes working with files in S3 very easy. However, the file
globbing available on most Unix/Linux systems is not quite as easy to use
with the AWS CLI. S3 doesn’t have folders, but it does use the concept of
folders by using the "/" character in S3 object keys as a folder delimiter.

To copy all objects in an S3 bucket to your local machine simply use the aws
s3 cp command with the `--recursive` option.

For example `aws s3 cp s3://my-s3-bucket/ ./ --recursive` will copy all
files from the "my-s3-bucket" bucket to the current working directory on
your local machine. If there are folders represented in the object keys (keys
containing "/" characters), they will be downloaded as separate directories
in the target location.

1. First, establish your AWS credentials for the environment
1. Then run the copy command with the `--recursive` flag

Example:

```bash
mkdir ./cla-project-logo-staging
aws s3 cp s3://cla-project-logo-staging/ ./cla-project-logo-staging --recursive
mkdir ./cla-signature-files-staging
aws s3 cp s3://cla-signature-files-staging/ ./cla-signature-files-staging --recursive
```

Optionally, we wrote a convenience script:

```bash
# Usage:
backup-s3-buckets.sh <environment>

# Examples:
backup-s3-buckets.sh dev
backup-s3-buckets.sh staging
backup-s3-buckets.sh prod
```

> Note: when using the scripts, you may want to review and adjust the local
> backup folder name.

## Restore Files From Local to S3

Occasionally, there is a need to backup and restore data from a S3 bucket.
This section describes the restore procedure going from a local folder to a
S3 bucket.

1. First, establish your AWS credentials for the environment
1. Then run the copy command with the `--recursive` flag

```bash
aws s3 mb s3://cla-project-logo-staging
aws s3 cp cla-project-logo-staging/ s3://cla-project-logo-staging/ --recursive
aws s3 mb s3://cla-signature-files-staging
aws s3 cp cla-signature-files-staging/ s3://cla-signature-files-staging/ --recursive
```

[Full instructions](https://aws.amazon.com/getting-started/tutorials/backup-to-s3-cli/)

Optionally, we wrote a convenience script:

```bash
# Usage:
restore-s3-buckets.sh <environment>

# Examples:
restore-s3-buckets.sh dev
restore-s3-buckets.sh staging
restore-s3-buckets.sh prod
```

> Note: when using the scripts, you may want to review and adjust the local
> backup folder name.
