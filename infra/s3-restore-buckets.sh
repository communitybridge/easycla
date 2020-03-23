#!/usr/bin/env bash

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
# Check if the previous backup directory already exists - we can't restore if
# it doesn't exist
#------------------------------------------------------------------------------
if [[ ! -d ./cla-project-log-${env} ]]; then
  log_warn "Backup folder ${_Y}./cla-project-log-${env}${_W} DOES NOT exist. Unable to perform restore."
  exit 1
fi
if [[ ! -d ./cla-signature-files-${env} ]]; then
  log_warn "Backup folder ${_Y}./cla-signature-files-${env}${_W} DOES NOT exist. Unable to perform restore."
  exit 1
fi

#------------------------------------------------------------------------------
# Perform the s3 restore
#------------------------------------------------------------------------------
log "Restoring ${_Y}cla-project-logo-${env}${_W}..."
aws s3 mb s3://cla-project-logo-${env}
aws s3 cp cla-project-logo-${env}/ s3://cla-project-logo-${env}/ --recursive

log "Restoring ${_Y}cla-signature-files-${env}${_W}..."
aws s3 mb s3://cla-signature-files-${env}
aws s3 cp cla-signature-files-${env}/ s3://cla-signature-files-${env}/ --recursive
