# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

import sys
sys.path.append('../')

import cla
from cla.utils import create_database, delete_database
delete_database()
create_database()
