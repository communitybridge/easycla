import time
import cla
from github import GithubIntegration, Github
from github.GithubException import UnknownObjectException
from jose import jwt


class GitHubInstallation(object):

    @property
    def app_id(self):
        return cla.conf['GITHUB_APP_ID']

    @property
    def private_key(self):
        return open(cla.conf['GITHUB_APP_PRIVATE_KEY_PATH']).read()

    @property
    def repos(self):
        return self.api_object.get_installation(self.installation_id).get_repos()

    def __init__(self, installation_id=None):
        if installation_id is None:
            self.installation_id = cla.conf['GITHUB_MAIN_INSTALLATION_ID']
        else:
            self.installation_id = installation_id
        integration = GithubCLAIntegration(self.app_id, self.private_key)
        auth = integration.get_access_token(self.installation_id)
        self.token = auth.token
        self.api_object = Github(self.token)

    def namespace_exists(self, namespace):
        """
        :param namespace: The name of the Github account.
        :type namespace: string
        :return: Whether or not the Github namespace exists.
        :rtype: bool
        """
        try:
            self.api_object.get_user(namespace)
            cla.log.info('Github User/Organization %s exists', namespace)
            return True
        except UnknownObjectException:
            cla.log.info('Github User/Organization %s does not exist', namespace)
            return False

    def get_namespace(self, namespace):
        """
        :param namespace: The name of the Github account.
        :type namespace: string
        :return: Dict of information on the account/organization.
        :rtype: dict
        """
        try:
            named_user = self.api_object.get_user(namespace)
            data = {'bio': named_user.bio,
                    'company': named_user.company,
                    'email': named_user.email,
                    'created_at': named_user.created_at,
                    'location': named_user.location,
                    'login': named_user.login,
                    'type': named_user.type}
            cla.log.info('Github User/Organization %s data: %s', namespace, data)
            return data
        except UnknownObjectException:
            cla.log.info('Github User/Organization %s does not exist - could not get data', namespace)
            return {'errors': {'namespace': 'Invalid GitHub account namespace'}}

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
