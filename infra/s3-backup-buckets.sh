#!/usr/bin/env bash
# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

set -o nounset -o pipefail
declare -r SCRIPT_DIR="$( cd "$( dirname "$0" )" && pwd )"

#------------------------------------------------------------------------------
# Load helper scripts
#------------------------------------------------------------------------------
if [[ -f "${SCRIPT_DIR}/colors.sh" ]]; then
  source "${SCRIPT_DIR}/colors.sh"
else
  echo "Unable to load script: ${SCRIPT_DIR}/colors.sh"
  exit 1
fi

if [[ -f "${SCRIPT_DIR}/logger.sh" ]]; then
  source "${SCRIPT_DIR}/logger.sh"
else
  echo "Unable to load script: ${SCRIPT_DIR}/logger.sh"
  exit 1
fi

if [[ -f "${SCRIPT_DIR}/utils.sh" ]]; then
  source "${SCRIPT_DIR}/utils.sh"
else
  echo "Unable to load script: ${SCRIPT_DIR}/utils.sh"
  exit 1
fi

#------------------------------------------------------------------------------
# Check command line arguments
#------------------------------------------------------------------------------
if [[ $# -eq 0 ]]; then
  echo "Missing environment parameter. Expecting one of: 'dev', 'staging', or 'prod'."
  echo "usage:   $0 [environment]"
  echo "example: $0 dev"
  echo "example: $0 staging"
  echo "example: $0 prod"
  exit 1
fi

declare -r env="${1}"
if [[ "${env}" == 'dev' || "${env}" == 'staging' || "${env}" == 'prod' ]]; then
  echo "Using environment '${env}'..."
else
  echo "Environment parameter does not match expected values. Expecting one of: 'dev', 'staging', or 'prod'."
  echo "usage:   $0 [environment]"
  echo "example: $0 dev"
  echo "example: $0 staging"
  echo "example: $0 prod"
  exit 1
fi

#------------------------------------------------------------------------------
# Common variables
#------------------------------------------------------------------------------
declare -r region="us-east-1"
declare -r profile="lfproduct-${env}"

#------------------------------------------------------------------------------
# Check if backup directory already exists - we don't want to over-write previous backup
#------------------------------------------------------------------------------
if [[ -d ./cla-project-log-${env} ]]; then
  log_warn "Backup folder ${_Y}./cla-project-log-${env}${_W} exists. Unable to perform backup."
  exit 1
fi
if [[ -d ./cla-signature-files-${env} ]]; then
  log_warn "Backup folder ${_Y}./cla-signature-files-${env}${_W} exists. Unable to perform backup."
  exit 1
fi

#------------------------------------------------------------------------------
# Perform the s3 backup
#------------------------------------------------------------------------------
log "Backing up s3 bucket: ${_Y}s3://cla-project-logo-${env}/${_W}..."
mkdir -p ./cla-project-logo-${env}
aws --profile lfproduct-${env} s3 cp s3://cla-project-logo-${env}/ ./cla-project-logo-${env} --recursive

log "Backing up s3 bucket: ${_Y}s3://cla-signature-files-${env}/${_W}..."
mkdir -p ./cla-signature-files-${env}
aws --profile lfproduct-${env} s3 cp s3://cla-signature-files-${env}/ ./cla-signature-files-${env} --recursive
