"""
Holds the GitHub repository service.
"""

import uuid
import falcon
import github
from github.GithubException import UnknownObjectException, BadCredentialsException
from requests_oauthlib import OAuth2Session
import cla
from cla.models import repository_service_interface, DoesNotExist


class GitHub(repository_service_interface.RepositoryService):
    """
    The GitHub repository service.
    """
    def __init__(self):
        self.client = None

    def initialize(self, config):
        username = config['GITHUB_USERNAME']
        token = config['GITHUB_TOKEN']
        self.client = self._get_github_client(username, token)

    def _get_github_client(self, username, token): # pylint: disable=no-self-use
        return github.Github(username, token)

    def get_repository_id(self, repo_name):
        """
        Helper method to get a GitHub repository ID based on repository name.

        :param repo_name: The name of the repository, example: 'linuxfoundation/cla'.
        :type repo_name: string
        :return: The repository ID.
        :rtype: integer
        """
        try:
            return self.client.get_repo(repo_name).id
        except github.GithubException as err:
            cla.log.error('Could not find GitHub repository (%s), ensure it exists and that '
                          'your personal access token is configured with the repo scope', repo_name)
        except Exception as err:
            cla.log.error('Unknown error while getting GitHub repository ID for repository %s: %s',
                          repo_name, str(err))

    def received_activity(self, data):
        cla.log.debug('Received GitHub activity: %s', data)
        if 'pull_request' not in data:
            cla.log.debug('Activity not related to pull request')
            return {'message': 'Not a pull request - no action performed'}
        if data['action'] == 'opened':
            cla.log.debug('Handling opened pull request')
            return self.process_opened_pull_request(data)
        elif data['action'] == 'reopened':
            cla.log.debug('Handling reopened pull request')
            return self.process_reopened_pull_request(data)
        elif data['action'] == 'closed':
            cla.log.debug('Handling closed pull request')
            return self.process_closed_pull_request(data)
        elif data['action'] == 'synchronize':
            cla.log.debug('Handling synchronized pull request')
            return self.process_synchronized_pull_request(data)

    def sign_request(self, repository_id, change_request_id, request):
        cla.log.info('Initiating GitHub sign request for repository %s', repository_id)
        session = self._get_request_session(request)
        repository = cla.utils.get_repository_instance()
        try:
            repository.load(repository_id)
        except DoesNotExist as err:
            cla.log.error('Error when initiating GitHub sign request for repository %s and ' + \
                          'change request %s: %s', repository_id, change_request_id, str(err))
            return {'errors': {'repository_id': str(err)}}
        if 'github_oauth2_token' in session:
            cla.log.info('Using existing session OAuth2 token')
            return self.sign_workflow(repository, change_request_id, request)
        else:
            cla.log.info('Initiating GitHub OAuth2 exchange')
            authorization_url, state = self.get_authorization_url_and_state(repository_id,
                                                                            change_request_id,
                                                                            ['user:email'])
            session['github_oauth2_state'] = state
            cla.log.info('GitHub OAuth2 request with state %s - sending user to %s',
                         state, authorization_url)
            raise falcon.HTTPFound(authorization_url)

    def _get_request_session(self, request): # pylint: disable=no-self-use
        """
        Mockable method used to get the current user session.
        """
        return request.context['session']

    def get_authorization_url_and_state(self, repository_id, pull_request_number, scope):
        """
        Helper method to get the GitHub OAuth2 authorization URL and state.

        This will be used to get the user email from GitHub.

        :param repository_id: The ID of the repository this request was initiated in.
        :type repository_id: int
        :param pull_request_number: The PR number this request was generated in.
        :type pull_request_number: int
        :param scope: The list of OAuth2 scopes to request from GitHub.
        :type scope: [string]
        """
        client_id = cla.conf['GITHUB_CLIENT_ID']
        redirect_uri = cla.utils.get_redirect_uri('github', repository_id, pull_request_number)
        return self._get_authorization_url_and_state(client_id,
                                                     redirect_uri,
                                                     scope,
                                                     cla.conf['GITHUB_OAUTH_AUTHORIZE_URL'])

    def _get_authorization_url_and_state(self, client_id, redirect_uri, scope, authorize_url): # pylint: disable=no-self-use
        """
        Mockable helper method to do the fetching of the authorization URL and state from GitHub.
        """
        return cla.utils.get_authorization_url_and_state(client_id, redirect_uri,
                                                         scope, authorize_url)

    def oauth2_redirect(self, state, code, repository_id, change_request_id, request): # pylint: disable=too-many-arguments
        """
        This is where the user will end up after having authorized the CLA system
        to get information such as email address. GitHub redirects the user here
        based on the redirect_uri we have provided in self.get_authorization_url_state()

        It will handle storing the OAuth2 session information for this user for
        further requests and initiate the signing workflow.
        """
        cla.log.info('Handling GitHub OAuth2 redirect')
        session = self._get_request_session(request)
        if state != session.get('github_oauth2_state', None):
            cla.log.warning('Invalid GitHub OAuth2 state')
            raise falcon.HTTPBadRequest('Invalid OAuth2 state', state)
        cla.log.info('Attempting to fetch OAuth2 token for state %s', state)
        client_id = cla.conf['GITHUB_CLIENT_ID']
        state = session.get('github_oauth2_state')
        token_url = cla.conf['GITHUB_OAUTH_TOKEN_URL']
        client_secret = cla.conf['GITHUB_SECRET']
        token = self._fetch_token(client_id, state, token_url, client_secret, code)
        cla.log.info('OAuth2 token received for state %s: %s', state, token)
        session['github_oauth2_token'] = token
        # Load the repository.
        repository = cla.utils.get_repository_instance()
        try:
            repository.load(repository_id)
        except DoesNotExist as err:
            cla.log.error('Repository not found during GitHub OAuth2 redirect for change ' + \
                          'request ID %s: %s',
                          change_request_id, repository_id)
            return {'errors': {'repository_id': str(err)}}
        return self.sign_workflow(repository, change_request_id, request)

    def _fetch_token(self, client_id, state, token_url, client_secret, code): # pylint: disable=too-many-arguments,no-self-use
        """
        Mockable method to fetch a OAuth2Session token.
        """
        return cla.utils.fetch_token(client_id, state, token_url, client_secret, code)

    def sign_workflow(self, repository, pull_request_number, request):
        """
        Once we have the 'github_oauth2_token' value in the user's session, we can initiate the
        signing workflow.
        """
        github_repository_id = repository.get_repository_external_id()
        cla.log.info('Initiating GitHub signing workflow for GitHub repo %s PR: %s',
                     github_repository_id, pull_request_number)
        user = self.get_or_create_user(request)
        agreement = cla.utils.get_user_agreement_by_repository(repository, user)

        if agreement is not None:
            cla.utils.redirect_user_by_agreement(user, agreement)
        else:
            # Agreement not found, create new one and send user to sign.
            cla.utils.request_signature(repository, user, pull_request_number)

    def process_opened_pull_request(self, data):
        """
        Helper method to handle a webhook fired from GitHub for an opened PR.

        :param data: The data returned from GitHub on this webhook.
        :type data: dict
        """
        pull_request_id = data['pull_request']['number']
        github_repository_id = data['repository']['id']
        repository_instance = cla.utils.get_repository_instance()
        repository = repository_instance.get_repository_by_external_id(github_repository_id,
                                                                       'github')
        if repository is None:
            cla.log.info('PR fired for GitHub repository %s, but not found in ' + \
                         'CLA system', github_repository_id)
            if not cla.conf['AUTO_CREATE_REPOSITORY']:
                cla.log.warning('AUTO_CREATE_REPOSITORY is set to False - please manually ' + \
                                'create this GitHub project\'s repository in the CLA system ' + \
                                'database')
                return
            else:
                cla.log.info('AUTO_CREATE_REPOSITORY is set to True, creating repository for ' + \
                             'this GitHub project automatically')
                repository = create_repository(data)
        if repository is not None:
            self.update_change_request(repository, pull_request_id)

    def get_return_url(self, repository_id, change_request_id):
        pull_request = self.get_pull_request(repository_id, change_request_id)
        return pull_request.html_url

    def update_change_request(self, repository, change_request_id):
        github_repository_id = repository.get_repository_external_id()
        pull_request = self.get_pull_request(github_repository_id, change_request_id)
        # Get all unique users involved in this PR.
        commit_authors = get_pull_request_commit_authors(pull_request)
        # Find users who have signed and who have not signed.
        signed = []
        missing = []
        for commit, commit_author in commit_authors:
            if isinstance(commit_author, github.NamedUser.NamedUser):
                # Deal with GitHub user.
                handle_commit_from_github_user(repository,
                                               commit,
                                               commit_author,
                                               signed,
                                               missing)
            elif isinstance(commit_author, github.GitAuthor.GitAuthor):
                # Deal with non-github user (just email and name in commit).
                handle_commit_from_git_author(repository,
                                              commit,
                                              commit_author,
                                              signed,
                                              missing)
            else:
                # Couldn't find any author information.
                missing.append((commit.sha, None))
        update_pull_request(repository.get_repository_id(),
                            pull_request,
                            signed=signed,
                            missing=missing)

    def get_pull_request(self, repository_id, pull_request_number):
        """
        Helper method to get the pull request object from GitHub.

        :param repository_id: The ID of the GitHub repository.
        :type repository_id: int
        :param pull_request_number: The number (not ID) of the GitHub PR.
        :type pull_request_number: int
        """
        cla.log.debug('Getting PR %s from repository %s', pull_request_number, repository_id)
        repo = self.client.get_repo(int(repository_id))
        try:
            return repo.get_pull(int(pull_request_number))
        except UnknownObjectException:
            cla.log.error('Could not find pull request %s for repository %s - ensure it ' + \
                          'exists and that your personal access token has the "repo" scope enabled',
                          pull_request_number, repository_id)
        except BadCredentialsException as err:
            cla.log.error('Invalid GitHub credentials provided: %s', str(err))

    def get_or_create_user(self, request):
        """
        Helper method to either get or create a user based on the GitHub request made.

        :param request: The hug request object for this API call.
        :type request: Request
        """
        session = self._get_request_session(request)
        github_user = self.get_user_data(session, cla.conf['GITHUB_CLIENT_ID'])
        if 'error' in github_user:
            # Could not get GitHub user data - maybe user revoked CLA app permissions?
            session = self._get_request_session(request)
            del session['github_oauth2_state']
            del session['github_oauth2_token']
            cla.log.warning('Deleted OAuth2 session data - retrying token exchange next time')
            raise falcon.HTTPError('400 Bad Request', 'github_oauth2_token',
                                   'Token permissions have been rejected, please try again')
        if github_user['email'] is None:
            cla.log.warning('GitHub user has no verified or public email address: %s (%s)',
                            github_user['name'], github_user['login'])
            raise falcon.HTTPError(
                '412 Precondition Failed', 'email',
                'Please verify and make public at least one email address with GitHub')
        cla.log.debug('Trying to load GitHub user by GitHub ID: %s', github_user['id'])
        user = cla.utils.get_user_instance().get_user_by_github_id(github_user['id'])
        if user is not None:
            cla.log.info('Loaded GitHub user by GitHub ID: %s - %s (%s)',
                         user.get_user_name(),
                         user.get_user_email(),
                         user.get_user_github_id())
            return user
        # User not found by GitHub ID, trying by email.
        cla.log.debug('Could not find GitHub user by GitHub ID: %s', github_user['id'])
        user = cla.utils.get_user_instance().get_user_by_email(github_user['email'])
        if user is not None:
            # Found user by email, set the GitHub ID
            user.set_user_github_id(github_user['id'])
            user.save()
            cla.log.info('Loaded GitHub user by email: %s - %s (%s)',
                         user.get_user_name(),
                         user.get_user_email(),
                         user.get_user_github_id())
            return user
        # User not found, create.
        cla.log.debug('Could not find GitHub user by email: %s', github_user['email'])
        cla.log.info('Creating new GitHub user %s - %s (%s)',
                     github_user['name'],
                     github_user['email'],
                     github_user['id'])
        user = cla.utils.get_user_instance()
        user.set_user_id(str(uuid.uuid4()))
        user.set_user_email(github_user['email'])
        user.set_user_name(github_user['name'])
        user.set_user_github_id(github_user['id'])
        user.save()
        return user

    def get_user_data(self, session, client_id): # pylint: disable=no-self-use
        """
        Mockable method to get user data. Returns all GitHub user data we have
        on the user based on the current OAuth2 session.

        :param session: The current user session.
        :type session: dict
        :param client_id: The GitHub OAuth2 client ID.
        :type session: string
        """
        token = session['github_oauth2_token']
        oauth2 = OAuth2Session(client_id, token=token)
        request = oauth2.get('https://api.github.com/user')
        github_user = request.json()
        cla.log.debug('GitHub user data: %s', github_user)
        if 'message' in github_user:
            cla.log.error('Could not get user data with OAuth2 token: %s', github_user['message'])
            return {'error': 'Could not get user data: %s' %github_user['message']}
        return github_user

    def process_reopened_pull_request(self, data):
        """
        Helper method to process a re-opened GitHub PR.

        Simply calls the self.process_opened_pull_request() method with the data provided.

        :param data: The data provided by the GitHub webhook.
        :type data: dict
        """
        return self.process_opened_pull_request(data)

    def process_closed_pull_request(self, data):
        """
        Helper method to process a closed GitHub PR.

        :param data: The data provided by the GitHub webhook.
        :type data: dict
        """
        pass

    def process_synchronized_pull_request(self, data):
        """
        Helper method to process a synchronized GitHub PR.

        Should be called when a new commit comes through on the PR.
        Simply calls the self.process_opened_pull_request() method with the data provided.
        This should re-check all commits for author information.

        :param data: The data provided by the GitHub webhook.
        :type data: dict
        """
        return self.process_opened_pull_request(data)


