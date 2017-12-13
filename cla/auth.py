"""
auth.py contains all necessary objects and functions to perform authentication.
"""

import hug
from jose.exceptions import ExpiredSignatureError, JWSSignatureError, JWSAlgorithmError, \
                            JWTClaimsError, JWTSignatureError, JWTError
from keycloak import KeycloakOpenID
import cla

kc = KeycloakOpenID()
kc.init(cla.conf['KEYCLOAK_ENDPOINT'],
        cla.conf['KEYCLOAK_CLIENT_ID'],
        cla.conf['KEYCLOAK_REALM'],
        cla.conf['KEYCLOAK_CLIENT_SECRET'])
kc_keys = kc.certs()

def token_verify(token):
    if token.startswith('Bearer '):
        token = token[7:]
        try:
            return kc.decode_token(token, kc_keys)
        except (ExpiredSignatureError, JWSSignatureError, JWSAlgorithmError, JWTClaimsError, \
                JWTSignatureError, JWTError) as e:
            cla.log.warning('Invalid token (%s): %s' %(e, token))
    return False

token_authentication = hug.authentication.token(token_verify)
