import requests
import json
import os

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
        response = requests.post(oauth_url, data=data, auth=(self.client_id, self.client_secret)).json()

        return response['access_token']

    # get LDAP group 
    def get_group(self, group_id):
        access_token = self._get_access_token()
        headers = { 'Authorization': 'bearer ' + access_token }
        get_group_url = os.path.join(self.lf_base_url, 'rest/auth0/og/', str(group_id))
        response = requests.get(get_group_url, headers=headers)

        if response.status_code == 200:
            return response.json()
        else:
            return {'error' : 'The LDAP Group does not exist for this group ID.'}

    # add user to LDAP group
    def add_user_to_group(self, group_id, username): 
        access_token = self._get_access_token()
        headers = {
            'Authorization': 'bearer ' + access_token,
            'Content-Type': 'application/json',
            'cache-control': 'no-cache',
        }
        data = { "username": username }
        add_user_url = os.path.join(self.lf_base_url, 'rest/auth0/og/', str(group_id))
        response = requests.put(add_user_url, headers=headers, data=json.dumps(data))

        if response.status_code == 200:
            return response.json()
        else:
            return {'error' : 'failed to add a user to the ldap group.'}
