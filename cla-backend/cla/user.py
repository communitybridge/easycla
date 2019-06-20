# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
user.py contains the user class and hug directive.
"""

from hug.directives import _built_in_directive
import cla
from jose import jwt

@_built_in_directive
def cla_user(default=None, request=None, **kwargs):
    """Returns the current logged in CLA user"""

    headers = request.headers
    if headers is None:
        cla.log.error('Error reading headers')
        return default

    bearer_token = headers.get('Authorization') or headers.get('AUTHORIZATION') 

    if bearer_token is None:
        cla.log.error('Error reading authorization header')
        return default

    bearer_token = bearer_token.replace('Bearer ', '')
    try:
        token_params = jwt.get_unverified_claims(bearer_token)
    except Exception as e:
        cla.log.error('Error parsing Bearer token: {}'.format(e))
        return default

    if token_params is not None:
        return CLAUser(token_params)
    cla.log.error('Failed to get user information from request')
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
