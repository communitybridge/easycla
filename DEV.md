# CLA Local Development

Copyright The Linux Foundation and each contributor to CommunityBridge.

SPDX-License-Identifier: CC-BY-4.0

## Backend

### cla-backend

#### Dependencies

Python 3.6.5
virtualenv
Serverless 1.32.0
nodejs (npm)

#### Setup

Create a virtual environment and install dependencies

```bash
# project root
yarn install
cd cla-backend
make setup
```

#### Local Installation

The local development environment makes use of three plugins for the Serverless
framework that emulate the AWS services locally. `serverless-wsgi` emulates AWS
λ and API Gateway, `serverless-dynamodb-local` emulates DynamoDb and
`serverless-s3-local` emulates S3. These plugins leverage Serverless in the same
way as our deployment infrastructure.

#### Plugins

* [`serverless-wsgi`](https://www.npmjs.com/package/serverless-wsgi)
* [`serverless-dynamodb-local`](https://www.npmjs.com/package/serverless-dynamodb-local)
* [`serverless-s3-local`](https://www.npmjs.com/package/serverless-s3-local)

#### Run Locally

1. Make sure Python3.6, python3-pip and `virtualenv` are installed on build
   machine. It appears Hug isn't fully compatible with Python3.7. [More info
   here](https://github.com/timothycrosley/hug/issues/631)
2. Run `export AWS_PROFILE=lf_dev` to set your AWS profile.
3. (optional) If you wish to override backend env vars, configure them in
   `env.json` using the env var names specified in `serverless.yml`.
4. Run `make setup` (This only needs to be run once.)
5. Run `make run_dynamo`
6. Run `make run_s3`
7. Run `make run_lambda`
8. Run `make add_project_manager_permission username=<username>
   project_sfdc_id=<sfdc_id> bearer_token=<token> base_url=http://localhost:5000`
   to add configure your user for the Project Management Console.

### cla-backend-go

### Dependencies

* Golang 1.11
* [dep](https://github.com/golang/dep)

#### Build And Run

Create a config file name `env.json` that contains the config entries specified
in the [config struct](/cla-backend-go/config/config.go).

Run the following commands:

```bash
cd cla-backend-go
make setup # Only needed on initial setup
make swagger
make build
./cla server --config env.json
```

## Frontend

#### Requirements

* Ionic ^3.20.0
* @ionic/app-scripts ^3.2.0
* Node ^8.11.1 at least 6.0 required.

### Install Top Level Dependencies

```bash
# Top level - project root
yarn install

# Project Management Console
cd cla-frontend-project-console
make setup

# Corporate Console
cd cla-frontend-corporate-console
make setup

# Contributor Console
cd cla-frontend-contributor-console
make setup

# or simply...
for i in cla-frontend-project-console cla-frontend-corporate-console cla-frontend-contributor-console; do
    pushd ${i} && make setup; popd;
done
```

### Config variables provision

Due to Ionic custom webpack limitations, in order to provision config variables,
we leverage `pre-build` process to trigger variable provision scripts before
executing build. For each CLA console, at `./src/package.json`, to find which
script is executed at `pre-build` process for each stage. All provision scripts
can be found under `./src/config/scripts`.

There are two ways to provision config values. The first way is to use a script to
download the values from a running environment (such as `DEV`, `STAGING`, or
`PROD`, for example) or build the file manually. Overall, both approaches should 
result in a config file named `./src/config/cla-env-config.json`. This file should
NOT be version controlled.

Using the script approach, we can fetch the configuration values from the `DEV`
environment for the project management console by doing:

```bash
# First, make sure your have your AWS profile setup and established
# Note: AWS MFA/assumed role from the command line doesn't seem to work with the node.js lib
export AWS_PROFILE=<your_aws_profile>
cd cla-frontend-project-console/src
# The following target is defined in the `cla-frontend-project-console/package.json` file
yarn prebuild:dev
```

This will generate the `cla-frontend-project-console/src/config/cla-env-config.json`
file for running in the `DEV` environment.

If the build is against qa/staging/prod stages (`qa`, `staging` and `prod`), the AWS
[parameter
store](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html)
should be used to provision these values by using a script. Go to AWS parameter store and make
sure provision property values with name convention `{variable name}-{stage
name}`. The values for each console needs are listed as following:

* CLA management console - cla-frontend-project-console

```json
{
  "auth0-clientId":"auth0_client_id",
  "auth0-domain":"linuxfoundation.auth0.com",
  "cla-api-url":"http://cla_api_endpoint",
  "cinco-api-url":"https://cinco_api_endpoint",
  "analytics-api-url":"https://analytics_api_endpoint",
  "gh-app-public-link":"https://github_app_public_link"
}
```

* CLA corporate console - cla-frontend-corporate-console

```json
{
  "auth0-clientId":"auth0_client_id",
  "auth0-domain":"linuxfoundation.auth0.com",
  "cla-api-url":"http://cla_api_endpoint",
  "cinco-api-url":"https://cinco_api_endpoint"
}
```

* CLA contributor console - cla-frontend-contributor-console

```json
{
  "auth0-clientId":"auth0_client_id",
  "auth0-domain":"linuxfoundation.auth0.com",
  "cla-api-url":"http://cla_api_endpoint",
  "corp-console-link":"https://corporate_console_link"
}
```

If the build is supposed to be a dev build hosted on a developer local machine,
directly create or edit (if it already exists) the `cla-env-config.json` configuration
file and adjust the config variables as needed.

To run the UI in local mode, run:

> Note: typically auth0 needs to be configured to support the localhost:8100
callback - many testing environments only support this in the DEV environment.

```bash
# For dev
yarn serve:dev

# For staging
yarn serve:staging
```