def create_repository(data):
    """
    Helper method to create a repository object in the CLA database given PR data.

    :param data: The data provided by the GitHub webhook.
    :type data: dict
    :return: The newly created repository object - already in the DB.
    :rtype: cla.models.model_interfaces.Repository
    """
    try:
        repository = cla.utils.get_repository_instance()
        repository.set_repository_id(str(uuid.uuid4()))
        # TODO: Need to use an ID unique across all repository providers instead of namespace.
        full_name = data['repository']['full_name']
        namespace = full_name.split('/')[0]
        repository.set_repository_project_id(namespace)
        repository.set_repository_external_id(data['repository']['id'])
        repository.set_repository_name(full_name)
        repository.set_repository_type('github')
        repository.set_repository_url(data['repository']['html_url'])
        repository.save()
        return repository
    except Exception as err:
        cla.log.error('Could not create GitHub repository automatically: %s', str(err))
        return None


def handle_commit_from_github_user(repository, commit, author, signed, missing): # pylint: disable=too-many-arguments
    """
    Helper method to triage commits between signed and not-signed user agreements.

    This method deals with GitHub users found in the commit information.

    :param repository: The repository this commit belongs to.
    :type repository: cla.models.model_interfaces.Repository
    :param commit: Commit object that we're handling.
    :type commit: github.Commit.Commit
    :param author: Author object holding information on the GitHub commit author.
    :type author: github.NamedUser.GitNamedUser
    :param signed: Reference to a list of signed authors so far. Should be modified
      in-place to add a signer if found.
    :type signed: [github.GitAuthor.GitAuthor | github.NamedUser.NamedUser]
    :param missing: Reference to a list of authors who have not signed yet.
        Should be modified in-place to add a missing signer if found.
    :type missing: [github.GitAuthor.GitAuthor | github.NamedUser.NamedUser]
    """
    user = cla.utils.get_user_instance().get_user_by_github_id(author.id)
    if user is None:
        # GitHub user not in system yet, agreement does not exist for this user.
        cla.log.info('GitHub user (%s - %s - %s) not found',
                     author.id, author.login, author.email)
        missing.append((commit.sha, author.name))
    else:
        cla.log.info('GitHub user found (%s - %s)',
                     user.get_user_email(), user.get_user_github_id())
        if cla.utils.user_signed_project_agreement(user, repository):
            signed.append((commit.sha, author.name))
        else:
            missing.append((commit.sha, author.name))


