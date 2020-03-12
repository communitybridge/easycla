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
declare -r profile="lfproduct-${env}"
declare -r logo_bucket_name="cla-project-logo-${env}"
declare -r signature_bucket_name="cla-signature-files-${env}"

#------------------------------------------------------------------------------
# Perform the meta-data update
#------------------------------------------------------------------------------
log "Setting ${_Y}'image/png'${_W} content type for PNG objects in s3 bucket: ${_Y}s3://${logo_bucket_name}/${_W}..."
aws --profile "${profile}" \
  s3 cp \
  "s3://${logo_bucket_name}/" \
  "s3://${logo_bucket_name}/" \
  --exclude '*' \
  --include '*.png' \
  --no-guess-mime-type \
  --content-type="image/png" \
  --metadata-directive="REPLACE" \
  --recursive

log "Setting ${_Y}'application/pdf'${_W} content type for PDF objects in s3 bucket: ${_Y}s3://${signature_bucket_name}/${_W}..."
aws --profile "${profile}" \
  s3 cp \
  "s3://${signature_bucket_name}/" \
  "s3://${signature_bucket_name}/" \
  --exclude '*' \
  --include '*.pdf' \
  --no-guess-mime-type \
  --content-type="application/pdf" \
  --metadata-directive="REPLACE" \
  --recursive
