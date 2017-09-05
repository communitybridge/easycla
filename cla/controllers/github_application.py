import cla
from github import GithubIntegration, Github


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

    def __init__(self, installation_id):
        integration = GithubIntegration(self.app_id, self.private_key)
        auth = integration.get_access_token(installation_id)
        self.api_object = Github(auth.token)
        self.installation_id = installation_id