def handle_commit_from_git_author(repository, commit, author, signed, missing):
    """
    Helper method to triage commits between signed and not-signed user agreements.

    This method deals with non-GitHub users found in the commit information.

    :param repository: The repository this commit belongs to.
    :type repository: cla.models.model_interfaces.Repository
    :param commit: Commit object that we're handling.
    :type commit: github.Commit.Commit
    :param author: Author object holding information on the non-github commit author.
    :type author: github.GitAuthor.GitAuthor
    :param signed: Reference to a list of signed authors so far. Should be modified
      in-place to add a signer if found.
    :type signed: [(github.Commit.Commit, name)]
    :param missing: Reference to a list of authors who have not signed yet.
        Should be modified in-place to add a missing signer if found.
    :type missing: [(github.Commit.Commit, name)]
    """
    user = cla.utils.get_user_instance().get_user_by_email(author.email)
    if user is not None:
        cla.log.info('Git commit user found %s', user.get_user_email())
        # For now, accept non-github users as legitimate users.
        if cla.utils.user_signed_project_agreement(user, repository):
            signed.append((commit.sha, author.name))
        else:
            missing.append((commit.sha, author.name))
    else:
        cla.log.info('Git commit user (%s <%s>) not found',
                     author.name, author.email)
        missing.append((commit.sha, author.name))


