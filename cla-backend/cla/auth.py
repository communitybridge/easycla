"""
auth.py contains all necessary objects and functions to perform authentication and authorization.
"""

import hug
import falcon
import requests
# from jose.exceptions import ExpiredSignatureError, JWSSignatureError, JWSAlgorithmError, \
#                             JWTClaimsError, JWTSignatureError, JWTError
# from keycloak import KeycloakOpenID
import cla
import functools

# kc = KeycloakOpenID()
# kc.init(cla.conf['KEYCLOAK_ENDPOINT'],
#         cla.conf['KEYCLOAK_CLIENT_ID'],
#         cla.conf['KEYCLOAK_REALM'],
#         cla.conf['KEYCLOAK_CLIENT_SECRET'])
# kc_keys = kc.certs()

# def decode_token(token):
#     """
#     Wrapper around the Keycloak decode_token function that will handle the pmc audience.
#     """
#     kc.client_id = 'pmc'
#     token = kc.decode_token(token, kc_keys)
#     kc.client_id = 'cla'
#     return token
#
# def token_verify(token):
#     if token.startswith('Bearer '):
#         token = token[7:]
#         try:
#             return decode_token(token)
#         except (ExpiredSignatureError, JWSSignatureError, JWSAlgorithmError, JWTClaimsError, \
#                 JWTSignatureError, JWTError) as e:
#             cla.log.warning('Invalid token (%s): %s' %(e, token))
#     return False
#
# def pm_verify(user, project_id):
#     """
#     Helper function to ensure the request is made by a user who is manager for the project in question.
#     """
#     project = cla.utils.get_project_instance()
#     project.load(str(project_id))
#     external_project_id = project.get_project_external_id()
#     pm_verify_external_id(user, external_project_id)
#
# def pm_verify_external_id(user, external_project_id):
#     """
#     Same as pm_verify() except it deals with external_project_id (CINCO project_id).
#     """
#     # TODO: Need to look into this - will fetching a new token every request cause issues?
#     cla_token = kc.token('client_credentials')
#     headers = {'Authorization': 'Bearer ' + cla_token['access_token']}
#     req = requests.get(cla.conf['CINCO_ENDPOINT'] + '/project/' + external_project_id + '/managers', headers=headers)
#     data = req.json()
#     if req.status_code is 200:
#         if user.user_id is not None and user.user_id in data:
#             return True
#         else:
#             raise falcon.HTTPError('403 Not Authorized', 'authorization', 'Insufficient permissions')
#     else:
#         cla.log.error('Could not fetch project managers from CINCO: %s - %s', data['code'], data['message'])
#         raise falcon.HTTPError('500 Internal Server Error', 'cinco_communication', 'Could not fetch project managers')
#
# def staff_verify(user):
#     if 'engineering-team' in user.roles or \
#        'STAFF_SUPER_ADMIN' in user.roles or \
#        'integration-api-role' in user.roles:
#         return True
#     raise falcon.HTTPError('403 Not Authorized', 'authorization', 'Insufficient permissions')
#
# def company_manager_verify(user, company_id):
#     user_obj = cla.utils.get_user_instance()
#     user_obj = user_obj.get_user_by_email(user.email)
#     comp = cla.utils.get_company_instance()
#     companies = comp.get_companies_by_manager_id(user_obj.get_user_id())
#     for company in companies:
#         if company['company_id'] == company_id:
#             return True
#     raise falcon.HTTPError('403 Not Authorized', 'authorization', 'Not company manager')
#
# token_authentication = hug.authentication.token(token_verify)


#
# Verify staff decorator
#
def staff_required(func):
    @functools.wraps(func)
    def inner(request, *args, **kwargs):
        #so here we can check user role from API-GW-Auth.
        print(request.env["API_GATEWAY_AUTHORIZER"])
        return func(*args, **kwargs)
    return inner
#
# Verify company manager decorator
#
def require_company_manager():
    pass

#
# Verify project manager decorator
#
def require_project_manager():
    pass


##

# Following code is not used
# This implemation is started to use `requires` in `@hug.get(requires=func)`.
# But stopped due to lack of access to `{company_id}` from path
# Decorator chaining also not giving access to varaibles in path

# from falcon import HTTPUnauthorized
# @hug.authentication.token
# def company_owner(request, response, *args, **kwargs):

#     user_id = cla_user(request=request)
#     if user_id in "company_acl":
#         return True
#     else:
#         return False
#         raise HTTPUnauthorized('Invalid Authentication',
#                             'Provided Token credentials were invalid')

# Example usage

# @hug.put('/endpoing/{test_value}', versions=1, requires=company_owner)
# def endpoint_with_requires(user: cla_user, test_value):
#     return {'ok': 'fine'}
