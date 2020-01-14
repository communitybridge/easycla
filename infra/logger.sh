#!/usr/bin/env bash

###############################################################################
# Logger definitions
###############################################################################
debug_enabled=1

function log() {
  declare -r date_format='%Y-%m-%dT%T%z'
  echo -e "${_W}[${_N}${_Y}`date +${date_format}`${_N}${_W}][${_N}${_GREEN}INFO${_N}${_W}]${_N} ${_W}$@${_N}"
}

function log_warn() {
  declare -r date_format='%Y-%m-%dT%T%z'
  echo -e "${_W}[${_N}${_Y}`date +${date_format}`${_N}${_W}][${_N}${_RED}WARN${_N}${_W}]${_N} ${_W}$@${_N}"
}

function log_debug() {
  if [ ${debug_enabled} -ne 0 ]; then
    declare -r date_format='%Y-%m-%dT%T%z'
    echo -e "${_W}[${_N}${_Y}`date +${date_format}`${_N}${_W}][${_N}${_BLUE}DEBUG${_N}${_W}]${_N} ${_W}$@${_N}"
  fi
}

# Export to allow other programs to access it
export -f log
export -f log_warn
export -f log_debug