def get_pull_request_commit_authors(pull_request):
    """
    Helper function to extract all committer information for a GitHub PR.

    :param pull_request: A GitHub pull request to examine.
    :type pull_request: GitHub.PullRequest
    :return: A list of tuples containing (commit, author).
    :rtype: [(github.Commit.Commit, string)]
    """
    commit_authors = []
    for commit in pull_request.get_commits():
        if commit.author is not None:
            cla.log.debug('GitHub author found for commit SHA %s: %s <%s>',
                          commit.sha, commit.author.id, commit.author.email)
            commit_authors.append((commit, commit.author))
        elif commit.commit.author is not None:
            # For now, trust that git commit author information is enough for verification.
            # TODO: This probably isn't enough - need to verify user somehow.
            cla.log.debug('No GitHub author found for commit SHA %s, using git author info: ' + \
                          '%s <%s>', commit.sha, commit.commit.author.name,
                          commit.commit.author.email)
            commit_authors.append((commit, commit.commit.author))
        else:
            cla.log.warning('Could not find commit author for SHA %s in PR %s',
                            commit.sha, pull_request.number)
            commit_authors.append((commit, None))
    return commit_authors


def update_pull_request(repository_id, pull_request, signed, missing): # pylint: disable=too-many-locals
    """
    Helper function to update a PR's comment/status based on the list of signers.

    :param repository_id: The ID of the repository this PR belongs to.
    :type repository_id: int
    :param pull_request: The GitHub PullRequest object for this PR.
    :type pull_request: GitHub.PullRequest
    :param signed: The list of (commit hash, author name) tuples that have signed an
        agreement for this PR.
    :type signed: [(string, string)]
    :param missing: The list of (commit hash, author name) tuples that have not signed
        an agreement for this PR.
    :type missing: [(string, string)]
    """
    notification = cla.conf['GITHUB_PR_NOTIFICATION']
    both = notification == 'status+comment' or notification == 'comment+status'
    if both or notification == 'comment':
        body = cla.utils.assemble_cla_comment('github', repository_id, pull_request.number,
                                              signed, missing)
        update_cla_comment(pull_request, body)
    if both or notification == 'status':
        state = 'pending'
        for commit, author_name in missing:
            body = cla.utils.assemble_cla_status(author_name, signed=False)
            sign_url = cla.utils.get_full_sign_url('github', repository_id, pull_request.number)
            cla.log.info('Creating new CLA status on commit %s: %s', commit, state)
            create_commit_status(pull_request, commit, state, sign_url, body)
        state = 'success'
        for commit, author_name in signed:
            body = cla.utils.assemble_cla_status(author_name, signed=True)
            sign_url = cla.utils.get_full_sign_url('github', repository_id, pull_request.number)
            cla.log.info('Creating new CLA status on commit %s: %s', commit, state)
            create_commit_status(pull_request, commit, state, sign_url, body)
        num_missing = len(missing)
        if num_missing > 0:
            # Need to update the last status message to prevent merging the PR.
            last_commit = pull_request.get_commits().reversed[0]
            signed_commits = [item[0] for item in signed]
            if last_commit.sha in signed_commits:
                num_signed = len(signed)
                total = num_signed + len(missing)
                last_commit.create_status('pending', sign_url,
                                          'Missing CLA signatures (%s/%s)' %(num_signed, total))


