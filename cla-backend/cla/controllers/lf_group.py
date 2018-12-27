import requests 
import json 

class LFGroup():
    def __init__(self, lf_base_url, client_id, client_secret, refresh_token):
        self.lf_base_url = lf_base_url
        data = {
            'grant_type': 'refresh_token',
            'refresh_token': refresh_token,
            'scope': 'manage_groups'
        }
        oauth_url = self.lf_base_url + "oauth2/token"
        response = requests.post(oauth_url, data=data, auth=(client_id, client_secret)).json()
        self.access_token = response['access_token']
    

    # get LDAP group 
    def get_group(self, group_id):
        headers = { 'Authorization': 'bearer ' + self.access_token } 
        response = requests.get(self.lf_base_url + 'rest/auth0/og/' + str(group_id), headers=headers)
        if response.status_code == 200:
            return response.json()
        else:
            return {'error' : 'The LDAP Group does not exist for this group ID.'}

    # add user to LDAP group
    def add_user_to_group(self, group_id, username): 
        headers = {
            'Authorization': 'bearer ' + self.access_token,
            'Content-Type': 'application/json',
            'cache-control': 'no-cache',
        }
        data = { "username": username }
        response = requests.put(self.lf_base_url + 'rest/auth0/og/' + str(group_id), headers=headers, data=json.dumps(data)).json()
        return response