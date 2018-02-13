import sys
sys.path.append('../')

from keycloak import KeycloakOpenID
import cla

kc = KeycloakOpenID(cla.conf['KEYCLOAK_ENDPOINT'],
                    cla.conf['KEYCLOAK_CLIENT_ID'],
                    cla.conf['KEYCLOAK_REALM'],
                    cla.conf['KEYCLOAK_CLIENT_SECRET'])
certs = kc.certs()
token = kc.token('password', '***REMOVED***', '<password-here>') # Password is same as username for sandbox.
print(token)
