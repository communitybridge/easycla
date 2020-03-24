#!/usr/bin/env bash
# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

set -o nounset -o pipefail
declare -r SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

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
declare -r profile="lfproduct-${env}"
declare -r logo_bucket_name="cla-project-logo-${env}"
declare -r signature_bucket_name="cla-signature-files-${env}"

#------------------------------------------------------------------------------
# Update the content type for the specified files identified by the filter
#------------------------------------------------------------------------------
function update_content_type() {
  if [[ $# -ne 3 ]]; then
    log_warn "Missing function parameters for ${FUNCNAME[0]} - expecting 3"
    log "Usage  : ${FUNCNAME[0]} <s3-bucket> <content-type> <filename_filter>"
    log "Example: ${FUNCNAME[0]} cla-project-logo-dev 'image/png' and '*.png'"
    log "Example: ${FUNCNAME[0]} cla-signature-files-dev 'application/pdf' and '*.pdf'"
    return
  fi

  # Grab the arguments
  bucket_name="${1}"
  content_type="${2}"
  filter="${3}"

  log "Setting ${_Y}'${content_type}'${_W} content type for ${_Y}'${filter}'${_W} objects in s3 bucket: ${_Y}s3://${bucket_name}/${_W}..."
  aws --profile "${profile}" \
    s3 cp \
    "s3://${bucket_name}/" \
    "s3://${bucket_name}/" \
    --exclude '*' \
    --include "${filter}" \
    --no-guess-mime-type \
    --content-type="${content_type}" \
    --metadata-directive="REPLACE" \
    --recursive
}

#------------------------------------------------------------------------------
# Main function
#------------------------------------------------------------------------------
function main() {
  update_content_type "${logo_bucket_name}" 'image/png' '*.png'
  update_content_type "${signature_bucket_name}" 'application/pdf' '*.pdf'
}

#------------------------------------------------------------------------------
# Application entry point - call the main function
#------------------------------------------------------------------------------
main