def create_commit_status(pull_request, commit_hash, state, sign_url, body):
    """
    Helper function to create a pull request commit status message given the PR and commit hash.

    :param pull_request: The GitHub Pull Request object.
    :type pull_request: github.PullRequest
    :param commit_hash: The commit hash to post a status on.
    :type commit_hash: string
    :param state: The state of the status.
    :type state: string
    :param sign_url: The link the user will be taken to when clicking on the status message.
    :type sign_url: string
    :param body: The contents of the status message.
    :type body: string
    """
    commit_obj = None
    for commit in pull_request.get_commits():
        if commit.sha == commit_hash:
            commit_obj = commit
            break
    if commit_obj is None:
        cla.log.error('Could not post status on PR %s: Commit %s not found',
                      pull_request.number, commit_hash)
        return
    commit_obj.create_status(state, sign_url, body)


def update_cla_comment(pull_request, body):
    """
    Helper function to create/edit a comment on the GitHub PR.

    :param pull_request: The PR object in question.
    :type pull_request: GitHub.PullRequest
    :param body: The contents of the comment.
    :type body: string
    """
    comment = get_existing_cla_comment(pull_request)
    if comment is not None:
        cla.log.info('Updating existing CLA comment for PR: %s', pull_request.number)
        cla.log.debug('Updated CLA comment for PR %s: %s', pull_request.number, body)
        comment.edit(body)
    else:
        cla.log.info('Creating new CLA comment for PR: %s', pull_request.number)
        cla.log.debug('New comment for PR %s: %s', pull_request.number, body)
        pull_request.create_issue_comment(body)


