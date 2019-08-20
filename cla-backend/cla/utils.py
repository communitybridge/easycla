# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Utility functions for the CLA project.
"""

import inspect
import json
import os
import urllib.parse

import falcon
from hug.middleware import SessionMiddleware
from requests_oauthlib import OAuth2Session

import cla
from cla.models import DoesNotExist
from cla.models.dynamo_models import Repository, GitHubOrg

api_base_url = os.environ.get('CLA_API_BASE', '')
cla_logo_url = os.environ.get('CLA_BUCKET_LOGO_URL', '')


def get_cla_path():
    """Returns the CLA code root directory on the current system."""
    cla_folder_dir = os.path.dirname(os.path.abspath(inspect.getfile(inspect.currentframe())))
    cla_root_dir = os.path.dirname(cla_folder_dir)
    return cla_root_dir


def get_session_middleware():
    """Prepares the hug middleware to manage key-value session data."""
    store = get_key_value_store_service()
    return SessionMiddleware(store, context_name='session', cookie_name='cla-sid',
                             cookie_max_age=300, cookie_domain=None, cookie_path='/',
                             cookie_secure=False)


def create_database(conf=None):
    """
    Helper function to create the CLA database. Will utilize the appropriate database
    provider based on configuration.

    :param conf: Configuration dictionary/object - typically parsed from the CLA config file.
    :type conf: dict
    """
    if conf is None:
        conf = cla.conf
    cla.log.info('Creating CLA database in %s', conf['DATABASE'])
    if conf['DATABASE'] == 'DynamoDB':
        from cla.models.dynamo_models import create_database as cd
    else:
        raise Exception('Invalid database selection in configuration: %s' % conf['DATABASE'])
    cd()


def delete_database(conf=None):
    """
    Helper function to delete the CLA database. Will utilize the appropriate database
    provider based on configuration.

    :WARNING: Use with caution.

    :param conf: Configuration dictionary/object - typically parsed from the CLA config file.
    :type conf: dict
    """
    if conf is None:
        conf = cla.conf
    cla.log.warning('Deleting CLA database in %s', conf['DATABASE'])
    if conf['DATABASE'] == 'DynamoDB':
        from cla.models.dynamo_models import delete_database as dd
    else:
        raise Exception('Invalid database selection in configuration: %s' % conf['DATABASE'])
    dd()


def get_database_models(conf=None):
    """
    Returns the database models based on the configuration dict provided.

    :param conf: Configuration dictionary/object - typically parsed from the CLA config file.
    :type conf: dict
    :return: Dictionary of all the supported database object classes (User, Signature, Repository,
        company, Project, Document) - keyed by name:

            {'User': cla.models.model_interfaces.User,
             'Signature': cla.models.model_interfaces.Signature,...}

    :rtype: dict
    """
    if conf is None:
        conf = cla.conf
    if conf['DATABASE'] == 'DynamoDB':
        from cla.models.dynamo_models import User, Signature, Repository, \
            Company, Project, Document, \
            GitHubOrg, Gerrit, UserPermissions
        return {'User': User, 'Signature': Signature, 'Repository': Repository,
                'Company': Company, 'Project': Project, 'Document': Document,
                'GitHubOrg': GitHubOrg, 'Gerrit': Gerrit, 'UserPermissions': UserPermissions}
    else:
        raise Exception('Invalid database selection in configuration: %s' % conf['DATABASE'])


def get_user_instance(conf=None):
    """
    Helper function to get a database User model instance based on CLA configuration.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: A User model instance based on configuration specified.
    :rtype: cla.models.model_interfaces.User
    """
    return get_database_models(conf)['User']()


def get_signature_instance(conf=None):
    """
    Helper function to get a database Signature model instance based on CLA configuration.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: An Signature model instance based on configuration.
    :rtype: cla.models.model_interfaces.Signature
    """
    return get_database_models(conf)['Signature']()


def get_repository_instance(conf=None):
    """
    Helper function to get a database Repository model instance based on CLA configuration.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: A Repository model instance based on configuration specified.
    :rtype: cla.models.model_interfaces.Repository
    """
    return get_database_models(conf)['Repository']()


def get_github_organization_instance(conf=None):
    """
    Helper function to get a database GitHubOrg model instance based on CLA configuration.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: A Repository model instance based on configuration specified.
    :rtype: cla.models.model_interfaces.GitHubOrg
    """
    return get_database_models(conf)['GitHubOrg']()


def get_gerrit_instance(conf=None):
    """
    Helper function to get a database Gerrit model based on CLA configuration.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: A Gerrit model instance based on configuration specified.
    :rtype: cla.models.model_interfaces.Gerrit
    """
    return get_database_models(conf)['Gerrit']()


def get_company_instance(conf=None):
    """
    Helper function to get a database company model instance based on CLA configuration.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: A company model instance based on configuration specified.
    :rtype: cla.models.model_interfaces.Company
    """
    return get_database_models(conf)['Company']()


def get_project_instance(conf=None):
    """
    Helper function to get a database Project model instance based on CLA configuration.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: A Project model instance based on configuration specified.
    :rtype: cla.models.model_interfaces.Project
    """
    return get_database_models(conf)['Project']()


def get_document_instance(conf=None):
    """
    Helper function to get a database Document model instance based on CLA configuration.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: A Document model instance based on configuration specified.
    :rtype: cla.models.model_interfaces.Document
    """
    return get_database_models(conf)['Document']()


def get_email_service(conf=None, initialize=True):
    """
    Helper function to get the configured email service instance.

    :param conf: Same as get_database_models().
    :type conf: dict
    :param initialize: Whether or not to run the initialize method on the instance.
    :type initialize: boolean
    :return: The email service model instance based on configuration specified.
    :rtype: EmailService
    """
    if conf is None:
        conf = cla.conf
    email_service = conf['EMAIL_SERVICE']
    if email_service == 'SMTP':
        from cla.models.smtp_models import SMTP as email
    elif email_service == 'MockSMTP':
        from cla.models.smtp_models import MockSMTP as email
    elif email_service == 'SES':
        from cla.models.ses_models import SES as email
    elif email_service == 'MockSES':
        from cla.models.ses_models import MockSES as email
    else:
        raise Exception('Invalid email service selected in configuration: %s' % email_service)
    email_instance = email()
    if initialize:
        email_instance.initialize(conf)
    return email_instance


def get_signing_service(conf=None, initialize=True):
    """
    Helper function to get the configured signing service instance.

    :param conf: Same as get_database_models().
    :type conf: dict
    :param initialize: Whether or not to run the initialize method on the instance.
    :type initialize: boolean
    :return: The signing service instance based on configuration specified.
    :rtype: SigningService
    """
    if conf is None:
        conf = cla.conf
    signing_service = conf['SIGNING_SERVICE']
    if signing_service == 'DocuSign':
        from cla.models.docusign_models import DocuSign as signing
    elif signing_service == 'MockDocuSign':
        from cla.models.docusign_models import MockDocuSign as signing
    else:
        raise Exception('Invalid signing service selected in configuration: %s' % signing_service)
    signing_service_instance = signing()
    if initialize:
        signing_service_instance.initialize(conf)
    return signing_service_instance


def get_storage_service(conf=None, initialize=True):
    """
    Helper function to get the configured storage service instance.

    :param conf: Same as get_database_models().
    :type conf: dict
    :param initialize: Whether or not to run the initialize method on the instance.
    :type initialize: boolean
    :return: The storage service instance based on configuration specified.
    :rtype: StorageService
    """
    if conf is None:
        conf = cla.conf
    storage_service = conf['STORAGE_SERVICE']
    if storage_service == 'LocalStorage':
        from cla.models.local_storage import LocalStorage as storage
    elif storage_service == 'S3Storage':
        from cla.models.s3_storage import S3Storage as storage
    elif storage_service == 'MockS3Storage':
        from cla.models.s3_storage import MockS3Storage as storage
    else:
        raise Exception('Invalid storage service selected in configuration: %s' % storage_service)
    storage_instance = storage()
    if initialize:
        storage_instance.initialize(conf)
    return storage_instance


def get_pdf_service(conf=None, initialize=True):
    """
    Helper function to get the configured PDF service instance.

    :param conf: Same as get_database_models().
    :type conf: dict
    :param initialize: Whether or not to run the initialize method on the instance.
    :type initialize: boolean
    :return: The PDF service instance based on configuration specified.
    :rtype: PDFService
    """
    if conf is None:
        conf = cla.conf
    pdf_service = conf['PDF_SERVICE']
    if pdf_service == 'DocRaptor':
        from cla.models.docraptor_models import DocRaptor as pdf
    elif pdf_service == 'MockDocRaptor':
        from cla.models.docraptor_models import MockDocRaptor as pdf
    else:
        raise Exception('Invalid PDF service selected in configuration: %s' % pdf_service)
    pdf_instance = pdf()
    if initialize:
        pdf_instance.initialize(conf)
    return pdf_instance


def get_key_value_store_service(conf=None):
    """
    Helper function to get the configured key-value store service instance.

    :param conf: Same as get_database_models().
    :type conf: dict
    :return: The key-value store service instance based on configuration specified.
    :rtype: KeyValueStore
    """
    if conf is None:
        conf = cla.conf
    keyvalue = cla.conf['KEYVALUE']
    if keyvalue == 'Memory':
        from hug.store import InMemoryStore as Store
    elif keyvalue == 'DynamoDB':
        from cla.models.dynamo_models import Store
    else:
        raise Exception('Invalid key-value store selected in configuration: %s' % keyvalue)
    return Store()


def get_supported_repository_providers():
    """
    Returns a dict of supported repository service providers.

    :return: Dictionary of supported repository service providers in the following
        format: {'<provider_name>': <provider_class>}
    :rtype: dict
    """
    from cla.models.github_models import GitHub, MockGitHub
    # from cla.models.gitlab_models import GitLab, MockGitLab
    # return {'github': GitHub, 'mock_github': MockGitHub,
    # 'gitlab': GitLab, 'mock_gitlab': MockGitLab}
    return {'github': GitHub, 'mock_github': MockGitHub}


def get_repository_service(provider, initialize=True):
    """
    Get a repository service instance by provider name.

    :param provider: The provider to load.
    :type provider: string
    :param initialize: Whether or not to call the initialize() method on the object.
    :type initialize: boolean
    :return: A repository provider instance (GitHub, Gerrit, etc).
    :rtype: RepositoryService
    """
    providers = get_supported_repository_providers()
    if provider not in providers:
        raise NotImplementedError('Provider not supported')
    instance = providers[provider]()
    if initialize:
        instance.initialize(cla.conf)
    return instance


def get_repository_service_by_repository(repository, initialize=True):
    """
    Helper function to get a repository service provider instance based
    on a repository.

    :param repository: The repository object or repository_id.
    :type repository: cla.models.model_interfaces.Repository | string
    :param initialize: Whether or not to call the initialize() method on the object.
    :type initialize: boolean
    :return: A repository provider instance (GitHub, Gerrit, etc).
    :rtype: RepositoryService
    """
    repository_model = get_database_models()['Repository']
    if isinstance(repository, repository_model):
        repo = repository
    else:
        repo = repository_model()
        repo.load(repository)
    provider = repo.get_repository_type()
    return get_repository_service(provider, initialize)


def get_supported_document_content_types():  # pylint: disable=invalid-name
    """
    Returns a list of supported document content types.

    :return: List of supported document content types.
    :rtype: dict
    """
    return ['pdf', 'url+pdf', 'storage+pdf']


def get_project_document(project, document_type, major_version, minor_version):
    """
    Helper function to get the specified document from a project.

    :param project: The project model object to look in.
    :type project: cla.models.model_interfaces.Project
    :param document_type: The type of document (individual or corporate).
    :type document_type: string
    :param major_version: The major version number to look for.
    :type major_version: integer
    :param minor_version: The minor version number to look for.
    :type minor_version: integer
    :return: The document model if found.
    :rtype: cla.models.model_interfaces.Document
    """
    if document_type == 'individual':
        documents = project.get_project_individual_documents()
    else:
        documents = project.get_project_corporate_documents()
    for document in documents:
        if document.get_document_major_version() == major_version and \
                document.get_document_minor_version() == minor_version:
            return document
    return None


def get_project_latest_individual_document(project_id):
    """
    Helper function to return the latest individual document belonging to a project.

    :param project_id: The project ID in question.
    :type project_id: string
    :return: Latest ICLA document object for this project.
    :rtype: cla.models.model_instances.Document
    """
    project = get_project_instance()
    project.load(str(project_id))
    document_models = project.get_project_individual_documents()
    major, minor = get_last_version(document_models)
    return project.get_project_individual_document(major, minor)


# TODO Heller remove
def get_project_latest_corporate_document(project_id):
    """
    Helper function to return the latest corporate document belonging to a project.

    :param project_id: The project ID in question.
    :type project_id: string
    :return: Latest CCLA document object for this project.
    :rtype: cla.models.model_instances.Document
    """
    project = get_project_instance()
    project.load(str(project_id))
    document_models = project.get_project_corporate_documents()
    major, minor = get_last_version(document_models)
    return project.get_project_corporate_document(major, minor)


def get_last_version(documents):
    """
    Helper function to get the last version of the list of documents provided.

    :param documents: List of documents to check.
    :type documents: [cla.models.model_interfaces.Document]
    :return: 2-item tuple containing (major, minor) version number.
    :rtype: tuple
    """
    last_major = 0  # 0 will be returned if no document was found.
    last_minor = -1  # -1 will be returned if no document was found.
    for document in documents:
        current_major = document.get_document_major_version()
        current_minor = document.get_document_minor_version()
        if current_major > last_major:
            last_major = current_major
            last_minor = current_minor
            continue
        if current_major == last_major and current_minor > last_minor:
            last_minor = current_minor
    return (last_major, last_minor)


def user_signed_project_signature(user, project_id, latest_major_version=True):
    """
    Helper function to check if a user has signed a project signature tied to a repository.
    Will consider both ICLA and employee signatures.

    :param user: The user object to check for.
    :type user: cla.models.model_interfaces.User
    :param project_id: The project to check for.
    :type project_id: string
    :param latest_major_version: True means only the latest document major version will be considered.
    :type latest_major_version: bool
    :return: Whether or not the user has an signature that's signed and approved
        for this project.
    :rtype: boolean
    """

    # Check if we have an ICLA for this user
    cla.log.debug('checking to see if user has signed an ICLA, user: {}, project: {}'.
                  format(user, project_id))

    signature = user.get_latest_signature(project_id)
    if signature is not None:
        cla.log.debug('ICLA signature found for user: {} on project: {}, signature_id: {}'.
                      format(user.get_user_id(), project_id, signature.get_signature_id()))

        if latest_major_version:  # Ensure it's latest signature.
            project = get_project_instance()
            project.load(str(project_id))
            document_models = project.get_project_individual_documents()
            major, _ = get_last_version(document_models)
            if signature.get_signature_document_major_version() != major:
                cla.log.debug('User: {} only has an old document version signed (v{}) - needs a new version'.
                              format(user, signature.get_signature_document_major_version()))
                return False

        if signature.get_signature_signed() and signature.get_signature_approved():
            # Signature found and signed/approved.
            cla.log.debug('User: {} has ICLA signed and approved signature_id: {} '
                          'for project: {}'.
                          format(user, signature.get_signature_id(), project_id))
            return True
        elif signature.get_signature_signed():  # Not approved yet.
            cla.log.debug('User: {} has ICLA signed with signature_id: {}, project: {}, '
                          'but has not been approved yet'.
                          format(user, signature.get_signature_id(), project_id))
            return False
        else:  # Not signed or approved yet.
            cla.log.debug('User: {} has ICLA with signature_id: {}, project: {}, '
                          'but has not been signed or approved yet'.
                          format(user, signature.get_signature_id(), project_id))
            return False
    else:
        cla.log.debug('ICLA signature NOT found for User: {} on project: {}'.
                      format(user, project_id))

    # Check if we have an CCLA for this user
    company_id = user.get_user_company_id()
    cla.log.debug('checking to see if user has signed an CCLA, user: {}, project_id: {}, company_id: {}'.
                  format(user, project_id, company_id))

    if company_id is not None:
        signature = user.get_latest_signature(project_id, company_id=company_id)

        # Don't check the version for employee signatures.
        if signature is not None:
            cla.log.debug('CCLA signature found for user: {} on project: {}, signature_id: {}'.
                          format(user, project_id, signature.get_signature_id()))

            if signature.get_signature_signed and signature.get_signature_approved:
                cla.log.debug('User: {} has a signed and approved CCLA for project: {}'.
                              format(user, project_id))
                return True

            if signature.get_signature_signed():
                cla.log.debug('User: {} has CCLA signed with signature_id: {}, project: {}, '
                              'but has not been approved yet'.
                              format(user, signature.get_signature_id(), project_id))
                return False
            else:  # Not signed or approved yet.
                cla.log.debug('User: {} has CCLA with signature_id: {}, project: {}, '
                              'but has not been signed or approved yet'.
                              format(user, signature.get_signature_id(), project_id))
                return False
    else:
        cla.log.debug('User: {} is NOT associated with a company - unable to check for a CCLA.'.format(user))
        return False

    cla.log.debug('User: {} has not signed an ICLA or CCLA'.format(user))
    return False


def get_redirect_uri(repository_service, installation_id, github_repository_id, change_request_id):
    """
    Function to generate the redirect_uri parameter for a repository service's OAuth2 process.

    :param repository_service: The repository service provider we're currently initiating the
        OAuth2 process with. Currently only supports 'github' and 'gitlab'.
    :type repository_service: string
    :param installation_id: The EasyCLA GitHub application ID
    :type installation_id: string
    :param github_repository_id: The ID of the repository object that applies for this OAuth2 process.
    :type github_repository_id: string
    :param change_request_id: The ID of the change request in question. Is a PR number if
        repository_service is 'github'. Is a merge request if the repository_service is 'gitlab'.
    :type change_request_id: string
    :return: The redirect_uri parameter expected by the OAuth2 process.
    :rtype: string
    """
    params = {'installation_id': installation_id,
              'github_repository_id': github_repository_id,
              'change_request_id': change_request_id}
    params = urllib.parse.urlencode(params)
    return '{}/v2/repository-provider/{}/oauth2_redirect?{}'.format(cla.conf['API_BASE_URL'], repository_service,
                                                                    params)


def get_full_sign_url(repository_service, installation_id, github_repository_id, change_request_id):
    """
    Helper function to get the full sign URL that the user should click to initiate the signing
    workflow.

    :TODO: Update comments.

    :param repository_service: The repository service provider we're getting the sign url for.
        Should be one of the supported repository providers ('github', 'gitlab', etc).
    :type repository_service: string
    :param installation_id: The EasyCLA GitHub application ID
    :type installation_id: string
    :param github_repository_id: The ID of the repository for this signature (used in order to figure out
        where to send the user once signing is complete.
    :type github_repository_id: int
    :param change_request_id: The change request ID for this signature (used in order to figure out
        where to send the user once signing is complete. Should be a pull request number when
        repository_service is 'github'. Should be a merge request ID when repository_service is
        'gitlab'.
    :type change_request_id: int
    """
    return '{}/v2/repository-provider/{}/sign/{}/{}/{}'.format(cla.conf['API_BASE_URL'], repository_service,
                                                               str(installation_id), str(github_repository_id),
                                                               str(change_request_id))


def get_comment_badge(repository_type, all_signed, sign_url):
    """
    Returns the CLA badge that will appear on the change request comment (PR for 'github', merge
    request for 'gitlab', etc)

    :param repository_type: The repository service provider we're getting the badge for.
        Should be one of the supported repository providers ('github', 'gitlab', etc).
    :type repository_type: string
    :param all_signed: Whether or not all committers have signed the change request.
    :type all_signed: boolean
    :param sign_url: The URL for the user to click in order to initiate signing.
    :type sign_url: string
    """

    if all_signed:
        badge_url = '{}/cla-signed.png'.format(cla_logo_url)
        badge_hyperlink = 'https://lfcla.com'
    else:
        badge_url = '{}/cla-notsigned.png'.format(cla_logo_url)
        badge_hyperlink = sign_url
    return '[![CLA Check](' + badge_url + ')](' + badge_hyperlink + ')'


def assemble_cla_status(author_name, signed=False):
    """
    Helper function to return the text that will display on a change request status.

    For GitLab there isn't much space here - we rely on the user hovering their mouse over the icon.
    For GitHub there is a 140 character limit.

    :param author_name: The name of the author of this commit.
    :type author_name: string
    :param signed: Whether or not the author has signed an signature.
    :type signed: boolean
    """
    if author_name is None:
        author_name = 'Unknown'
    if signed:
        return (author_name, 'EasyCLA check passed. You are authorized to contribute.')
    return (author_name, 'Missing CLA Authorization.')


def assemble_cla_comment(repository_type, installation_id, github_repository_id, change_request_id, signed, missing):
    """
    Helper function to generate a CLA comment based on a a change request.

    :TODO: Update comments

    :param repository_type: The type of repository this comment will be posted on ('github',
        'gitlab', etc).
    :type repository_type: string
    :param installation_id: The EasyCLA GitHub application ID
    :type installation_id: string
    :param github_repository_id: The ID of the repository for this change request.
    :type github_repository_id: int
    :param change_request_id: The repository service's ID of this change request.
    :type change_request_id: id
    :param signed: The list of commit hashes and authors that have signed an signature for this
        change request.
    :type signed: [(string, string)]
    :param missing: The list of commit hashes and authors that have not signed for this
        change request.
    :type missing: [(string, string)]
    """
    num_missing = len(missing)
    sign_url = get_full_sign_url(repository_type, installation_id, github_repository_id, change_request_id)
    comment = get_comment_body(repository_type, sign_url, signed, missing)
    all_signed = num_missing == 0
    badge = get_comment_badge(repository_type, all_signed, sign_url)
    return badge + '<br />' + comment


def get_comment_body(repository_type, sign_url, signed, missing):
    """
    Returns the CLA comment that will appear on the repository provider's change request item.

    :param repository_type: The repository type where this comment will be posted ('github',
        'gitlab', etc).
    :type repository_type: string
    :param sign_url: The URL for the user to click in order to initiate signing.
    :type sign_url: string
    :param signed: List of tuples containing the commit and author name of signers.
    :type signed: [(string, string)]
    :param missing: List of tuples containing the commit and author name of not-signed users.
    :type missing: [(string, string)]
    """
    cla.log.info('Getting comment body for repository type: %s', repository_type)
    failed = ':x:'
    success = ':white_check_mark:'
    committers_comment = ''
    num_signed = len(signed)
    num_missing = len(missing)
    if num_signed > 0:
        # Group commits by author.
        committers = {}
        for commit, author in signed:
            if author is None:
                author = 'Unknown'
            if author not in committers:
                committers[author] = []
            committers[author].append(commit)
        # Print author commit information.
        committers_comment += '<ul>'
        for author, commit_hashes in committers.items():
            committers_comment += '<li>' + success + '  ' + author + \
                                  ' (' + ", ".join(commit_hashes) + ')</li>'
        committers_comment += '</ul>'
    if num_missing > 0:
        text = 'One or more committers are not authorized under a signed CLA as indicated below. ' + \
               '[Please click here to be authorized](' + sign_url + ').'
        # Group commits by author.
        committers = {}
        for commit, author in missing:
            if author is None:
                author = 'Unknown'
            if author not in committers:
                committers[author] = []
            committers[author].append(commit)
        # Print author commit information.
        committers_comment += '<ul>'
        for author, commit_hashes in committers.items():
            committers_comment += '<li>[' + failed + '](' + sign_url + ')  ' + \
                                  author + ' (' + ", ".join(commit_hashes) + ')</li>'
        committers_comment += '</ul>'
        return text + committers_comment
    text = 'The committers are authorized under a signed CLA.'
    return text + committers_comment


def get_authorization_url_and_state(client_id, redirect_uri, scope, authorize_url):
    """
    Helper function to get an OAuth2 session authorization URL and state.

    :param client_id: The client ID for this OAuth2 session.
    :type client_id: string
    :param redirect_uri: The redirect URI to specify in this OAuth2 session.
    :type redirect_uri: string
    :param scope: The list of scope items to use for this OAuth2 session.
    :type scope: [string]
    :param authorize_url: The URL to submit the OAuth2 request.
    :type authorize_url: string
    """
    oauth = OAuth2Session(client_id, redirect_uri=redirect_uri, scope=scope)
    authorization_url, state = oauth.authorization_url(authorize_url)
    cla.log.debug(
        'utils.py - get_authorization_url_and_state - authorization_url: {}, state: {}'.format(authorization_url,
                                                                                               state))
    return authorization_url, state


def fetch_token(client_id, state, token_url, client_secret, code,
                redirect_uri=None):  # pylint: disable=too-many-arguments
    """
    Helper function to fetch a OAuth2 session token.

    :param client_id: The client ID for this OAuth2 session.
    :type client_id: string
    :param state: The OAuth2 session state.
    :type state: string
    :param token_url: The token URL for this OAuth2 session.
    :type token_url: string
    :param client_secret: the client secret
    :type client_secret: string
    :param code: The OAuth2 session code.
    :type code: string
    :param redirect_uri: The redirect URI for this OAuth2 session.
    :type redirect_uri: string
    """
    if redirect_uri is not None:
        oauth2 = OAuth2Session(client_id, state=state, scope=['user:email'], redirect_uri=redirect_uri)
    else:
        oauth2 = OAuth2Session(client_id, state=state, scope=['user:email'])
    cla.log.debug('utils.py - oauth2.fetch_token - token_url: {}, client_id: {}, client_secret: {}, code: {}'.
                  format(token_url, client_id, client_secret, code))
    return oauth2.fetch_token(token_url, client_secret=client_secret, code=code)


def redirect_user_by_signature(user, signature):
    """
    Helper method to redirect a user based on their signature status and return_url.

    :param user: The user object for this redirect.
    :type user: cla.models.model_interfaces.User
    :param signature: The signature object for this user.
    :type signature: cla.models.model_interfaces.Signature
    """
    return_url = signature.get_signature_return_url()
    if signature.get_signature_signed() and signature.get_signature_approved():
        # Signature already signed and approved.
        # TODO: Notify user of signed and approved signature somehow.
        cla.log.info('Signature already signed and approved for user: %s, %s',
                     user.get_user_emails(), signature.get_signature_id())
        if return_url is None:
            cla.log.info('No return_url set in signature object - serving success message')
            return {'status': 'signed and approved'}
        else:
            cla.log.info('Redirecting user back to %s', return_url)
            raise falcon.HTTPFound(return_url)
    elif signature.get_signature_signed():
        # Awaiting approval.
        # TODO: Notify user of pending approval somehow.
        cla.log.info('Signature signed but not approved yet: %s',
                     signature.get_signature_id())
        if return_url is None:
            cla.log.info('No return_url set in signature object - serving pending message')
            return {'status': 'pending approval'}
        else:
            cla.log.info('Redirecting user back to %s', return_url)
            raise falcon.HTTPFound(return_url)
    else:
        # Signature awaiting signature.
        sign_url = signature.get_signature_sign_url()
        signature_id = signature.get_signature_id()
        cla.log.info('Signature exists, sending user to sign: %s (%s)', signature_id, sign_url)
        raise falcon.HTTPFound(sign_url)


def get_active_signature_metadata(user_id):
    """
    When a user initiates the signing process, the CLA system must store information on this
    signature - such as where the user came from, what repository it was initiated on, etc.
    This information is temporary while the signature is in progress. See the Signature object
    for information on this signature once the signing is complete.

    :param user_id: The ID of the user in question.
    :type user_id: string
    :return: Dict of data on the signature request from this user.
    :rtype: dict
    """
    store = get_key_value_store_service()
    key = 'active_signature:' + str(user_id)
    if store.exists(key):
        return json.loads(store.get(key))
    return None


def set_active_signature_metadata(user_id, project_id, repository_id, pull_request_id):
    """
    When a user initiates the signing process, the CLA system must store information on this
    signature - such as where the user came from, what repository it was initiated on, etc.
    This is a helper function to perform the storage of this information.

    :param user_id: The ID of the user beginning the signing process.
    :type user_id: string
    :param project_id: The ID of the project this signature is for.
    :type project_id: string
    :param repository_id: The repository where the signature is coming from.
    :type repository_id: string
    :param pull_request_id: The PR where this signature request is coming from (where the user
        clicked on the 'Sign CLA' badge).
    :type pull_request_id: string
    """
    store = get_key_value_store_service()
    key = 'active_signature:' + str(user_id)  # Should have been set when user initiated the signature.
    value = json.dumps({'user_id': user_id,
                        'project_id': project_id,
                        'repository_id': repository_id,
                        'pull_request_id': pull_request_id})
    store.set(key, value)
    cla.log.info('Stored active signature details for user %s: Key - %s  Value - %s', user_id, key, value)


def delete_active_signature_metadata(user_id):
    """
    Helper function to delete all metadata regarding the active signature request for the user.

    :param user_id: The ID of the user in question.
    :type user_id: string
    """
    store = get_key_value_store_service()
    key = 'active_signature:' + str(user_id)
    store.delete(key)
    cla.log.info('Deleted stored active signature details for user %s', user_id)


def get_active_signature_return_url(user_id, metadata=None):
    """
    Helper function to get a user's active signature return URL.

    :param user_id: The user ID in question.
    :type user_id: string
    :return: The URL the user will be redirected to upon successful signature.
    :rtype: string
    """
    if metadata is None:
        metadata = get_active_signature_metadata(user_id)
    if metadata is None:
        cla.log.error('Could not find active signature for user %s, return URL request failed' % user_id)
        return None

    # Get Github ID from metadata
    github_repository_id = metadata['repository_id']

    # Get installation id through a helper function
    installation_id = get_installation_id_from_github_repository(github_repository_id)
    if installation_id is None:
        cla.log.error('Could not find installation ID that is configured for this repository ID: %s',
                      github_repository_id)
        return None

    github = cla.utils.get_repository_service('github')
    return github.get_return_url(metadata['repository_id'],
                                 metadata['pull_request_id'],
                                 installation_id)


def get_installation_id_from_github_repository(github_repository_id):
    # Get repository ID that references the github ID. 
    try:
        repository = Repository().get_repository_by_external_id(github_repository_id, 'github')
    except DoesNotExist:
        return None

    # Get Organization from this repository
    organization = GitHubOrg()
    try:
        organization.load(repository.get_repository_organization_name())
    except DoesNotExist:
        return None

    # Get this organization's installation ID 
    return organization.get_organization_installation_id()


def get_project_id_from_github_repository(github_repository_id):
    # Get repository ID that references the github ID. 
    try:
        repository = Repository().get_repository_by_external_id(github_repository_id, 'github')
    except DoesNotExist:
        return None

    # Get project ID (contract group ID) of this repository
    return repository.get_repository_project_id()


def get_individual_signature_callback_url(user_id, metadata=None):
    """
    Helper function to get a user's active signature callback URL.

    :param user_id: The user ID in question.
    :type user_id: string
    :return: The callback URL that will be hit by the signing service provider.
    :rtype: string
    """
    if metadata is None:
        metadata = get_active_signature_metadata(user_id)
    if metadata is None:
        cla.log.error('Could not find active signature for user %s, callback URL request failed' % user_id)
        return None

    # Get Github ID from metadata
    github_repository_id = metadata['repository_id']

    # Get installation id through a helper function
    installation_id = get_installation_id_from_github_repository(github_repository_id)
    if installation_id is None:
        cla.log.error('Could not find installation ID that is configured for this repository ID: %s',
                      github_repository_id)
        return None

    return os.path.join(api_base_url, 'v2/signed/individual', str(installation_id), str(metadata['repository_id']),
                        str(metadata['pull_request_id']))


def request_individual_signature(installation_id, github_repository_id, user, change_request_id, callback_url=None):
    """
    Helper function send the user off to sign an signature based on the repository.

    :TODO: Update comments.

    :param repository: The repository object in question.
    :type repository: cla.models.model_interfaces.Repository
    :param user: The user in question.
    :type user: cla.models.model_interfaces.User
    :param change_request_id: The change request ID (used to redirect the user after signing).
    :type change_request_id: string
    :param callback_url: Optionally provided a callback_url. Will default to
        <SIGNED_CALLBACK_URL>/<repo_id>/<change_request_id>.
    :type callback_url: string
    """
    project_id = get_project_id_from_github_repository(github_repository_id)
    repo_service = get_repository_service('github')
    return_url = repo_service.get_return_url(github_repository_id,
                                             change_request_id,
                                             installation_id)
    if callback_url is None:
        callback_url = os.path.join(api_base_url, 'v2/signed/individual', str(installation_id), str(change_request_id))

    signing_service = get_signing_service()
    return_url_type = 'Github'
    signature_data = signing_service.request_individual_signature(project_id,
                                                                  user.get_user_id(),
                                                                  return_url_type,
                                                                  return_url,
                                                                  callback_url)
    if 'sign_url' in signature_data:
        raise falcon.HTTPFound(signature_data['sign_url'])
    cla.log.error('Could not get sign_url from signing service provider - sending user ' + \
                  'to return_url instead')
    raise falcon.HTTPFound(return_url)


def get_oauth_client():
    return OAuth2Session(os.environ['GH_OAUTH_CLIENT_ID'])
