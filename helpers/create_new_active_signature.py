"""
Convenience script to create a new user signature request (simulate a user clicking on the sign icon in GitHub).
"""
import sys
sys.path.append('../')

import cla
from cla.utils import get_user_instance

# The user.
user = get_user_instance()
user = user.get_user_by_github_id(123)

repository_id = '96820382'
pull_request = '4'
cla.log.info('Creating new active signature for repository %s and PR %s', repository_id, pull_request)
# Store data on signature.
store = cla.utils.get_key_value_store_service()
key = 'active_signature:' + str(user.get_user_id())
store.set(key, repository_id + '|' + pull_request)
