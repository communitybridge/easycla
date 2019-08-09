# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Holds the GitHub repository service.
"""

import os
import uuid

import falcon
import github
from github.GithubException import UnknownObjectException, BadCredentialsException
from requests_oauthlib import OAuth2Session

import cla
from cla.controllers.github_application import GitHubInstallation
from cla.models import repository_service_interface, DoesNotExist
from cla.models.dynamo_models import Repository, GitHubOrg


class GitHub(repository_service_interface.RepositoryService):
    """
    The GitHub repository service.
    """
    def __init__(self):
        self.client = None

    def initialize(self, config):
        # username = config['GITHUB_USERNAME']
        # token = config['GITHUB_TOKEN']
        # self.client = self._get_github_client(username, token)
        pass

    def _get_github_client(self, username, token):  # pylint: disable=no-self-use
        return github.Github(username, token)

    def get_repository_id(self, repo_name, installation_id=None):
        """
        Helper method to get a GitHub repository ID based on repository name.

        :param repo_name: The name of the repository, example: 'linuxfoundation/cla'.
        :type repo_name: string
        :param installation_id: The github installation id
        :type installation_id: string
        :return: The repository ID.
        :rtype: integer
        """
        if installation_id is not None:
            self.client = get_github_integration_client(installation_id)
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

    def sign_request(self, installation_id, github_repository_id, change_request_id, request):
        """
        This method gets called when the OAuth2 app (NOT the GitHub App) needs to get info on the
        user trying to sign. In this case we begin an OAuth2 exchange with the 'user:email' scope.
        """
        cla.log.info('Initiating GitHub sign request for repository %s', github_repository_id)
        # Not sure if we need a different token for each installation ID...
        session = self._get_request_session(request)
        session['github_installation_id'] = installation_id
        session['github_repository_id'] = github_repository_id
        session['github_change_request_id'] = change_request_id

        origin_url = self.get_return_url(github_repository_id, change_request_id, installation_id)
        session['github_origin_url'] = origin_url
        if 'github_oauth2_token' in session:
            cla.log.info('Using existing session OAuth2 token')
            return self.redirect_to_console(installation_id, github_repository_id, change_request_id, origin_url, request)
        else:
            cla.log.info('Initiating GitHub OAuth2 exchange')
            authorization_url, state = self.get_authorization_url_and_state(installation_id,
                                                                            github_repository_id,
                                                                            change_request_id,
                                                                            ['user:email'])
            session['github_oauth2_state'] = state
            cla.log.info('GitHub OAuth2 request with state %s - sending user to %s',
                         state, authorization_url)
            raise falcon.HTTPFound(authorization_url)

    def _get_request_session(self, request):  # pylint: disable=no-self-use
        """
        Mockable method used to get the current user session.
        """
        return request.context['session']

    def get_authorization_url_and_state(self, installation_id, github_repository_id, pull_request_number, scope):
        """
        Helper method to get the GitHub OAuth2 authorization URL and state.

        This will be used to get the user's emails from GitHub.

        :TODO: Update comments.

        :param repository_id: The ID of the repository this request was initiated in.
        :type repository_id: int
        :param pull_request_number: The PR number this request was generated in.
        :type pull_request_number: int
        :param scope: The list of OAuth2 scopes to request from GitHub.
        :type scope: [string]
        """
        # Get the PR's html_url property.
        # origin = self.get_return_url(github_repository_id, pull_request_number, installation_id)
        # Add origin to user's session here?
        api_base_url = os.environ.get('CLA_API_BASE', '')
        return self._get_authorization_url_and_state(os.environ['GH_OAUTH_CLIENT_ID'],
                                                     os.path.join(api_base_url, 'v2/github/installation'),
                                                     scope,
                                                     cla.conf['GITHUB_OAUTH_AUTHORIZE_URL'])

    def _get_authorization_url_and_state(self, client_id, redirect_uri, scope, authorize_url):  # pylint: disable=no-self-use
        """
        Mockable helper method to do the fetching of the authorization URL and state from GitHub.
        """
        return cla.utils.get_authorization_url_and_state(client_id, redirect_uri,
                                                         scope, authorize_url)

    def oauth2_redirect(self, state, code, request):  # pylint: disable=too-many-arguments
        """
        This is where the user will end up after having authorized the CLA system
        to get information such as email address.

        It will handle storing the OAuth2 session information for this user for
        further requests and initiate the signing workflow.
        """
        cla.log.info('Handling GitHub OAuth2 redirect')
        session = self._get_request_session(request)

        cla.log.debug('State: %s', state)
        cla.log.debug('Code: %s', code)
        cla.log.debug('Session: %s', session)

        if 'github_oauth2_state' in session:
            session_state = session['github_oauth2_state']
        else:
            session_state = None
            cla.log.warning('github_oauth2_state not set in session')

        if state != session_state:
            cla.log.warning('Invalid GitHub OAuth2 state %s expecting %s',
                            session_state, state)
            raise falcon.HTTPBadRequest('Invalid OAuth2 state', state)

        # Get session information for this request.
        cla.log.info('Attempting to fetch OAuth2 token for state %s', state)
        installation_id = session.get('github_installation_id', None)
        github_repository_id = session.get('github_repository_id', None)
        change_request_id = session.get('github_change_request_id', None)
        origin_url = session.get('github_origin_url', None)
        state = session.get('github_oauth2_state')
        token_url = cla.conf['GITHUB_OAUTH_TOKEN_URL']
        client_id = os.environ['GH_OAUTH_CLIENT_ID']
        client_secret = os.environ['GH_OAUTH_SECRET']
        token = self._fetch_token(client_id, state, token_url, client_secret, code)
        cla.log.info('OAuth2 token received for state %s: %s', state, token)
        session['github_oauth2_token'] = token
        return self.redirect_to_console(installation_id, github_repository_id, change_request_id, origin_url, request)

    def redirect_to_console(self, installation_id, repository_id, pull_request_id, redirect, request):
        console_endpoint = cla.conf['CONTRIBUTOR_BASE_URL']
        # Get repository using github's repository ID.
        repository = Repository().get_repository_by_external_id(repository_id, "github")
        if repository is None:  
            cla.log.error('Could not find repository with the following repository_id: %s', repository_id)
            return None
        
        # Get project ID from this repository
        project_id = repository.get_repository_project_id()

        user = self.get_or_create_user(request)
        # Ensure user actually requires a signature for this project.
        # TODO: Skipping this for now - we can do this for ICLAs but there's no easy way of doing
        # the check for CCLAs as we need to know in advance what the company_id is that we're checking
        # the CCLA signature for.
        # We'll have to create a function that fetches the latest CCLA regardless of company_id.
        # icla_signature = cla.utils.get_user_signature_by_github_repository(installation_id, user)
        # ccla_signature = cla.utils.get_user_signature_by_github_repository(installation_id, user, company_id=?)
        #try:
            #document = cla.utils.get_project_latest_individual_document(project_id)
        #except DoesNotExist:
            #cla.log.info('No ICLA for project %s' %project_id)
        #if signature is not None and \
            #signature.get_signature_document_major_version() == document.get_document_major_version():
            #return cla.utils.redirect_user_by_signature(user, signature)
        # Store repository and PR info so we can redirect the user back later.
        cla.utils.set_active_signature_metadata(user.get_user_id(), project_id, repository_id, pull_request_id)
        # Generate console URL
        console_url = 'https://' + console_endpoint + \
                      '/#/cla/project/' + project_id + \
                      '/user/' + user.get_user_id() + \
                      '?redirect=' + redirect
        raise falcon.HTTPFound(console_url)

    def _fetch_token(self, client_id, state, token_url, client_secret, code):  # pylint: disable=too-many-arguments,no-self-use
        """
        Mockable method to fetch a OAuth2Session token.
        """
        return cla.utils.fetch_token(client_id, state, token_url, client_secret, code)

    def sign_workflow(self, installation_id, github_repository_id, pull_request_number, request):
        """
        Once we have the 'github_oauth2_token' value in the user's session, we can initiate the
        signing workflow.
        """
        cla.log.info('Initiating GitHub signing workflow for GitHub repo %s PR: %s',
                     github_repository_id, pull_request_number)
        user = self.get_or_create_user(request)
        signature = cla.utils.get_user_signature_by_github_repository(installation_id, user)
        project_id = cla.utils.get_project_id_from_installation_id(installation_id)
        document = cla.utils.get_project_latest_individual_document(project_id)
        if signature is not None and \
           signature.get_signature_document_major_version() == document.get_document_major_version():
            return cla.utils.redirect_user_by_signature(user, signature)
        else:
            # Signature not found or older version, create new one and send user to sign.
            cla.utils.request_individual_signature(installation_id, github_repository_id, user, pull_request_number)

    def process_opened_pull_request(self, data):
        """
        Helper method to handle a webhook fired from GitHub for an opened PR.

        :param data: The data returned from GitHub on this webhook.
        :type data: dict
        """
        pull_request_id = data['pull_request']['number']
        github_repository_id = data['repository']['id']
        installation_id = data['installation']['id']
        self.update_change_request(installation_id, github_repository_id, pull_request_id)

    def get_return_url(self, github_repository_id, change_request_id, installation_id):
        pull_request = self.get_pull_request(github_repository_id, change_request_id, installation_id)
        return pull_request.html_url

    def update_change_request(self, installation_id, github_repository_id, change_request_id):
        pull_request = self.get_pull_request(github_repository_id,
                                             change_request_id,
                                             installation_id)
        # Get all unique users involved in this PR.
        commit_authors = get_pull_request_commit_authors(pull_request)
        # Get existing repository info using the repository's external ID, which is the repository ID assigned by github. 
        
        try: 
            repository = Repository().get_repository_by_external_id(github_repository_id, "github")
        except DoesNotExist:
            cla.log.error('Could not find repository with the repository ID: %s', github_repository_id)
            cla.log.error('Failed to update change request %s of repository %s', change_request_id, github_repository_id)
            return

        # Get Github Organization name that the repository is configured to. 
        organization_name = repository.get_repository_organization_name()

        # Check that the Github Organization exists.
        github_org = GitHubOrg()
        try:
            github_org.load(organization_name)
        except DoesNotExist:
            cla.log.error('Could not find Github Organization with the following organization name: %s', organization_name)
            cla.log.error('Failed to update change request %s of repository %s', change_request_id, github_repository_id)
            return 
    
        # Ensure that installation ID for this organization matches the given installation ID
        if github_org.get_organization_installation_id() != installation_id:
            cla.log.error('The installation ID: %s of this organization does not match installation ID: %s given by the pull request.', 
                                                                        github_org.get_organization_installation_id(), installation_id)
            cla.log.error('Failed to update change request %s of repository %s', change_request_id, github_repository_id)
            return

        # Retrieve project ID from the repository. 
        project_id = repository.get_repository_project_id() 

        # Find users who have signed and who have not signed.
        signed = []
        missing = []
        for commit, commit_author in commit_authors:
            if isinstance(commit_author, github.NamedUser.NamedUser):
                # Handle GitHub user.
                cla.log.info("Handle GitHub user")
                handle_commit_from_github_user(project_id,
                                               commit,
                                               commit_author,
                                               signed,
                                               missing)
            elif isinstance(commit_author, github.GitAuthor.GitAuthor):
                # Handle non-github user (just email and name in commit).
                cla.log.info("Handle non-github user (just email and name in commit)")
                handle_commit_from_git_author(project_id,
                                              commit,
                                              commit_author,
                                              signed,
                                              missing)
            else:
                # Couldn't find any author information.
                cla.log.info("Couldn't find any author information for author: {}.".format(commit_author))
                if commit_author is not None:
                    missing.append((commit.sha, commit_author))
                else:
                    missing.append((commit.sha, None))

        cla.log.debug('updating github pull request for repo: {}, '
                      'pr: {} with signed authors: {} with missing authors: {}'.
                      format(github_repository_id, pull_request, signed, missing))
        update_pull_request(installation_id,
                            github_repository_id,
                            pull_request,
                            signed=signed,
                            missing=missing)

    def get_pull_request(self, github_repository_id, pull_request_number, installation_id):
        """
        Helper method to get the pull request object from GitHub.

        :param github_repository_id: The ID of the GitHub repository.
        :type github_repository_id: int
        :param pull_request_number: The number (not ID) of the GitHub PR.
        :type pull_request_number: int
        :param installation_id: The ID of the GitHub application installed on this repository.
        :type installation_id: int | None
        """
        cla.log.debug('Getting PR %s from GitHub repository %s', pull_request_number, github_repository_id)
        if self.client is None:
            self.client = get_github_integration_client(installation_id)
        repo = self.client.get_repo(int(github_repository_id))
        try:
            return repo.get_pull(int(pull_request_number))
        except UnknownObjectException:
            cla.log.error('Could not find pull request %s for repository %s - ensure it '
                          'exists and that your personal access token has the "repo" scope enabled',
                          pull_request_number, github_repository_id)
        except BadCredentialsException as err:
            cla.log.error('Invalid GitHub credentials provided: %s', str(err))

    def get_or_create_user(self, request):
        """
        Helper method to either get or create a user based on the GitHub request made.

        :param request: The hug request object for this API call.
        :type request: Request
        """
        session = self._get_request_session(request)
        github_user = self.get_user_data(session, os.environ['GH_OAUTH_CLIENT_ID'])
        if 'error' in github_user:
            # Could not get GitHub user data - maybe user revoked CLA app permissions?
            session = self._get_request_session(request)
            del session['github_oauth2_state']
            del session['github_oauth2_token']
            cla.log.warning('Deleted OAuth2 session data - retrying token exchange next time')
            raise falcon.HTTPError('400 Bad Request', 'github_oauth2_token',
                                   'Token permissions have been rejected, please try again')

        emails = self.get_user_emails(session, os.environ['GH_OAUTH_CLIENT_ID'])
        if len(emails) < 1:
            cla.log.warning('GitHub user has no verified email address: %s (%s)',
                            github_user['name'], github_user['login'])
            raise falcon.HTTPError(
                '412 Precondition Failed', 'email',
                'Please verify at least one email address with GitHub')

        cla.log.debug('Trying to load GitHub user by GitHub ID: %s', github_user['id'])
        user = cla.utils.get_user_instance().get_user_by_github_id(github_user['id'])
        if user is not None:
            cla.log.info('Loaded GitHub user by GitHub ID: %s - %s (%s)',
                         user.get_user_name(),
                         user.get_user_emails(),
                         user.get_user_github_id())
            user.set_user_emails(emails)
            user.save()
            return user

        # User not found by GitHub ID, trying by email.
        cla.log.debug('Could not find GitHub user by GitHub ID: %s', github_user['id'])
        # TODO: This is very slow and needs to be improved - may need a DB schema change.
        found = None
        user = cla.utils.get_user_instance()
        for email in emails:
            found = user.get_user_by_email(email)
            if found is not None:
                break

        if found is not None:
            # Found user by email, set the GitHub ID
            found.set_user_github_id(github_user['id'])
            found.set_user_emails(emails)
            found.save()
            cla.log.info('Loaded GitHub user by email: %s - %s (%s)',
                         found.get_user_name(),
                         found.get_user_emails(),
                         found.get_user_github_id())
            return found

        # User not found, create.
        cla.log.debug('Could not find GitHub user by email: %s', emails)
        cla.log.info('Creating new GitHub user %s - %s (%s)',
                     github_user['name'],
                     emails,
                     github_user['id'])
        user = cla.utils.get_user_instance()
        user.set_user_id(str(uuid.uuid4()))
        user.set_user_emails(emails)
        user.set_user_name(github_user['name'])
        user.set_user_github_id(github_user['id'])
        user.set_user_github_username(github_user['login'])
        user.save()
        return user

    def get_user_data(self, session, client_id):  # pylint: disable=no-self-use
        """
        Mockable method to get user data. Returns all GitHub user data we have
        on the user based on the current OAuth2 session.

        :param session: The current user session.
        :type session: dict
        :param client_id: The GitHub OAuth2 client ID.
        :type client_id: string
        """
        token = session['github_oauth2_token']
        oauth2 = OAuth2Session(client_id, token=token)
        request = oauth2.get('https://api.github.com/user')
        github_user = request.json()
        cla.log.debug('GitHub user data: %s', github_user)
        if 'message' in github_user:
            cla.log.error('Could not get user data with OAuth2 token: %s', github_user['message'])
            return {'error': 'Could not get user data: %s' % github_user['message']}
        return github_user

    def get_user_emails(self, session, client_id):  # pylint: disable=no-self-use
        """
        Mockable method to get all user emails based on OAuth2 session.

        :param session: The current user session.
        :type session: dict
        :param client_id: The GitHub OAuth2 client ID.
        :type client_id: string
        """
        token = session['github_oauth2_token']
        oauth2 = OAuth2Session(client_id, token=token)
        request = oauth2.get('https://api.github.com/user/emails')
        emails = request.json()
        cla.log.debug('GitHub user emails: %s', emails)
        if 'message' in emails:
            cla.log.error('Could not get user emails with OAuth2 token: %s', emails['message'])
            return {'error': 'Could not get user emails: %s' % emails['message']}
        return [item['email'] for item in emails if item['verified']]

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


def handle_commit_from_github_user(project_id, commit, author, signed, missing):  # pylint: disable=too-many-arguments
    """
    Helper method to triage commits between signed and not-signed user signatures.

    This method deals with GitHub users found in the commit information.

    :param project_id: The project ID for this github PR organization.
    :type project_id: string
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
    
    # Validate author name to show
    if author.name is not None:
        author_name = author.name # set name when available
    else:
        author_name = author.login # set username (login) when name not available

    user = cla.utils.get_user_instance().get_user_by_github_id(author.id)
    if user is None:
        # GitHub user not in system yet, signature does not exist for this user.
        cla.log.info('GitHub user (%s - %s - %s) not found, looking up user by email',
                     author.id, author.login, author.email)
        # Try looking up user by email as a fallback
        handle_commit_from_git_author(project_id, commit, author, signed, missing)
    else:
        cla.log.info('GitHub user found (%s - %s)',
                     user.get_user_emails(), user.get_user_github_id())
        if cla.utils.user_signed_project_signature(user, project_id):
            signed.append((commit.sha, author_name))
        else:
            missing.append((commit.sha, author_name))

