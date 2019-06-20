# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Holds the GitLab repository service.
"""

import uuid
import falcon
import gitlab
from requests_oauthlib import OAuth2Session
import cla
from cla.models import repository_service_interface, DoesNotExist

class GitLab(repository_service_interface.RepositoryService):
    """
    The GitLab repository service.
    """
    def __init__(self):
        self.client = None

    def initialize(self, config):
        domain = config['GITLAB_DOMAIN']
        token = config['GITLAB_TOKEN']
        self.client = self._get_gitlab_client(domain, token)

    def _get_gitlab_client(self, domain, token): # pylint: disable=no-self-use
        return gitlab.Gitlab(domain, token, api_version=4)

    def get_repository_id(self, repo_name):
        """
        Helper method to get a GitLab repository ID based on the repository name.

        :param repo_name: The name of the repository, example: 'repo/cla-enabled-repo'.
        :type repo_name: string
        :return: The repository ID for that repository in GitLab.
        :rtype: integer
        """
        try:
            return self.client.projects.get(repo_name).id
        except Exception as err:
            cla.log.error('Unkown error while getting GitLab repository ID for repository %s: %s',
                          repo_name, str(err))

    def received_activity(self, data):
        cla.log.debug('Received GitLab activity: %s', data)
        if 'object_kind' not in data or data['object_kind'] != 'merge_request':
            cla.log.debug('Activity not related to merge request')
            return {'message': 'Not a merge request - no action performed'}
        if data['object_attributes']['action'] == 'open':
            cla.log.debug('Handling opened merge request')
            return self.process_opened_merge_request(data)
        elif data['object_attributes']['action'] == 'reopen':
            cla.log.debug('Handling re-opened merge request')
            return self.process_reopened_merge_request(data)
        elif data['object_attributes']['action'] == 'close':
            cla.log.debug('Handling closed merge request')
            return self.process_closed_merge_request(data)

    def sign_request(self, repository_id, change_request_id, request):
        cla.log.info('Initiating GitLab sign request for repository %s', repository_id)
        session = self._get_request_session(request)
        repository = cla.utils.get_repository_instance()
        try:
            # Load the repository in question.
            repository.load(repository_id)
        except DoesNotExist as err:
            cla.log.error('Error when initiating GitLab sign request for repository %s and ' + \
                          'change request %s: %s', repository_id, change_request_id, str(err))
            return {'errors': {'repository_id': str(err)}}
        if 'gitlab_oauth2_token' in session:
            cla.log.info('Using existing session OAuth2 token')
            return self.sign_workflow(repository, change_request_id, request)
        else:
            cla.log.info('Initiating GitLab OAuth2 exchange')
            authorization_url, state = self.get_authorization_url_and_state(repository_id,
                                                                            change_request_id,
                                                                            ['read_user'])
            session['gitlab_oauth2_state'] = state
            cla.log.info('GitLab OAuth2 request with state %s - sending user to %s',
                         state, authorization_url)
            raise falcon.HTTPFound(authorization_url)

    def _get_request_session(self, request): # pylint: disable=no-self-use
        """
        Mockable method used to get the current user session.
        """
        return request.context['session']

    def get_authorization_url_and_state(self, repository_id, merge_request_id, scope):
        """
        Helper method to get the GitLab OAuth2 authorization URL and state.

        This will be used to get the user email from GitLab.

        :param repository_id: The ID of the repository this request was initiated in.
        :type repository_id: int
        :param merge_request_id: The merge request ID this request was generated in.
        :type merge_request_id: int
        :param scope: The list of OAuth2 scopes to request from GitLab.
        :type scope: [string]
        """
        client_id = cla.conf['GITLAB_CLIENT_ID']
        redirect_uri = cla.utils.get_redirect_uri('gitlab', repository_id, merge_request_id)
        return self._get_authorization_url_and_state(client_id,
                                                     redirect_uri,
                                                     scope,
                                                     cla.conf['GITLAB_OAUTH_AUTHORIZE_URL'])

    def _get_authorization_url_and_state(self, client_id, redirect_uri, scope, authorize_url): # pylint: disable=no-self-use
        """
        Mockable helper method to do the fetching of the authorization URL and state from GitLab.
        """
        return cla.utils.get_authorization_url_and_state(client_id, redirect_uri,
                                                         scope, authorize_url)

    def oauth2_redirect(self, state, code, repository_id, change_request_id, request): # pylint: disable=too-many-arguments
        """
        This is where the user will end up after having authorized the CLA system to get
        information such as email address. GitLab redirects the user here based on the redirect_uri
        we have provided in the self.get_authorization_url_state().

        It will handle storing the OAuth2 session information for this user for further requests
        and initiate the signing workflow.
        """
        cla.log.info('Handling GitLab OAuth2 redirect')
        session = self._get_request_session(request)
        if state != session.get('gitlab_oauth2_state', None):
            cla.log.warning('Invalid GitLab OAuth2 state')
            raise falcon.HTTPBadRequest('Invalid OAuth2 state', state)
        cla.log.info('Attempting to fetch OAuth2 token for state %s', state)
        client_id = cla.conf['GITLAB_CLIENT_ID']
        state = session.get('gitlab_oauth2_state')
        token_url = cla.conf['GITLAB_OAUTH_TOKEN_URL']
        client_secret = cla.conf['GITLAB_SECRET']
        redirect_uri = cla.utils.get_redirect_uri('gitlab', repository_id, change_request_id)
        token = self._fetch_token(client_id, state, token_url, client_secret, code, redirect_uri)
        cla.log.info('OAuth2 token received for state %s: %s', state, token)
        session['gitlab_oauth2_token'] = token
        repository = cla.utils.get_repository_instance()
        try:
            repository.load(repository_id)
        except DoesNotExist as err:
            cla.log.error('Repository not found during GitLab OAuth2 redirect for change ' + \
                          'request ID %s: %s', change_request_id, repository_id)
            return {'errors': {'repository_id': str(err)}}
        return self.sign_workflow(repository, change_request_id, request)

    def _fetch_token(self, client_id, state, token_url, client_secret, code, redirect_uri): # pylint: disable=too-many-arguments,no-self-use
        """
        Mockable method to fetch a OAuth2Session token.
        """
        return cla.utils.fetch_token(client_id, state, token_url, client_secret, code, redirect_uri)

    def sign_workflow(self, repository, change_request_id, request):
        """
        Once we have the 'gitlab_oauth2_token' value in the user's session, we can initiate the
        signing workflow.
        """
        gitlab_repository_id = repository.get_repository_external_id()
        cla.log.info('Initiating GitLab signing workflow for GitLab repo %s MR: %s',
                     gitlab_repository_id, change_request_id)
        user = self.get_or_create_user(request)
        signature = cla.utils.get_user_signature_by_repository(repository, user)
        if signature is not None:
            cla.utils.redirect_user_by_signature(user, signature)
        else:
            # Signature not found, create new one and send user to sign.
            cla.utils.request_individual_signature(repository, user, change_request_id)

    def process_opened_merge_request(self, data):
        """
        Helper method to handle a webhook fired from GitLab for an opened merge request.

        :param data: The data returned from GitLab on this webhook.
        :type data: dict
        """
        merge_request_id = data['object_attributes']['iid']
        gitlab_repository_id = data['object_attributes']['target_project_id']
        repository_instance = cla.utils.get_repository_instance()
        repository = repository_instance.get_repository_by_external_id(gitlab_repository_id,
                                                                       'gitlab')
        if repository is None:
            cla.log.info('Merge request fired for GitLab repository %s, but not found in ' + \
                         'CLA system', gitlab_repository_id)
            if not cla.conf['AUTO_CREATE_REPOSITORY']:
                cla.log.warning('AUTO_CREATE_REPOSITORY is set to False - please manually ' + \
                                'create this GitLab project\'s repository in the CLA system ' + \
                                'database')
                return
            else:
                cla.log.info('AUTO_CREATE_REPOSITORY is set to True, creating repository for ' + \
                             'this GitLab project automatically')
                repository = create_repository(data)
        if repository is not None:
            self.update_change_request(repository, merge_request_id)

    def get_return_url(self, repository_id, change_request_id):
        merge_request = self.get_merge_request(repository_id, change_request_id)
        return merge_request.web_url

    def update_change_request(self, repository, change_request_id):
        gitlab_repository_id = repository.get_repository_external_id()
        merge_request = self.get_merge_request(gitlab_repository_id, change_request_id)
        # Commits have the following format:
        # {'author_email': '<email-in-commit>',
        #  'author_name': 'Commit Author Name',
        #  'committer_email': '<gitlab-user-email>',
        #  'committer_name': 'GitLab User Name',
        #  'gitlab': <gitlab.Gitlab>,
        #  '_module': <module 'gitlab.v4.objects'>,
        #  '_from_api': True,
        #  'short_id': 'a7595176',
        #  'id': 'a75951768ea8ba8871f98d59fee78eec8667c684',
        #  'created_at': '2017-07-12T14:11:36.000+00:00',
        #  'title': 'Commit title',
        #  'message': 'Commit message'}
        commits = merge_request.commits()
        # Find users who have signed and who have not signed.
        signed = []
        missing = []
        for commit in commits:
            # Store the project_id in individual commits to create status messages later.
            commit.project_id = merge_request.project_id
            if commit.committer_email is not None and commit.committer_name is not None:
                # Deal with GitLab user.
                handle_commit_from_gitlab_user(repository,
                                               commit,
                                               signed,
                                               missing)
            elif commit.author_email is not None and commit.author_name is not None:
                # Deal with non-github user (just email and name in commit).
                handle_commit_from_git_author(repository,
                                              commit,
                                              signed,
                                              missing)
            else:
                # Couldn't find any author information.
                missing.append((commit.id, None))
        update_merge_request(repository.get_repository_id(),
                             merge_request,
                             signed=signed,
                             missing=missing)

    def get_merge_request(self, repository_id, merge_request_id):
        """
        Helper method to get the merge request object from GitLab.

        Merge Request Example Data:
            {'milestone': None,
             'project_id': 12,
             'target_project_id': 12,
             'iid': 2,
             'id': 2,
             'sha': '3c731aa3d9da6a60774fb110be1b71f591326800',
             'source_branch': 'test',
             'description': '',
             'assignee': None,
             'state': 'reopened',
             'title': 'Update testing',
             'labels': [],
             'user_notes_count': 0,
             'should_remove_source_branch': None,
             'gitlab': <gitlab.Gitlab>,
             'subscribed': True,
             'target_branch': 'master',
             'author': <User id:2>,
             'source_project_id': 12,
             'work_in_progress': False,
             'updated_at': '2017-07-12T21:11:28.724Z',
             'created_at': '2017-07-12T17:57:35.601Z',
             'merge_when_build_succeeds': False,
             'web_url': 'http://<domain>/<user-or-project>/<repository>/merge_requests/2',
             'merge_commit_sha': None,
             'force_remove_source_branch': False,
             'upvotes': 0,
             'downvotes': 0,
             'merge_status': 'unchecked',
             '_from_api': True,
             '_module': <module 'gitlab.v4.objects'>}

        :param repository_id: The ID of the GitLab repository.
        :type repository_id: int
        :param merge_request_id: The ID of the GitLab merge request.
        :type merge_request_id: int
        """
        cla.log.debug('Getting merge request %s from repository %s',
                      merge_request_id, repository_id)
        #try:
        return self.client.project_mergerequests.get(merge_request_id, project_id=repository_id) # pylint: disable=no-member
        #except Exception as err:
            #cla.log.error('Could not find merge request %s for repository %s: %s',
                          #merge_request_id, repository_id, str(err))

    def get_or_create_user(self, request):
        """
        Helper method to either get or create a user based on the GitLab request made.

        :param request: The hug request object for this API call.
        :type request: Request
        """
        session = self._get_request_session(request)
        gitlab_user = self.get_user_data(session, cla.conf['GITLAB_CLIENT_ID'])
        if 'error' in gitlab_user:
            # Could not get GitLab user data - maybe user revoked CLA app permissions?
            session = self._get_request_session(request)
            del session['gitlab_oauth2_state']
            del session['gitlab_oauth2_token']
            cla.log.warning('Deleted OAuth2 session data - retrying token exchange next time')
            raise falcon.HTTPError('400 Bad Request', 'gitlab_oauth2_token',
                                   'Token permissions have been rejected, please try again')
        if gitlab_user['email'] is None:
            cla.log.warning('GitLab user has no verified or public email address: %s (%s)',
                            gitlab_user['name'], gitlab_user['login'])
            raise falcon.HTTPError(
                '412 Precondition Failed', 'email',
                'Please verify and make public at least one email address with GitLab')
        cla.log.debug('Trying to load GitLab user by email: %s', gitlab_user['email'])
        # Not storing GitLab user IDs like we are with GitHub, but maybe we should?
        user = cla.utils.get_user_instance().get_user_by_email(gitlab_user['email'])
        if user is not None:
            cla.log.info('Loaded GitLab user by email: %s <%s>',
                         user.get_user_name(),
                         user.get_user_email())
            return user
        # User not found, create.
        cla.log.debug('Could not find GitLab user by email: %s', gitlab_user['email'])
        cla.log.info('Creating new GitLab user %s <%s>',
                     gitlab_user['name'],
                     gitlab_user['email'])
        user = cla.utils.get_user_instance()
        user.set_user_id(str(uuid.uuid4()))
        user.set_user_email(gitlab_user['email'])
        user.set_user_name(gitlab_user['name'])
        user.save()
        return user

    def get_user_data(self, session, client_id): # pylint: disable=no-self-use
        """
        Mockable method to get user data. Returns all GitLab user data we have on the user based on
        the current OAuth2 session.

        GitLab User Data Example:
            {'email': 'user@email.com',
             'username': 'eddie',
             'identities': [],
             'id': 2,
             'name': 'User Name',
             'state': 'active',
             'external': False,
             'location': '',
             'linkedin': '',
             'skype': '',
             'twitter': '',
             'organization': '',
             'web_url': 'http://<gitlab-domain>/<user-name>',
             'avatar_url': 'http://www.gravatar.com/avatar/<id>?s=80&d=identicon',
             'website_url': '',
             'can_create_group': True,
             'confirmed_at': '2017-02-24T04:00:57.765Z',
             'last_sign_in_at': '2017-07-11T16:04:46.654Z',
             'current_sign_in_at': '2017-07-12T15:52:43.874Z',
             'created_at': '2017-02-24T04:00:57.765Z',
             'can_create_project': True,
             'projects_limit': 50,
             'bio': '',
             'theme_id': 2,
             'color_scheme_id': 1,
             'is_admin': True,
             'two_factor_enabled': False}
        """
        token = session['gitlab_oauth2_token']
        oauth2 = OAuth2Session(client_id, token=token)
        domain = cla.conf['GITLAB_DOMAIN']
        request = oauth2.get(domain + '/api/v4/user')
        gitlab_user = request.json()
        cla.log.debug('GitLab user data: %s', gitlab_user)
        if 'error' in gitlab_user:
            cla.log.error('Could not get GitLab user data with OAuth2 token: %s',
                          gitlab_user['error'])
            return {'error': 'Could not get GitLab user data: %s' %gitlab_user['error_description']}
        return gitlab_user

    def process_reopened_merge_request(self, data):
        """
        Helper method to process a re-opened GitLab merge request.

        Simply calls the self.process_opened_merge_request() method with the data provided.

        :param data: The data provided by the GitLab webhook.
        :type data: dict:
        """
        return self.process_opened_merge_request(data)

    def process_closed_merge_request(self, data):
        """
        Helper method to process the closed GitLab merge request.

        :param data: The data provided by the GitLab webhook.
        :type data: dict:
        """
        pass

def handle_commit_from_gitlab_user(repository, commit, signed, missing):
    """
    Helper method to triage commits between signed and not-signed user signatures.

    This method deals with GitLab users found in the commit information.

    :param repository: The repository this commit belongs to.
    :type repository: cla.models.model_interfaces.Repository
    :param commit: Commit object that we're handling.
    :type commit: gitlab.ProjectCommit
    :param signed: List of information on contributors who have signed.
        Should be modified in-place to add a signer if found.
    :type signed: [dict]
    :param missing: List of information on contributors who have not yet signed.
        Should be modified in-place to add a missing signer if found.
    :type missing: [dict]
    """
    user = cla.utils.get_user_instance().get_user_by_email(commit.committer_email)
    if user is None:
        # GitLab user not in system yet, signature does not exist for this user.
        cla.log.info('GitLab user not found: %s <%s>',
                     commit.committer_name, commit.committer_email)
        missing.append((commit.id, commit.committer_name))
    else:
        cla.log.info('GitLab user found: %s <%s>',
                     user.get_user_name(), user.get_user_email())
        if cla.utils.user_signed_project_signature(user, repository):
            signed.append((commit.id, user.get_user_name()))
        else:
            missing.append((commit.id, None))

def handle_commit_from_git_author(repository, commit, signed, missing):
    """
    Helper method to triage commits between signed and not-signed user signatures.

    This method deals with non-GitLab users found in the commit information.

    :param repository: The repository this commit belongs to.
    :type repository: cla.models.model_interfaces.Repository
    :param commit: Commit object that we're handling.
    :type commit: gitlab.ProjectCommit
    :param signed: List of information on contributors who have signed.
        Should be modified in-place to add a signer if found.
    :type signed: [dict]
    :param missing: List of information on contributors who have not yet signed.
        Should be modified in-place to add a missing signer if found.
    :type missing: [dict]
    """
    user = cla.utils.get_user_instance().get_user_by_email(commit.author_email)
    if user is None:
        # Git commit author not in system yet, signature does not exist for this user.
        cla.log.info('Git commit author not found: %s <%s>',
                     commit.author_name, commit.author_email)
        missing.append((commit.id, commit.author_name))
    else:
        cla.log.info('Git commit author found: %s <%s>',
                     user.get_user_name(), user.get_user_email())
        if cla.utils.user_signed_project_signature(user, repository):
            signed.append((commit.id, user.get_user_name()))
        else:
            missing.append((commit.id, None))

def update_merge_request(repository_id, merge_request, signed, missing): # pylint: disable=too-many-locals
    """
    Helper function to update the merge request comment and status based on the list of signers.

    :param repository_id: The ID of the repository this merge request belongs to.
    :type repository_id: int
    :param merge_request: The GitLab merge request object.
    :type merge_request: gitlab.ProjectMergeRequest
    :param signed: The list of (commit, author info) tuples that have signed an
        signature for this merge request.
    :type signed: [(gitlab.ProjectCommit, dict)]
    :param missing: The list of (commit, author info) tuples that have not signed
        an signature for this merge request.
    :type missing: [(gitlab.ProjectCommit, dict)]
    """
    notification = cla.conf['GITLAB_MR_NOTIFICATION']
    both = notification == 'status+comment' or notification == 'comment+status'
    if both or notification == 'comment':
        body = cla.utils.assemble_cla_comment('gitlab', repository_id, merge_request.id,
                                              signed, missing)
        update_cla_comment(merge_request, body)
    if both or notification == 'status':
        state = 'pending'
        for commit, author in missing:
            body = cla.utils.assemble_cla_status(author, signed=False)
            sign_url = cla.utils.get_full_sign_url('gitlab', repository_id, merge_request.id)
            cla.log.info('Creating new CLA status on commit %s: %s', commit, state)
            # TODO: Need to only create new status if doesn't already exist.
            create_commit_status(merge_request, commit, state, sign_url, body)
        state = 'success'
        for commit, author in signed:
            body = cla.utils.assemble_cla_status(author, signed=True)
            sign_url = cla.utils.get_full_sign_url('gitlab', repository_id, merge_request.id)
            cla.log.info('Creating new CLA status on commit %s: %s', commit, state)
            # TODO: Need to only create new status if doesn't already exist.
            create_commit_status(merge_request, commit, state, sign_url, body)
        num_missing = len(missing)
        if num_missing > 0:
            # Need to update the last status message to prevent merging the PR.
            last_commit = merge_request.commits(all=True)[-1]
            cla.log.debug('Using last commit to set merge request status: %s', last_commit.short_id)
            signed_commits = [item[0].short_id for item in signed]
            if last_commit.short_id in signed_commits:
                num_signed = len(signed)
                total = num_signed + len(missing)
                # TODO: Need to only create new status if doesn't already exist.
                last_commit.statuses.create({'state': 'pending',
                                             'name': 'Missing CLA signatures (%s/%s)' \
                                              %(num_signed, total),
                                             'target_url': sign_url})

def update_cla_comment(merge_request, body):
    """
    Helper function to create/edit a comment on the GitLab merge request.

    :param merge_request: The merge request object in question.
    :type merge_request: gitlab.ProjectMergeRequest
    :param body: The contents of the comment.
    :type body: string
    """
    comment = get_existing_cla_comment(merge_request)
    if comment is not None:
        cla.log.info('Updating existing CLA comment for merge request: %s', merge_request.id)
        cla.log.debug('Updated CLA comment for merge request %s: %s', merge_request.id, body)
        comment.body = body
        comment.save()
    else:
        cla.log.info('Creating new CLA comment for merge request: %s', merge_request.id)
        cla.log.debug('New comment for merge request %s: %s', merge_request.id, body)
        merge_request.notes.create({'body': body})

def get_existing_cla_comment(merge_request):
    """
    Helper function to get an existing comment from the CLA system in a GitLab merge request.

    Comment Example Data:
        {'attachment': None,
         'upvote?': False,
         'downvote?': False,
         'body': 'closed',
         'noteable_type': 'MergeRequest',
         'id': 28,
         'project_id': 12,
         'merge_request_iid': 2,
         'noteable_id': 2,
         'created_at': '2017-07-12T20:01:58.366Z',
         'updated_at': '2017-07-12T20:01:58.366Z',
         'author': <User id:2>,
         'gitlab': <gitlab.Gitlab>,
         'system': True,
         '_module': <module 'gitlab.v4.objects'>,
         '_from_api': True}

    :param merge_request: The merge request object in question.
    :type merge_request: gitlab.ProjectMergeRequest
    """
    # Will only look at the first page of comments.
    comments = merge_request.notes.list()
    for comment in comments:
        if '[![CLA Check](' in comment.body:
            cla.log.info('Found matching CLA comment for merge request: %s', merge_request.id)
            return comment

def create_commit_status(merge_request, commit_hash, state, sign_url, body):
    """
    Helper function to create a pull request commit status message given the merge request
    commit hash.

    :param merge_request: The GitLab Merge Request object.
    :type merge_request: gitlab.ProjectMergeRequest
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
    for commit in merge_request.commits(all=True):
        if commit.id == commit_hash:
            commit_obj = commit
            break
    if commit_obj is None:
        cla.log.error('Could not post status on merge request %s: Commit %s not found',
                      merge_request.id, commit_hash)
        return
    commit_obj.project_id = merge_request.project_id
    commit_obj.statuses.create({'state': state, 'name': body, 'target_url': sign_url})

