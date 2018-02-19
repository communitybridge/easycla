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
        self.user_id = data.get('sub', None)
        self.name = data.get('name', None)
        self.session_state = data.get('session_state', None)
        self.resource_access = data.get('resource_access', None)
        self.preferred_username = data.get('preferred_username', None)
        self.given_name = data.get('given_name', None)
        self.family_name = data.get('family_name', None)
        self.email = data.get('email', None)
        self.roles = data.get('realm_access', {}).get('roles', [])
