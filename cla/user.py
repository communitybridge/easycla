"""
user.py contains the user class and hug directive.
"""

from hug.directives import _built_in_directive
import cla

@_built_in_directive
def cla_user(default=None, request=None, **kwargs):
    """Returns the current logged in CLA user"""
    user = request.context.get('user', None)
    if user is not None:
        return CLAUser(user)
    return default

class CLAUser(object):
    """
    Data received from Keycloak:
    {'jti': '***REMOVED***', 'exp': 1513195124, 'nbf': 0, 'iat': 1513193324, 'iss': 'http://localhost:53235/auth/realms/LinuxFoundation', 'aud': 'cla', 'sub': '***REMOVED***', 'typ': 'Bearer', 'azp': 'cla', 'auth_time': 0, 'session_state': '***REMOVED***', 'acr': '1', 'allowed-origins': ['*'], 'realm_access': {'roles': ['engineering-team']}, 'resource_access': {}, 'name': '***REMOVED*** ***REMOVED***', 'preferred_username': '***REMOVED***', 'given_name': '***REMOVED***', 'family_name': '***REMOVED***', 'email': '***REMOVED***@linuxfoundation.org'}
    """
    def __init__(self, data):
        self.data = data
        self.name = data['name']
        self.session_state = data['session_state']
        self.resource_access = data['resource_access']
        self.preferred_username = data['preferred_username']
        self.given_name = data['given_name']
        self.family_name = data['family_name']
        self.email = data['email']