def handle_commit_from_git_author(project_id, commit, author, signed, missing):
    """
    Helper method to triage commits between signed and not-signed user signatures.

    This method deals with non-GitHub users found in the commit information.

    :param project_id: The project ID for this github PR organization.
    :type project_id: string
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
        cla.log.info('Git commit user found (lookup by email) %s', user.get_user_emails())
        # For now, accept non-github users as legitimate users.
        if cla.utils.user_signed_project_signature(user, project_id):
            signed.append((commit.sha, author.name))
        else:
            missing.append((commit.sha, author.name))
    else:
        cla.log.info('Git commit user (lookup by email) (%s <%s>) not found',
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
            cla.log.info('GitHub author found for commit SHA %s: %s <%s>',
                          commit.sha, commit.author.id, commit.author.email)
            # commit_authors.append((commit, commit.author.login))
            commit_authors.append((commit, commit.author))
        elif commit.commit.author is not None:
            # For now, trust that git commit author information is enough for verification.
            # TODO: This probably isn't enough - need to verify user somehow.
            cla.log.info('No GitHub author found for commit SHA %s, using git author info: ' + \
                          '%s <%s>', commit.sha, commit.commit.author.name,
                          commit.commit.author.email)
            # commit_authors.append((commit, commit.commit.author))
            commit_authors.append((commit, commit.commit))
        else:
            cla.log.warning('Could not find commit author for SHA %s in PR %s',
                            commit.sha, pull_request.number)
            commit_authors.append((commit, None))
    return commit_authors


def update_pull_request(installation_id, github_repository_id, pull_request, signed, missing):  # pylint: disable=too-many-locals
    """
    Helper function to update a PR's comment/status based on the list of signers.

    :TODO: Update comments.

    :param repository_id: The ID of the repository this PR belongs to.
    :type repository_id: int
    :param pull_request: The GitHub PullRequest object for this PR.
    :type pull_request: GitHub.PullRequest
    :param signed: The list of (commit hash, author name) tuples that have signed an
        signature for this PR.
    :type signed: [(string, string)]
    :param missing: The list of (commit hash, author name) tuples that have not signed
        an signature for this PR.
    :type missing: [(string, string)]
    """
    notification = cla.conf['GITHUB_PR_NOTIFICATION']
    both = notification == 'status+comment' or notification == 'comment+status'
    last_commit = pull_request.get_commits().reversed[0]

    # Here we update the PR status by adding/updating the PR body - this is the way the EasyCLA app
    # knows if it is pass/fail.

    if both or notification == 'comment':
        body = cla.utils.assemble_cla_comment('github', installation_id, github_repository_id, pull_request.number,
                                              signed, missing)
        update_cla_comment(pull_request, body)
        if not missing:
            cla.log.debug('EasyCLA App checks pass for PR: {} with authors: {}'.format(pull_request.number, signed))
        else:
            cla.log.debug('EasyCLA App checks fail for PR: {}. CLA signatures with signed authors: {} and '
                          'with missing authors: {}'.format(pull_request.number, signed, missing))

    if both or notification == 'status':
        state = 'failure'
        for commit, author_name in missing:
            context, body = cla.utils.assemble_cla_status(author_name, signed=False)
            sign_url = cla.utils.get_full_sign_url('github', installation_id, github_repository_id, pull_request.number)
            cla.log.info('Creating new CLA status on commit %s: %s', commit, state)
            create_commit_status(pull_request, last_commit.sha, state, sign_url, body, context)
        state = 'success'
        for commit, author_name in signed:
            context, body = cla.utils.assemble_cla_status(author_name, signed=True)
            # sign_url = cla.utils.get_full_sign_url('github', installation_id, github_repository_id, pull_request.number)
            cla.log.info('Creating new CLA status on commit %s: %s', commit, state)
            sign_url = "https://lfcla.com" # Remove this once signature detail page ready.
            create_commit_status(pull_request, last_commit.sha, state, sign_url, body, context)


def create_commit_status(pull_request, commit_hash, state, sign_url, body, context):
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
    # context is a string label to differentiate one signer status from another signer status.
    # committer name is used as context label
    commit_obj.create_status(state, sign_url, body, context)


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


def get_github_integration_client(installation_id):
    """
    GitHub App integration client used for authenticated client actions through an installed app.
    """
    return GitHubInstallation(installation_id).api_object

def get_github_client(organization_id):
    github_org = cla.utils.get_github_organization_instance()
    github_org.load(organization_id)
    installation_id = github_org.get_organization_installation_id()
    return get_github_integration_client(installation_id)


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

    def _fetch_token(self, client_id, state, token_url, client_secret, code):  # pylint: disable=too-many-arguments
        return 'random-token'

    def _get_request_session(self, request):
        if self.oauth2_token:
            return {'github_oauth2_token': 'random-token',
                    'github_oauth2_state': 'random-state',
                    'github_origin_url': 'http://github/origin/url',
                    'github_installation_id': 1}
        return {}

    def get_user_data(self, session, client_id):
        return {'email': 'test@user.com', 'name': 'Test User', 'id': 123}

    def get_user_emails(self, session, client_id):
        return [{'email': 'test@user.com', 'verified': True, 'primary': True, 'visibility': 'public'}]

    def get_pull_request(self, github_repository_id, pull_request_number, installation_id):
        return MockGitHubPullRequest(pull_request_number)

class MockGitHubClient(object):  # pylint: disable=too-few-public-methods
    """
    The GitHub Client object mock class for testing.
    """
    def __init__(self, username, token):
        self.username = username
        self.token = token

    def get_repo(self, repository_id):  # pylint: disable=no-self-use
        """
        Mock version of the GitHub Client object's get_repo method.
        """
        return MockGitHubRepository(repository_id)


class MockGitHubRepository(object):  # pylint: disable=too-few-public-methods
    """
    The GitHub Repository object mock class for testing.
    """
    def __init__(self, repository_id):
        self.id = repository_id

    def get_pull(self, pull_request_id):  # pylint: disable=no-self-use
        """
        Mock version of the GitHub Repository object's get_pull method.
        """
        return MockGitHubPullRequest(pull_request_id)


class MockGitHubPullRequest(object):  # pylint: disable=too-few-public-methods
    """
    The GitHub PullRequest object mock class for testing.
    """
    def __init__(self, pull_request_id):
        self.number = pull_request_id
        self.html_url = 'http://test-github.com/user/repo/' + str(self.number)

    def get_commits(self):  # pylint: disable=no-self-use
        """
        Mock version of the GitHub PullRequest object's get_commits method.
        """
        lst = MockPaginatedList()
        lst._elements = [MockGitHubCommit()]  # pylint: disable=protected-access
        return lst

    def get_issue_comments(self):  # pylint: disable=no-self-use
        """
        Mock version of the GitHub PullRequest object's get_issue_comments method.
        """
        return [MockGitHubComment()]

    def create_issue_comment(self, body):  # pylint: disable=no-self-use
        """
        Mock version of the GitHub PullRequest object's create_issue_comment method.
        """
        pass


class MockGitHubComment(object):  # pylint: disable=too-few-public-methods
    """
    A GitHub mock issue comment object for testing.
    """
    body = 'Test'


class MockPaginatedList(github.PaginatedList.PaginatedListBase):  # pylint: disable=too-few-public-methods
    """Mock GitHub paginated list for testing purposes."""
    def __init__(self):
        super().__init__()
        # Need to use our own elements list (self.__elements from PaginatedListBase does not
        # work as expected).
        self._elements = []

    @property
    def reversed(self):
        """Fake reversed property."""
        return [MockGitHubCommit()]

    def __iter__(self):
        for element in self._elements:
            yield element


class MockGitHubCommit(object):  # pylint: disable=too-few-public-methods
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


class MockGitHubAuthor(object):  # pylint: disable=too-few-public-methods
    """
    The GitHub Author object mock class for testing.
    """
    def __init__(self, author_id=1):
        self.id = author_id
        self.login = 'user'
        self.email = 'user@github.com'
