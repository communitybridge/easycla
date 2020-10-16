# SSM Routine

A simple command line tool to fetch a value from the AWS Systems Manager Parameter Store (SSM).

## Prerequisites

- Need to have an AWS account
- Need to have your AWS credentials setup - this application loads the default AWS configuration in your environment

## Building

```bash
go build main.go
```

## Running

```bash
# For help/usage
main --help

# Fetch a key value from SSM - using the default region and stage
main --key 'cla-gh-app-webhook-secret-dev'

# With additional options
main --key 'cla-gh-app-webhook-secret-dev' --region us-east-1 --stage dev
```
