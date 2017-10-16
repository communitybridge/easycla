"""
Convenience script to create a new user signature request (simulate a user clicking on the sign icon in GitHub).
"""
import sys
sys.path.append('../')

import uuid
import cla
from cla.utils import get_user_instance

# The user.
user = get_user_instance()
user = user.get_user_by_github_id(123)

# Store data on signature.
store = cla.utils.get_key_value_store_service()
key = 'active_signature:' + str(user.get_user_id())
# For now use PR 4 from the CLA-Test repository.
value = '96820382|4'
store.set(key, value)
