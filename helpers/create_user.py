"""
Convenience script to create a user and agreement in the CLA system.
"""
import sys
sys.path.append('../')

import uuid
import cla
from cla.utils import get_user_instance

# Test User.
cla.log.info('Creating new user')
user = get_user_instance()
user.set_user_id(str(uuid.uuid4()))
user.set_user_name('***REMOVED*** ***REMOVED***')
user.set_user_email('***REMOVED***@linuxfoundation.org')
user.set_user_github_id(123)
user.save()
