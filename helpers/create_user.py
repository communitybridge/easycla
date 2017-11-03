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
user = get_user_instance()
user.set_user_id(str(uuid.uuid4()))
user.set_user_name('***REMOVED*** ***REMOVED***')
user.set_user_email('***REMOVED***@linuxfoundation.org')
user.set_user_github_id(123)
user.save()
cla.log.info('Creating second user')
user = get_user_instance()
user.set_user_id(str(uuid.uuid4()))
user.set_user_name('Whitelisted User')
user.set_user_email('white@listed.org')
user.set_user_github_id(234)
user.save()
cla.log.info('Creating third user')
user = get_user_instance()
user.set_user_id(str(uuid.uuid4()))
user.set_user_name('Third User')
user.set_user_email('white@listed.org')
user.set_user_github_id(345)
user.save()