def create_repository(data):
    """
    Helper method to create a repository object in the CLA database given merge request data.

    :param data: The data provided by the GitLab webhook.
    :type data: dict
    :return: The newly created repository object - already in the DB.
    :rtype: cla.models.model_interfaces.Repository
    """
    try:
        repository = cla.utils.get_repository_instance()
        repository.set_repository_id(str(uuid.uuid4()))
        # TODO: Need to use an ID unique across all repository providers instead of namespace.
        repository.set_repository_project_id(data['project']['namespace'])
        # GitLab Project == Repository as far as CLA system is concerned.
        repository.set_repository_external_id(data['object_attributes']['target_project_id'])
        repository.set_repository_name(data['project']['name'])
        repository.set_repository_type('gitlab')
        repository.set_repository_url(data['project']['web_url'])
        repository.save()
        return repository
    except Exception as err:
        cla.log.error('Could not create GitLab repository automatically: %s', str(err))
        return None

class MockGitLab(GitLab):
    """
    The GitLab repository service mock class for testing.
    """
    def __init__(self, oauth2_token=False):
        super().__init__()
        self.oauth2_token = oauth2_token

    def _get_gitlab_client(self, domain, token):
        return MockGitLabClient(domain, token)

    def _get_authorization_url_and_state(self, client_id, redirect_uri, scope, authorize_url):
        authorization_url = 'http://authorization.url'
        state = 'random-state-here'
        return authorization_url, state

    def _fetch_token(self, client_id, state, token_url, client_secret, code, redirect_uri): # pylint: disable=too-many-arguments
        return 'random-token'

    def _get_request_session(self, request):
        if self.oauth2_token:
            return {'gitlab_oauth2_token': 'random-token',
                    'gitlab_oauth2_state': 'random-state'}
        return {}

    def get_user_data(self, session, client_id):
        return {'email': 'test@user.com', 'name': 'Test User', 'id': 123}