def get_existing_cla_comment(pull_request):
    """
    Helper function to get an existing comment from the CLA system in a GitHub PR.

    :param pull_request: The PR object in question.
    :type pull_request: GitHub.PullRequest
    """
    comments = pull_request.get_issue_comments()
    for comment in comments:
        if '[![CLA Check](' in comment.body:
            cla.log.info('Found matching CLA comment for PR: %s', pull_request.number)
            return comment


class MockGitHub(GitHub):
    """
    The GitHub repository service mock class for testing.
    """
    def __init__(self, oauth2_token=False):
        super().__init__()
        self.oauth2_token = oauth2_token

    def _get_github_client(self, username, token):
        return MockGitHubClient(username, token)

    def _get_authorization_url_and_state(self, client_id, redirect_uri, scope, authorize_url):
        authorization_url = 'http://authorization.url'
        state = 'random-state-here'
        return authorization_url, state

    def _fetch_token(self, client_id, state, token_url, client_secret, code): # pylint: disable=too-many-arguments
        return 'random-token'

    def _get_request_session(self, request):
        if self.oauth2_token:
            return {'github_oauth2_token': 'random-token',
                    'github_oauth2_state': 'random-state'}
        return {}

    def get_user_data(self, session, client_id):
        return {'email': 'test@user.com', 'name': 'Test User', 'id': 123}


