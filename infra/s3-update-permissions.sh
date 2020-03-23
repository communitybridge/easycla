#!/usr/bin/env bash

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
declare -r region="us-east-1"
declare -r profile="lfproduct-${env}"
# Enable/disable the dry-run flag to skip or process the s3 bucket from
# updating - helpful for testing
declare -r dry_run_logo="false"
declare -r dry_run_signatures="false"

#------------------------------------------------------------------------------
# Check and update S3 Project Logo bucket permissions
#------------------------------------------------------------------------------
function update_project_logos() {
  log "Querying s3 bucket: ${_Y}s3://cla-project-logo-${env}/${_W} for items..."
  objects=$(aws --profile lfproduct-${env} s3api list-objects --bucket cla-project-logo-${env} --query 'Contents[].{Key: Key}' | jq -r '.[]| .Key')
  # Convert the newline delimited output to a bash array so we can loop over it
  readarray -t array <<<"${objects}"
  updated_logo_permissions_count=0

  if [[ "${dry_run_logo}" == "true" ]]; then
    log "Skipping setting cla-project-logo-${env} bucket permissions - in ${_Y}dry-run${_W} mode."
  else
    for object in "${array[@]}"; do
      log "Setting public-read permission on object: ${_Y}${object}${_W}${_W} - ${_Y}$((updated_logo_permissions_count + 1))${_W} of ${_Y}${#array[@]}${_W}"
      aws --profile "lfproduct-${env}" s3api put-object-acl --bucket "cla-project-logo-${env}" --key "${object}" --acl public-read
      updated_logo_permissions_count=$((updated_logo_permissions_count + 1))
    done
  fi
  log "Updated ${_Y}${updated_logo_permissions_count}${_W} project logo objects."
}

#------------------------------------------------------------------------------
# Check and update S3 Project Template bucket permissions
#------------------------------------------------------------------------------
function update_project_templates() {
  log "Querying s3 bucket: ${_Y}s3://cla-signature-files-${env}/${_W} for project templates..."
  #objects=$(aws --profile "lfproduct-${env}" s3api list-objects --bucket "cla-signature-files-${env}" --query 'Contents[].{Key: Key}' | jq -r '.[]| .Key')
  objects=$(aws --profile "lfproduct-${env}" s3api list-objects --bucket "cla-signature-files-${env}" --query "Contents[?contains(Key,'template')].Key" | jq -r '.[]|.')

  # Convert the newline delimited output to a bash array so we can loop over it
  readarray -t array <<<"${objects}"
  updated_signature_permissions_count=0

  if [[ "${dry_run_signatures}" == "true" ]]; then
    log "Skipping setting cla-signature-files-${env} bucket permissions - in ${_Y}dry-run${_W} mode."
  else
    log "Processing ${_Y}${#array[@]}${_W} objects..."
    for object in "${array[@]}"; do
      #log_debug "Testing ${object}"
      if [[ "${object}" == *"template"* ]]; then
        log "Setting public-read permission on object: ${_Y}${object}${_W} - ${_Y}$((updated_signature_permissions_count + 1))${_W} of ${_Y}${#array[@]}${_W}"
        aws --profile "lfproduct-${env}" s3api put-object-acl --bucket "cla-signature-files-${env}" --key "${object}" --acl public-read
        updated_signature_permissions_count=$((updated_signature_permissions_count + 1))
      fi
    done
  fi
  log "Updated ${_Y}${updated_signature_permissions_count}${_W} project template objects."
}

#------------------------------------------------------------------------------
# Main function
#------------------------------------------------------------------------------
function main() {
  update_project_logos
  update_project_templates
}

#------------------------------------------------------------------------------
# Application entry point - call the main function
#------------------------------------------------------------------------------
main
