#!/usr/bin/env bash
# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

###############################################################################
# What to do if we die
###############################################################################
die() {
  log_warn "$@"
  exit 1
}

###############################################################################
# Formatted date/time string
###############################################################################
get_logger_timestamp() {
  declare -r date_format='%Y-%m-%dT%T%z'
  echo -e "`date +${date_format}`"
}

###############################################################################
# Retries a command a with backoff.
#
# The retry count is given by ATTEMPTS (default 5), the
# initial backoff timeout is given by TIMEOUT in seconds
# (default 1.)
#
# Successive backoffs double the timeout.
#
# Beware of set -e killing your whole script!
#
# Reference: https://coderwall.com/p/--eiqg/exponential-backoff-in-bash
###############################################################################
function with_backoff {
  local max_attempts=${ATTEMPTS-5}
  local timeout=${TIMEOUT-1}
  local attempt=0
  local exitCode=0

  while [[ $attempt < $max_attempts ]]; do
    "$@"
    exitCode=$?

    if [[ $exitCode == 0 ]]
    then
      break
    fi

    echo "Failure! Retrying in $timeout.." 1>&2
    sleep $timeout
    attempt=$(( attempt + 1 ))
    timeout=$(( timeout * 2 ))
  done

  if [[ $exitCode != 0 ]]; then
    echo "You've failed me for the last time! ($@)" 1>&2
  fi

  return $exitCode
}