class MockGitHubClient(object): # pylint: disable=too-few-public-methods
    """
    The GitHub Client object mock class for testing.
    """
    def __init__(self, username, token):
        self.username = username
        self.token = token

    def get_repo(self, repository_id): # pylint: disable=no-self-use
        """
        Mock version of the GitHub Client object's get_repo method.
        """
        return MockGitHubRepository(repository_id)


class MockGitHubRepository(object): # pylint: disable=too-few-public-methods
    """
    The GitHub Repository object mock class for testing.
    """
    def __init__(self, repository_id):
        self.id = repository_id

    def get_pull(self, pull_request_id): # pylint: disable=no-self-use
        """
        Mock version of the GitHub Repository object's get_pull method.
        """
        return MockGitHubPullRequest(pull_request_id)


class MockGitHubPullRequest(object): # pylint: disable=too-few-public-methods
    """
    The GitHub PullRequest object mock class for testing.
    """
    def __init__(self, pull_request_id):
        self.number = pull_request_id
        self.html_url = 'http://test-github.com/user/repo/' + str(self.number)

    def get_commits(self): # pylint: disable=no-self-use
        """
        Mock version of the GitHub PullRequest object's get_commits method.
        """
        lst = MockPaginatedList()
        lst._elements = [MockGitHubCommit()] # pylint: disable=protected-access
        return lst

    def get_issue_comments(self): # pylint: disable=no-self-use
        """
        Mock version of the GitHub PullRequest object's get_issue_comments method.
        """
        return [MockGitHubComment()]

    def create_issue_comment(self, body): # pylint: disable=no-self-use
        """
        Mock version of the GitHub PullRequest object's create_issue_comment method.
        """
        pass


class MockGitHubComment(object): # pylint: disable=too-few-public-methods
    """
    A GitHub mock issue comment object for testing.
    """
    body = 'Test'


class MockPaginatedList(github.PaginatedList.PaginatedListBase): # pylint: disable=too-few-public-methods
    """Mock GitHub paginated list for testing purposes."""
    def __init__(self):
        super().__init__()
        # Need to use our own elements list (self.__elements from PaginatedListBase does not
        # work as expected).
        self._elements = []

    @property
    def reversed(self):
        """Fake reversed propery."""
        return [MockGitHubCommit()]

    def __iter__(self):
        for element in self._elements:
            yield element


class MockGitHubCommit(object): # pylint: disable=too-few-public-methods
    """
    The GitHub Commit object mock class for testing.
    """
    def __init__(self):
        self.author = MockGitHubAuthor()
        self.sha = 'sha-test-commit'

    def create_status(self, state, sign_url, body):
        """
        Mock version of the GitHub Commit object's create_status method.
        """
        pass


class MockGitHubAuthor(object): # pylint: disable=too-few-public-methods
    """
    The GitHub Author object mock class for testing.
    """
    def __init__(self, author_id=1):
        self.id = author_id
        self.login = 'user'
        self.email = 'user@github.com'
