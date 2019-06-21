# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Convenience script to create a user and agreement in the CLA system.
"""
import sys
sys.path.append('../')

import uuid
import cla
from cla.utils import get_user_instance

# Test User.
cla.log.info('Creating first user')
user1 = get_user_instance()
user1.set_user_id(str(uuid.uuid4()))
user1.set_user_name('First User')
user1.set_user_email('firstuser@domain.org')
user1.set_user_email('foobarski@linuxfoundation.org')
user1.set_user_github_id(123)
user1.save()
cla.log.info('Creating second user')
user2 = get_user_instance()
user2.set_user_id(str(uuid.uuid4()))
user2.set_user_name('Second User')
user2.set_user_email('seconduser@listed.org')
user2.set_user_github_id(234)
user2.save()
cla.log.info('Creating third user')
user3 = get_user_instance()
user3.set_user_id(str(uuid.uuid4()))
user3.set_user_name('Third User')
user3.set_user_email('thirduser@listed.org')
user3.set_user_github_id(345)
user3.save()
