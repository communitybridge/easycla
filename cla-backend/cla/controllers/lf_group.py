# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

import requests
import json
import os
import cla

class LFGroup():
    def __init__(self, lf_base_url, client_id, client_secret, refresh_token):
        self.lf_base_url = lf_base_url
        self.client_id = client_id
        self.client_secret = client_secret
        self.refresh_token = refresh_token

    def _get_access_token(self):
        data = {
            'grant_type': 'refresh_token',
            'refresh_token': self.refresh_token,
            'scope': 'manage_groups'
        }
        oauth_url = os.path.join(self.lf_base_url, 'oauth2/token')

        try:
            response = requests.post(oauth_url, data=data, auth=(self.client_id, self.client_secret)).json()
        except requests.exceptions.RequestException as e:
            cla.log.error(e)
            return None

        return response.get('access_token')

    # get LDAP group 
    def get_group(self, group_id):
        access_token = self._get_access_token()
        if access_token is None:
            return {'error': 'Unable to retrieve access token'}

        headers = { 'Authorization': 'bearer ' + access_token }
        get_group_url = os.path.join(self.lf_base_url, 'rest/auth0/og/', str(group_id))

        try:
            response = requests.get(get_group_url, headers=headers)
        except requests.exceptions.RequestException as e:
            cla.log.error(e)
            return {'error': 'Unable to get group'}

        if response.status_code == 200:
            return response.json()
        else:
            return {'error' : 'The LDAP Group does not exist for this group ID.'}

    # add user to LDAP group
    def add_user_to_group(self, group_id, username): 
        access_token = self._get_access_token()
        if access_token is None:
            return {'error': 'Unable to retrieve access token'}

        headers = {
            'Authorization': 'bearer ' + access_token,
            'Content-Type': 'application/json',
            'cache-control': 'no-cache',
        }
        data = { "username": username }
        add_user_url = os.path.join(self.lf_base_url, 'rest/auth0/og/', str(group_id))

        try:
            response = requests.put(add_user_url, headers=headers, data=json.dumps(data))
        except requests.exceptions.RequestException as e:
            cla.log.error(e)
            return {'error': 'Unable to update group'}

        if response.status_code == 200:
            cla.log.info('LFGroup; Successfully added user %s into group %s', username, str(group_id))
            return response.json()
        else:
            cla.log.error('LFGroup; Failed adding user %s into group %s', username, str(group_id))
            return {'error' : 'failed to add a user to the ldap group.'}
