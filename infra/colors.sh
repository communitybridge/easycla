#!/usr/bin/env bash

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

# Colors consist of one or more of the following numeric codes:
# Attribute codes:
#   00=none 01=bold 04=underscore 05=blink 07=reverse 08=concealed
# Text color codes:
#   30=black 31=red 32=green 33=yellow 34=blue 35=magenta 36=cyan 37=white
# Background color codes:
#   40=black 41=red 42=green 43=yellow 44=blue 45=magenta 46=cyan 47=white

# Reference: https://wiki.archlinux.org/index.php/Color_Bash_Prompt
# Discussion on the escape codes: http://superuser.com/questions/33914/why-doesnt-echo-support-e-escape-when-using-the-e-argument-in-macosx
# Reset
export Color_Off='\033[0m'       # Text Reset
export _NC='\033[0m'             # Text Reset
export _N='\033[0m'              # Text Reset

# Regular Colors
export Black='\033[0;30m'        # Black
export _BLACK='\033[0;30m'       # Black
export _BK='\033[0;30m'          # Black
export Red='\033[0;31m'          # Red
export _RED='\033[0;31m'         # Red
export _R='\033[0;31m'           # Red
export Green='\033[0;32m'        # Green
export _GREEN='\033[0;32m'       # Green
export _G='\033[0;32m'           # Green
export Yellow='\033[0;33m'       # Yellow
export _YELLOW='\033[0;33m'      # Yellow
export _Y='\033[0;33m'           # Yellow
export Blue='\033[0;34m'         # Blue
export _BLUE='\033[0;34m'        # Blue
export _BU='\033[0;34m'          # Blue
export Purple='\033[0;35m'       # Purple
export _MAGENTA='\033[0;35m'     # Purple
export _M='\033[0;35m'           # Purple
export Cyan='\033[0;36m'         # Cyan
export _CYAN='\033[0;36m'        # Cyan
export _C='\033[0;36m'           # Cyan
export White='\033[0;37m'        # White
export _WHITE='\033[0;37m'       # White
export _W='\033[0;37m'           # White

# Bold
export BBlack='\033[1;30m'       # Black
export BRed='\033[1;31m'         # Red
export BGreen='\033[1;32m'       # Green
export BYellow='\033[1;33m'      # Yellow
export BBlue='\033[1;34m'        # Blue
export BPurple='\033[1;35m'      # Purple
export BCyan='\033[1;36m'        # Cyan
export BWhite='\033[1;37m'       # White

# Underline
export UBlack='\033[4;30m'       # Black
export URed='\033[4;31m'         # Red
export UGreen='\033[4;32m'       # Green
export UYellow='\033[4;33m'      # Yellow
export UBlue='\033[4;34m'        # Blue
export UPurple='\033[4;35m'      # Purple
export UCyan='\033[4;36m'        # Cyan
export UWhite='\033[4;37m'       # White

# Background
export On_Black='\033[40m'       # Black
export On_Red='\033[41m'         # Red
export On_Green='\033[42m'       # Green
export On_Yellow='\033[43m'      # Yellow
export On_Blue='\033[44m'        # Blue
export On_Purple='\033[45m'      # Purple
export On_Cyan='\033[46m'        # Cyan
export On_White='\033[47m'       # White

# High Intensity
export IBlack='\033[0;90m'       # Black
export IRed='\033[0;91m'         # Red
export IGreen='\033[0;92m'       # Green
export IYellow='\033[0;93m'      # Yellow
export IBlue='\033[0;94m'        # Blue
export IPurple='\033[0;95m'      # Purple
export ICyan='\033[0;96m'        # Cyan
export IWhite='\033[0;97m'       # White

# Bold High Intensity
export BIBlack='\033[1;90m'      # Black
export BIRed='\033[1;91m'        # Red
export BIGreen='\033[1;92m'      # Green
export BIYellow='\033[1;93m'     # Yellow
export BIBlue='\033[1;94m'       # Blue
export BIPurple='\033[1;95m'     # Purple
export BICyan='\033[1;96m'       # Cyan
export BIWhite='\033[1;97m'      # White

# High Intensity backgrounds
export On_IBlack='\033[0;100m'   # Black
export On_IRed='\033[0;101m'     # Red
export On_IGreen='\033[0;102m'   # Green
export On_IYellow='\033[0;103m'  # Yellow
export On_IBlue='\033[0;104m'    # Blue
export On_IPurple='\033[0;105m'  # Purple
export On_ICyan='\033[0;106m'    # Cyan
export On_IWhite='\033[0;107m'   # White

# Checkbox and X mark
export CHK="\xE2\x9C\x94"       # checkbox character (used with printf)
#export CHK="\u2714" # checkbox character (suitable with echo -e)
export EKS="\xE2\x9D\x8C"         # X character (used with printf)
#export EKS="\u274c"   # X character (suitable with echo -e)
