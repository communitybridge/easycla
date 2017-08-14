"""
Convenience script to create a user and agreement in the CLA system.
"""
import uuid
import sys
sys.path.append('../')

import cla
from cla.utils import get_user_instance, get_agreement_instance

test_user_name = '***REMOVED*** ***REMOVED***'
test_user_email = '***REMOVED***@linuxfoundation.org'

# Test User.
cla.log.info('Creating new user: %s <%s>', test_user_name, test_user_email)
user = get_user_instance()
user.set_user_id('023643e7-1964-4fea-a877-fce7e4bc8eb5') # Matches the send_document.py user.
user.set_user_name(test_user_name)
user.set_user_email(test_user_email)
user.save()
# Test Agreement.
cla.log.info('Creating CLA agreement for new user')
agreement = get_agreement_instance()
agreement.set_agreement_id(str(uuid.uuid4()))
agreement.set_agreement_project_id('test-project')
agreement.set_agreement_signed(True)
agreement.set_agreement_approved(True)
agreement.set_agreement_type('cla')
agreement.set_agreement_reference_id(user.get_user_id())
agreement.set_agreement_reference_type('user')
agreement.save()
cla.log.info('Done')
