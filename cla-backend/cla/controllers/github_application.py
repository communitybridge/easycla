import time
import cla
from github import GithubIntegration, Github
from github.GithubException import UnknownObjectException
from jose import jwt
import os

class GitHubInstallation(object):

    @property
    def app_id(self):
        return os.environ['GH_APP_ID']

    @property
    def private_key(self):
        return os.environ['GH_APP_PRIVATE_SECRET']

    @property
    def repos(self):
        return self.api_object.get_installation(self.installation_id).get_repos()

    def __init__(self, installation_id):
        self.installation_id = installation_id
        integration = GithubCLAIntegration(self.app_id, self.private_key)
        auth = integration.get_access_token(self.installation_id)
        self.token = auth.token
        self.api_object = Github(self.token)

class GithubCLAIntegration(GithubIntegration):
    """Custom GithubIntegration using python-jose instead of pyjwt for token creation."""
    def create_jwt(self):
        """
        Overloaded to use python-jose instead of pyjwt.
        Couldn't get it working with pyjwt.
        """
        now = int(time.time())
        payload = {
            "iat": now,
            "exp": now + 60,
            "iss": self.integration_id
        }
        return jwt.encode(payload, self.private_key, 'RS256')
