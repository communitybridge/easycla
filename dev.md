# DEV

## CLA Local Development

Copyright The Linux Foundation and each contributor to CommunityBridge.

SPDX-License-Identifier: CC-BY-4.0

## Prerequisites

- Go 1.13
- Python
- Node 8/12

## Node Setup

The frontend UI still requires an older version of nodejs - version 8. This is
due to compatibility issues with the Angular and Ionic toolchain. Use node
version 8.x for the frontend package installation and build/deploy commands.
For all other folders which use node, use version 12.x+.  The serverless
library picks up a newer version of semver which requires node 10+.

In order to quickly switch between node versions, use
[Node Version Manager - nvm](https://github.com/nvm-sh/nvm). The CircleCI
build configuration uses this approach to switch between node versions
within the build.

## Building and Running the Python Backend

These are the steps for setting up a local dev environment for the python
backend (work in progress).

### Install Python

You will need to install a local Python 3.6+ runtime - Python 3.7 will work now
that we've updated the `hug` library. 

One easy way to install a Python runtime on a Mac is via the `pipenv` tool.
Here, we'll use it to install a specific version of Python.  There are other
apporaches to installing Python. Feel free to add your preferred approach to the
information below.

```bash
brew install pipenv
```

To create a new Python installation using `pipenv`, specify the specific version
using the `--python VERSION` flag, like so:

```bash
pipenv --python 3.7
```

When given a Python version, like this, Pipenv will automatically scan your
system for a Python that matches that given version and download one if
necessary.

### Configure the Virtual Environment

You can use `pipenv` or any other tool to configure a local sandbox for testing.
We prefer to use `virtualenv`.

```bash
pip3 install virtualenv
```

Now setup the virtual env using python 3.6+:

```bash
# Check your version and grab the appropriate path
python --version
which python

# Setup the virtual env
cd cla-backend
virtualenv --python=/usr/local/Cellar/python3/3.7.4/bin/python3 .venv
```

This will create a `.venv` folder in the `cla-backend` project folder. We'll
want to load the sandbox environment by using:

```bash
source .venv/bin/activate
```

### Install the Requirements

If not already done, load the virtual environment. Then we can install the
required modules:

```bash
source .venv/bin/activate
pip3 install -r requirements.txt
```

This will install the dependencies in the `.venv` path.

### Setup Environment Variables

We'll need a few run-time environment variables to communicate with AWS.

- `AWS_REGION` - the AWS region, normally this should be set to `us-east-1`
- `AWS_ACCESS_KEY_ID` - AWS key, used to authenticate to AWS for DynamoDB and SSM
- `AWS_SECRET_ACCESS_KEY` - AWS secret key, used to authenticate to AWS for DynamoDB and SSM
- `PRODUCT_DOMAIN` - e.g. `export PRODUCT_DOMAIN=dev.lfcla.com`
- `ROOT_DOMAIN` - e.g. `export ROOT_DOMAIN=lfcla.dev.platform.linuxfoundation.org`

Optional environment variables:

- `PORT` - optional, the HTTP port when running in local mode. The default is 5000.
- `STAGE` - optional, specifies the environment stage. The default is typically `dev`.

For testing locally, you'll want to point to the dev environment or stand up a
local DynamoDB and S3 instance. Generally, we run the backend services and UI
locally and simply point to the DEV environment. The `STAGE` environment
variable controls where we point. Make sure you export/provide/setup the AWS
properties in order to connect.

## Run the Python Backend

```bash
# ok to use node v12+
cd cla-backend
yarn install
```

This will install the `serverless` tool and plugins.

```bash
# If using local DynamoDB and local S3:
node_modules/serverless/bin/serverless wsgi serve -s 'local' 

# If pointing to a DEV cluster:
node_modules/serverless/bin/serverless wsgi serve -s 'dev' 

# Alternatively, you can simply run:
# ok to use node v12+
yarn serve:dev
```

To test:

```bash
# Health Check
open http://localhost:5000/v2/health

# Get User by ID
open http://localhost:5000/v2/user/<some_uuid_from_users_table>
```

## Building and Running the Go Backend

### Install Go

Follow the [Go Getting Started](https://golang.org/doc/install) guide to install the tool.

After installation, you should have something similar to:

```bash
# After the GO installation, you should have GOROOT set
echo $GOROOT
/usr/local/opt/go/libexec

# And go should be in your path
go version
go version go1.14 darwin/amd64
```

### Setup Project Folders

To build the go stuff, the go tool is very specific on paths. If this isn't
correct, some of the build targets, such as `swagger` will fail with a
nondescript error message. Generally:

```bash
# GOPATH also should be set, depends on where your code would generally be located
echo $GOPATH
# response would be something like this:
/Users/<my-os-user-id>/projects/go

# if `GOPATH` is not set, you will need to explictly set it:
export GOPATH=$HOME/projects/go # or wherever your go projects will be located
# most sane people set this in their `~/.bashrc` or `~/.zshrc` file so they
# don't have to remember to set it again

# If the `GOPATH` variable is NOT set, run: export GOPATH=<your_go_path>
# Inside of this folder you should traditionally have `src`, `bin`, and `pkg` folders
# To recreate them, run (needed for the first-time setup only):
pushd $GOPATH && mkdir -p src bin pkg && popd

# If you have code checked out in a different path than the GOPATH (as in the
# case of CLA, probably), you will need to establish a soft link to point where
# you checked out the public easycla source - this must follow the EasyCLA go path
# structure (e.g. the directory path of the source code on disk must match the
# code import paths). To create the path structure, run:
mkdir -p $GOPATH/src/github.com/communitybridge
pushd $GOPATH/src/github.com/communitybridge
ln -s <path_to_where_you_really_have_easycla_checked_out>/easycla easycla
popd

# Confirm everything:
ls -la $GOPATH/src/github.com/communitybridge/easycla
# Should see something like:
lrwxr-xr-x  1 ddeal  staff  29 Jul  7 14:48 /Users/ddeal/projects/go/src/github.com/communitybridge/easycla@ -> /Users/ddeal/projects/easycla
```

### Configure Git

EasyCLA currently relies on two private libraries for user authorization and
common data models:

- github.com/LF-Engineering/lfx-kit (authorization library)
- github.com/LF-Engineering/lfx-models (common data models)

The near-term plan is to migrate these priavte libraries to the 
github.com/communitybridge organization as public repositories.

Until that time, each user needs to request access to this repository 
and create/update the following file: `~/.gitconfig` to include the following:

```code
[url "ssh://git@github.com/"]
  insteadOf = https://github.com/
```

This will allow the Golang `dep` tool to pull dependencies from private
repositories.

### Building Go Source

```bash
# Navigate to the go source code and build it - you must be in this path/directory
pushd $GOPATH/src/github.com/communitybridge/easycla/cla-backend-go

# First time only, you will need to install the dev tools.  Simply run:
make setup

# Once the tools are installed, you can run a clean build with all
# the bells and whistles. This will:
# - remove the old binary
# - generate go code based on the swagger specification (includes data models and API stuff)
# - download external dependencies
# - formats the code (should be already formatted) - run this before you checkin code - the linter will catch violations
# - build the code - use `build` or `build-mac` based on your platform
# - test will run unit tests
# - lint will run lint checks (do this before checking in code to avoid CI/CD catching the violations later)

# Linux only:
make clean swagger swagger-validate deps fmt build test lint

# Mac only:
make clean swagger swagger-validate deps fmt build-mac test lint
# or use the 'all' target (mac only)
make all

# After the above, you should have the binary now:
ls -lhF cla
-rwxr-xr-x  1 ddeal  staff    36M Jul 18 10:57 cla*

# Type is based on your OS
file cla
cla: Mach-O 64-bit executable x86_64
```

### Setup the Environment

To run the Go backend, it requires a few environment variables to be set:

- `AWS_REGION` - the AWS region, normally this should be set to `us-east-1`
- `AWS_ACCESS_KEY_ID` - AWS key, used to authenticate to AWS for DynamoDB and SSM
- `AWS_SECRET_ACCESS_KEY` - AWS secret key, used to authenticate to AWS for DynamoDB and SSM

Optional environment settings:

- `PORT` - optional, the HTTP port when running in local mode. The default is 8080.
- `STAGE` - optional, specifies the environment stage. The default is `dev`.
- `GH_ORG_VALIDATION` - set to `false` to test locally which will by-pass the GH auth checks and
   allow local functional tests (e.g. with cURL or Postman) - default is enabled/true

### Running

First build and setup the environment.  Then simply run it:

```bash
./cla
```

You should see the typical diagnostic details on startup indicating that it
started without errors. To confirm, connect to the health service using a
browser, curl or PostMan:

```bash
open http://localhost:8080/v3/ops/health
```

## Testing UI Locally

If testing in local mode, set the `USE_LOCAL_SERVICES=true` environment variable
and review the localhost values in the `cla.service.ts` implementation for each
console. Update the ports as necessary.

For example, to test the corporate console:

```bash
cd cla-frontend-corporate-console
# use node 8
yarn install-frontend

# move to the source folder
cd src

# Ensure the AWS environment is setup
export AWS_PROFILE=cla-dev-profile
export AWS_REGION=us-east-1

# To use local services, set the following value.
# You will need to have the go and python backend up and running locally,
# otherwise it will use the DEV, STAGING, or PROD environment services based on
# the STAGE environment variable and the specific yarn target that you invoke.
export USE_LOCAL_SERVICES=true

# Run locally (code changes auto-deploy) - dev environment settings
# use node 8
yarn serve:dev
```
