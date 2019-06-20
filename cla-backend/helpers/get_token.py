# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

import sys
sys.path.append('../')

from keycloak import KeycloakOpenID
import cla

# kc = KeycloakOpenID(cla.conf['KEYCLOAK_ENDPOINT'],
#                     cla.conf['KEYCLOAK_CLIENT_ID'],
#                     cla.conf['KEYCLOAK_REALM'],
#                     cla.conf['KEYCLOAK_CLIENT_SECRET'])
# certs = kc.certs()
# token = kc.token('password', '***REMOVED***', '***REMOVED***') # Password is same as username for sandbox.
# print(token)
# print(kc.decode_token(token['access_token'], certs))
# token = kc.token('client_credentials')
# print(token)
# print(kc.decode_token(token['access_token'], certs))
