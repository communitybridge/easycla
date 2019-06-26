# CLA Deployment

Copyright The Linux Foundation and each contributor to CommunityBridge.

SPDX-License-Identifier: CC-BY-4.0

The usual deployment flow begins with dev, which we deploy to when the feature or bug is ready for testing in an AWS environment. Once the changes are tested and appear complete, we deploy to staging for QA and acceptance. After the changes have been accepted, we deploy to prod.

## Staging & Prod

Staging and prod builds are triggered by creating a "tag" in GitHub. The tag must begin with a lowercase "v", and follows semantic versioning. e.g. "v1.0.0"

Once the tag is created, CircleCI will build the projects. Human approval is required to deploy to staging and prod. In the CircleCI console, select the workflow triggered by the tag, then click "approve_staging" or "approve_prod" to deploy to the appropriate environment.

## Dev

### Backend

`cla-backend` and `cla-backend-go` are both deployed using the `serverless.yml` located at `cla-backend/serverless.yml`. The golang backend executable must be built and copied to the `cla-backend` directory.

#### Requirements

* Serverless
* AWS CLI
* nodejs (npm)
* yarn

#### Deploy

```bash
# Build golang executable
cd cla-backend-go
make build_aws_lambda
# Copy golang executable to cla-backend
cp build_aws_lambda ../cla-backend
# Deploy
cd ../cla-backend
make deploy
```

#### Initial Deployment Steps

The following steps must be performed the first time a frontend application is deployed to a new AWS environment

##### Validate the SSL Certificate in Certificate Manager

1. Log in to the AWS console
2. Select `Services` > `Certificate Manager`
3. Expand the row `api.lfcla.dev.platform.linuxfoundation.org`
4. For each domain, click `Create Record in Route53`

##### Create subdomain records in Route53

1. Log in to the AWS console
2. Select `Services` > `Route53` > `Hosted Zones` > `dev.lfcla.com`
3. Click `Create Record Set`
4. Enter Name: `api`, Type: `A - IPv4 Address`, Alias: `Yes`
5. Click `Alias Target` and select the Project Management Console CloudFront Distribution

### Frontend

#### Requirements

* Serverless
* AWS CLI
* nodejs (npm)
* yarn

#### Deploying Webpage

```bash
make setup
export AWS_PROFILE=lf_dev
export STAGE=dev
export ROOT_DOMAIN=lfcla.dev.platform.linuxfoundation.org
export PRODUCT_DOMAIN=dev.lfcla.com
make deploy
```

#### Initial Deployment Steps

The following steps must be performed the first time a frontend application is deployed to a new AWS environment

##### Validate the SSL Certificate in Certificate Manager

1. Log in to the AWS console
2. Select `Services` > `Certificate Manager`
3. Expand the row `[project|corporate|contributor].lfcla.dev.platform.linuxfoundation.org`
4. For each domain, click `Create Record in Route53`

##### Create subdomain records in Route53

1. Log in to the AWS console
2. Select `Services` > `Route53` > `Hosted Zones` > `dev.lfcla.com`
3. Click `Create Record Set`
4. Enter Name: `[project|corporate|contributor]`, Type: `A - IPv4 Address`, Alias: `Yes`
5. Click `Alias Target` and select the Project Management Console CloudFront Distribution