class MockGitLabClient(object): # pylint: disable=too-few-public-methods
    """
    The GitLab Client object mock class for testing.
    """
    def __init__(self, domain, token):
        self.domain = domain
        self.token = token
        self.project_mergerequests = MockGitLabMergeRequest(None, None)

class MockGitLabMergeRequest(object):
    """Mock GitLab merge request object."""
    def __init__(self, merge_request_id, project_id=None):
        self.id = 1
        self.web_url = 'http://gitlab.com/test/1'
        self.notes = MockGitLabNotes()
        self.commits = MockGitLabCommits()
        self.merge_request_id = merge_request_id
        self.project_id = project_id

    def get(self, merge_request_id, project_id=None): # pylint: disable=no-self-use
        """Mock merge request get method."""
        return MockGitLabMergeRequest(merge_request_id, project_id)

class MockGitLabCommits(object): # pylint: disable=too-few-public-methods
    """Mock GitLab commits object."""
    def __iter__(self):
        yield MockGitLabCommit()

    def __call__(self, all=False): # pylint: disable=redefined-builtin
        return [MockGitLabCommit()]

    def get(self, sha): # pylint: disable=unused-argument,no-self-use
        """Mock GitLab commits get method."""
        return MockGitLabCommit()

class MockGitLabCommit(object): # pylint: disable=too-few-public-methods
    """Mock GitLab commit object."""
    def __init__(self):
        self.id = 'sha-here'
        self.short_id = self.id
        self.project_id = 1
        self.committer_name = 'GitLab User'
        self.committer_email = 'test@gitlab.com'
        self.statuses = MockGitLabStatuses()

class MockGitLabNotes(object):
    """Mock GitLab notes object."""
    def __init__(self):
        pass

    def list(self): # pylint: disable=no-self-use
        """Mock GitLab list of notes."""
        return []

    def create(self, data):
        """Mock GitLab note creation method."""
        pass

class MockGitLabStatuses(object): # pylint: disable=too-few-public-methods
    """Mock GitLab list of status objects."""
    def create(self, data):
        """Mock GitLab statuses create method."""
        pass
