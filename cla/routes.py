"""
The entry point for the CLA service. Lays out all routes and controller functions.
"""

import hug
import cla.hug_types
import cla.controllers.user
import cla.controllers.project
import cla.controllers.signing
import cla.controllers.agreement
import cla.controllers.repository
import cla.controllers.organization
import cla.controllers.repository_service
from cla.utils import get_supported_repository_providers, \
                      get_supported_document_content_types, \
                      get_session_middleware

#
# Middleware
#
hug.API('cla/routes').http.add_middleware(get_session_middleware())

#
# User routes.
#
@hug.get('/user', versions=1)
def get_users():
    """
    GET: /user

    Returns all CLA users.
    """
    return cla.controllers.user.get_users()

@hug.get('/user/{user_id}', versions=1)
def get_user(user_id: hug.types.uuid):
    """
    GET: /user/{user_id}

    Returns the requested user data based on ID.
    """
    return cla.controllers.user.get_user(user_id=user_id)

@hug.get('/user/email/{user_email}', versions=1)
def get_user_email(user_email: cla.hug_types.email):
    """
    GET: /user/email/{user_email}

    Returns the requested user data based on user email.
    """
    return cla.controllers.user.get_user(user_email=user_email)

@hug.get('/user/github/{user_github_id}', versions=1)
def get_user_github(user_github_id: hug.types.number):
    """
    GET: /user/github/{user_github_id}

    Returns the requested user data based on user GitHub ID.
    """
    return cla.controllers.user.get_user(user_github_id=user_github_id)

@hug.post('/user', versions=1,
          examples=" - {'user_email': 'user@email.com', 'user_name': 'User Name', \
                   'user_organization_id': '<org-id>', 'user_github_id': 12345)")
def post_user(user_email: cla.hug_types.email, user_name=None,
              user_organization_id=None, user_github_id=None):
    """
    POST: /user

    DATA: {'user_email': 'user@email.com', 'user_name': 'User Name',
           'user_organization_id': '<org-id>', 'user_github_id': 12345}

    Returns the data of the newly created user.
    """
    return cla.controllers.user.create_user(user_email=user_email,
                                            user_name=user_name,
                                            user_organization_id=user_organization_id,
                                            user_github_id=user_github_id)

@hug.put('/user', versions=1,
         examples=" - {'user_id': '<user-id>', 'user_github_id': 23456)")
def put_user(user_id: hug.types.uuid, user_email=None, user_name=None,
             user_organization_id=None, user_github_id=None):
    """
    PUT: /user

    DATA: {'user_id': '<user-id>', 'user_github_id': 23456}

    Supports all the same fields as the POST equivalent.

    Returns the data of the updated user.
    """
    return cla.controllers.user.update_user(user_id,
                                            user_email=user_email,
                                            user_name=user_name,
                                            user_organization_id=user_organization_id,
                                            user_github_id=user_github_id)

@hug.delete('/user/{user_id}', versions=1)
def delete_user(user_id: hug.types.uuid):
    """
    DELETE: /user/{user_id}

    Deletes the specified user.
    """
    return cla.controllers.user.delete_user(user_id)

@hug.get('/user/{user_id}/agreements', versions=1)
def get_user_agreements(user_id: hug.types.uuid):
    """
    GET: /user/{user_id}/agreements

    Returns a list of agreements associated with a user.
    """
    return cla.controllers.user.get_user_agreements(user_id)

@hug.get('/users/organization/{user_organization_id}', versions=1)
def get_users_organization(user_organization_id: hug.types.uuid):
    """
    GET: /users/organization/{user_organization_id}

    Returns a list of users associated with an organization.
    """
    return cla.controllers.user.get_users_organization(user_organization_id)

#
# Agreement Routes.
#
@hug.get('/agreement', versions=1)
def get_agreements():
    """
    GET: /agreement

    Returns all CLA agreements.
    """
    return cla.controllers.agreement.get_agreements()

@hug.get('/agreement/{agreement_id}', versions=1)
def get_agreement(agreement_id: hug.types.uuid):
    """
    GET: /agreement/{agreement_id}

    Returns the CLA agreement requested by UUID.
    """
    return cla.controllers.agreement.get_agreement(agreement_id)

@hug.post('/agreement', versions=1,
          examples=" - {'agreement_type': 'cla', 'agreement_signed': true, \
                        'agreement_approved': true, 'agreement_sign_url': 'http://sign.com/here', \
                        'agreement_return_url': 'http://cla-system.com/signed', \
                        'agreement_project_id': '<project-id>', \
                        'agreement_reference_id': '<ref-id>', \
                        'agreement_reference_type': 'individual'}")
def post_agreement(agreement_project_id: hug.types.text, # pylint: disable=too-many-arguments
                   agreement_reference_id: hug.types.text,
                   agreement_reference_type: hug.types.one_of(['organization', 'user']),
                   agreement_type: hug.types.one_of(['cla', 'dco']),
                   agreement_signed: hug.types.smart_boolean,
                   agreement_approved: hug.types.smart_boolean,
                   agreement_return_url: cla.hug_types.url,
                   agreement_sign_url: cla.hug_types.url):
    """
    POST: /agreement

    DATA: {'agreement_type': 'cla',
           'agreement_signed': true,
           'agreement_approved': true,
           'agreement_sign_url': 'http://sign.com/here',
           'agreement_return_url': 'http://cla-system.com/signed',
           'agreement_project_id': '<project-id>',
           'agreement_reference_id': '<ref-id>',
           'agreement_reference_type': 'individual'}

    agreement_reference_type is either 'individual' or 'corporate', depending on the CLA type.
    agreement_reference_id needs to reflect the user or organization tied to this agreement.

    Returns a CLA agreements that was created.
    """
    return cla.controllers.agreement.create_agreement(agreement_project_id,
                                                      agreement_reference_id,
                                                      agreement_reference_type,
                                                      agreement_type=agreement_type,
                                                      agreement_signed=agreement_signed,
                                                      agreement_approved=agreement_approved,
                                                      agreement_return_url=agreement_return_url,
                                                      agreement_sign_url=agreement_sign_url)

@hug.put('/agreement', versions=1,
         examples=" - {'agreement_id': '01620259-d202-4350-8264-ef42a861922d', \
                       'agreement_type': 'cla', 'agreement_signed': true}")
def put_agreement(agreement_id: hug.types.uuid, # pylint: disable=too-many-arguments
                  agreement_project_id=None,
                  agreement_reference_id=None,
                  agreement_reference_type=None,
                  agreement_type=None,
                  agreement_signed=None,
                  agreement_approved=None,
                  agreement_return_url=None,
                  agreement_sign_url=None):
    """
    PUT: /agreement

    DATA: {'agreement_id': '<agreement-id>',
           'agreement_type': 'cla', 'agreement_signed': true}

    Supports all the fields as the POST equivalent.

    Returns the CLA agreement that was just updated.
    """
    return cla.controllers.agreement.update_agreement(
        agreement_id,
        agreement_project_id=agreement_project_id,
        agreement_reference_id=agreement_reference_id,
        agreement_reference_type=agreement_reference_type,
        agreement_type=agreement_type,
        agreement_signed=agreement_signed,
        agreement_approved=agreement_approved,
        agreement_return_url=agreement_return_url,
        agreement_sign_url=agreement_sign_url)

@hug.delete('/agreement/{agreement_id}', versions=1)
def delete_agreement(agreement_id: hug.types.uuid):
    """
    DELETE: /agreement/{agreement_id}

    Deletes the specified agreement.
    """
    return cla.controllers.agreement.delete_agreement(agreement_id)

@hug.get('/agreements/user/{user_id}', versions=1)
def get_agreements_user(user_id: hug.types.uuid):
    """
    GET: /agreements/user/{user_id}

    Get all agreements for user specified.
    """
    return cla.controllers.agreement.get_user_agreements(user_id)

@hug.get('/agreements/organization/{organization_id}', versions=1)
def get_agreements_organization(organization_id: hug.types.uuid):
    """
    GET: /agreements/organization/{organization_id}

    Get all agreements for organization specified.
    """
    return cla.controllers.agreement.get_organization_agreements(organization_id)

@hug.get('/agreements/project/{project_id}', versions=1)
def get_agreements_project(project_id: hug.types.text):
    """
    GET: /agreements/project/{project_id}

    Get all agreements for project specified.
    """
    return cla.controllers.agreement.get_project_agreements(project_id)

#
# Repository Routes.
#
@hug.get('/repository', versions=1)
def get_repositories():
    """
    GET: /repository

    Returns all CLA repositories.
    """
    return cla.controllers.repository.get_repositories()

@hug.get('/repository/{repository_id}', versions=1)
def get_repository(repository_id: hug.types.text):
    """
    GET: /repository/{repository_id}

    Returns the CLA repository requested by UUID.
    """
    return cla.controllers.repository.get_repository(repository_id)

@hug.post('/repository', versions=1,
          examples=" - {'repository_project_id': '<project-id>', \
                        'repository_external_id': 'repo1', \
                        'repository_name': 'Repo Name', \
                        'repository_type': 'github', \
                        'repository_url': 'http://url-to-repo.com'}")
def post_repository(repository_project_id: hug.types.text, # pylint: disable=too-many-arguments
                    repository_name: hug.types.text,
                    repository_type: hug.types.one_of(get_supported_repository_providers().keys()),
                    repository_url: cla.hug_types.url,
                    repository_external_id=None):
    """
    POST: /repository

    DATA: {'repository_project_id': '<project-id>',
           'repository_external_id': 'repo1',
           'repository_name': 'Repo Name',
           'repository_type': 'github',
           'repository_url': 'http://url-to-repo.com'}

    repository_external_id is the ID of the repository given by the repository service provider.
    It is used to redirect the user back to the appropriate location once signing is complete.

    Returns the CLA repository that was just created.
    """
    return cla.controllers.repository.create_repository(repository_project_id,
                                                        repository_name,
                                                        repository_type,
                                                        repository_url,
                                                        repository_external_id)

@hug.put('/repository', versions=1,
         examples=" - {'repository_id': '<repo-id>', \
                       'repository_url': 'http://new-url-to-repository.com'}")
def put_repository(repository_id: hug.types.text, # pylint: disable=too-many-arguments
                   repository_project_id=None,
                   repository_name=None,
                   repository_type=None,
                   repository_url=None,
                   repository_external_id=None):
    """
    PUT: /repository

    DATA: {'repository_id': '<repo-id>',
           'repository_url': 'http://new-url-to-repository.com'}

    Returns the CLA repository that was just updated.
    """
    return cla.controllers.repository.update_repository(
        repository_id,
        repository_project_id=repository_project_id,
        repository_name=repository_name,
        repository_type=repository_type,
        repository_url=repository_url,
        repository_external_id=repository_external_id)

@hug.delete('/repository/{repository_id}', versions=1)
def delete_repository(repository_id: hug.types.text):
    """
    DELETE: /repository/{repository_id}

    Deletes the specified repository.
    """
    return cla.controllers.repository.delete_repository(repository_id)

#
# Organization Routes.
#
@hug.get('/organization', versions=1)
def get_organizations():
    """
    GET: /organization

    Returns all CLA organizations.
    """
    return cla.controllers.organization.get_organizations()

@hug.get('/organization/{organization_id}', versions=1)
def get_organization(organization_id: hug.types.text):
    """
    GET: /organization/{organization_id}

    Returns the CLA organization requested by UUID.
    """
    return cla.controllers.organization.get_organization(organization_id)

@hug.post('/organization', versions=1,
          examples=" - {'organization_name': 'Org Name', \
                        'organization_whitelist': ['safe.org'], \
                        'organization_exclude_patterns': ['^info@*']}")
def post_organization(organization_name: hug.types.text,
                      organization_whitelist: hug.types.multiple,
                      organization_exclude_patterns: hug.types.multiple):
    """
    POST: /organization

    DATA: {'organization_name': 'Org Name',
           'organization_whitelist': ['safe.org'],
           'organization_exclude_patterns': ['^info@*']}

    Returns the CLA organization that was just created.
    """
    return cla.controllers.organization.create_organization(
        organization_name=organization_name,
        organization_whitelist=organization_whitelist,
        organization_exclude_patterns=organization_exclude_patterns)

@hug.put('/organization', versions=1,
         examples=" - {'organization_id': '<org-id>', \
                       'organization_name': 'New Org Name'}")
def put_organization(organization_id: hug.types.uuid, # pylint: disable=too-many-arguments
                     organization_name=None,
                     organization_exclude_patterns=None,
                     organization_whitelist=None):
    """
    PUT: /organization

    DATA: {'organization_id': '<org-id>',
           'organization_name': 'New Org Name'}

    Returns the CLA organization that was just updated.
    """
    return cla.controllers.organization.update_organization(
        organization_id,
        organization_name=organization_name,
        organization_whitelist=organization_whitelist,
        organization_exclude_patterns=organization_exclude_patterns)

@hug.delete('/organization/{organization_id}', versions=1)
def delete_organization(organization_id: hug.types.text):
    """
    DELETE: /organization/{organization_id}

    Deletes the specified organization.
    """
    return cla.controllers.organization.delete_organization(organization_id)

#
# Project Routes.
#
@hug.get('/project', versions=1)
def get_projects():
    """
    GET: /project

    Returns all CLA projects.
    """
    return cla.controllers.project.get_projects()

@hug.get('/project/{project_id}', versions=1)
def get_project(project_id: hug.types.text):
    """
    GET: /project/{project_id}

    Returns the CLA project requested by ID.
    """
    return cla.controllers.project.get_project(project_id)

@hug.post('/project', versions=1,
          examples=" - {'project_name': 'Project Name'}")
def post_project(project_id: hug.types.text, project_name: hug.types.text):
    """
    POST: /project

    DATA: {'project_id': '<proj-id>', 'project_name': 'Project Name'}

    Returns the CLA project that was just created.
    """
    return cla.controllers.project.create_project(project_id, project_name=project_name)

@hug.put('/project', versions=1,
         examples=" - {'project_id': '<proj-id>', \
                       'project_name': 'New Project Name'}")
def put_project(project_id: hug.types.text, project_name=None):
    """
    PUT: /project

    DATA: {'project_id': '<project-id>',
           'project_name': 'New Project Name'}

    Returns the CLA project that was just updated.
    """
    return cla.controllers.project.update_project(project_id, project_name=project_name)

@hug.delete('/project/{project_id}', versions=1)
def delete_project(project_id: hug.types.text):
    """
    DELETE: /project/{project_id}

    Deletes the specified project.
    """
    return cla.controllers.project.delete_project(project_id)

@hug.get('/project/{project_id}/repositories', versions=1)
def get_project_repositories(project_id: hug.types.text):
    """
    GET: /project/{project_id}/repositories

    Gets the specified project's repositories.
    """
    return cla.controllers.project.get_project_repositories(project_id)

@hug.get('/project/{project_id}/document/{document_type}', versions=1)
def get_project_document(project_id: hug.types.text,
                         document_type: hug.types.one_of(['individual', 'corporate'])):
    """
    GET: /project/{project_id}/document/{document_type}

    Fetch a project's agreement document.
    """
    return cla.controllers.project.get_project_document(project_id, document_type)

@hug.post('/project/{project_id}/document/{document_type}', versions=1,
          examples=" - {'document_name': 'doc_name.pdf', \
                        'document_content_type': 'url+pdf', \
                        'document_content': 'http://url.com/doc.pdf'}")
def post_project_document(
        project_id: hug.types.text,
        document_type: hug.types.one_of(['individual', 'corporate']),
        document_name: hug.types.text,
        document_content_type: hug.types.one_of(get_supported_document_content_types()),
        document_content: hug.types.text):
    """
    POST: /project/{project_id}/document/{document_type}

    DATA: {'document_name': 'doc_name.pdf',
           'document_content_type': 'url+pdf',
           'document_content': 'http://url.com/doc.pdf'}

    Creates a new CLA document for a specified project.

    Will create a new revision of the individual or corporate document.

    If document_content_type starts with 'storage+', the document_content is assumed to be base64
    encoded binary data that will be saved in the CLA system's configured storage service.
    """
    return cla.controllers.project.post_project_document(
        project_id=project_id,
        document_type=document_type,
        document_name=document_name,
        document_content_type=document_content_type,
        document_content=document_content)

@hug.delete('/project/{project_id}/document/{document_type}/{revision}', versions=1)
def delete_project_document(project_id: hug.types.text,
                            document_type: hug.types.one_of(['individual', 'corporate']),
                            revision: hug.types.number):
    """
    DELETE: /project/{project_id}/document/{document_type}/{revision}

    Delete a project's agreement document by revision.
    """
    return cla.controllers.project.delete_project_document(project_id, document_type, revision)

#
# Document Signing Routes.
#
@hug.post('/request-signature', versions=1,
          examples=" - {'project_id': 'some-proj-id', \
                        'user_id': 'some-user-uuid', \
                        'return_url': 'https://github.com/linuxfoundation/cla, \
                        'callback_url': 'http://cla.system/signed-callback'}")
def request_signature(project_id: hug.types.text,
                      user_id: hug.types.uuid,
                      return_url: hug.types.text,
                      callback_url=None):
    """
    POST: /request-signature

    DATA: {'project_id': 'some-project-id',
           'user_id': 'some-user-uuid',
           'return_url': 'https://github.com/linuxfoundation/cla,
           'callback_url': 'http://cla.system/signed-callback'}

    Creates a new agreement given project and user IDs. The user will be redirected to the
    return_url once signature is complete. If the optional callback_url is provided, the signing
    service provider will hit that URL once user signature is confirmed (typically used to update
    the pull request/merge request/etc.

    Returns a dict of the format:

        {'user_id': <user_id>,
         'agreement_id': <agreement_id>,
         'project_id': <project_id>,
         'sign_url': <sign_url>}

    User should hit the provided URL to initiate the signing process through the
    signing service provider.
    """
    return cla.controllers.signing.request_signature(project_id,
                                                     user_id,
                                                     return_url,
                                                     callback_url)

@hug.post('/signed/{repository_id}/{change_request_id}', versions=1)
def post_signed(body, repository_id: hug.types.text, change_request_id: hug.types.text):
    """
    POST: /signed/{repository_id}/{change_request_id}

    Callback URL from signing service upon signature.

    If you want the repository service provider to be notified of a successful signature, this
    callback URL should be specified when creating a new signature request with the
    /request-signature endpoint.
    """
    content = body.read()
    return cla.controllers.signing.post_signed(content, repository_id, change_request_id)

@hug.get('/return-url/{agreement_id}', versions=1)
def get_return_url(agreement_id: hug.types.uuid, event=None):
    """
    GET: /return-url/{agreement_id}

    The endpoint the user will be redirected to upon completing signature. Will utilize the
    agreement's "agreement_return_url" field to redirect the user to the appropriate location.

    Will also capture the signing service provider's return GET parameters, such as DocuSign's
    'event' flag that describes the redirect reason.
    """
    return cla.controllers.signing.return_url(agreement_id, event)

#
# Repository Provider Routes.
#
@hug.get('/repository-provider/{provider}/sign/{repository_id}/{change_request_id}', versions=1)
def sign_request(provider: hug.types.one_of(get_supported_repository_providers().keys()),
                 repository_id: hug.types.text,
                 change_request_id: hug.types.text,
                 request):
    """
    GET: /repository-provider/{provider}/sign/{repository_id}/{change_request_id}

    The endpoint that will initiate a CLA signature for the user.
    """
    return cla.controllers.repository_service.sign_request(provider,
                                                           repository_id,
                                                           change_request_id,
                                                           request)

@hug.get('/repository-provider/{provider}/icon.svg', versions=1,
         output=hug.output_format.svg_xml_image) # pylint: disable=no-member
def change_icon(provider: hug.types.one_of(get_supported_repository_providers().keys()),
                signed: hug.types.smart_boolean):
    """
    GET: /repository-provider/{provider}/icon.svg

    Returns the CLA status image for the provider requested.
    """
    return cla.controllers.repository_service.change_icon(provider, signed)

@hug.get('/repository-provider/{provider}/oauth2_redirect', versions=1)
def oauth2_redirect(provider: hug.types.one_of(get_supported_repository_providers().keys()), # pylint: disable=too-many-arguments
                    state: hug.types.text,
                    code: hug.types.text,
                    repository_id: hug.types.text,
                    change_request_id: hug.types.text,
                    request=None):
    """
    GET: /repository-provider/{provider}/oauth2_redirect

    Handles the redirect from an OAuth2 provider when initiating a signature.
    """
    return cla.controllers.repository_service.oauth2_redirect(provider,
                                                              state,
                                                              code,
                                                              repository_id,
                                                              change_request_id,
                                                              request)

@hug.post('/repository-provider/{provider}/activity', versions=1)
def received_activity(body,
                      provider: hug.types.one_of(get_supported_repository_providers().keys())):
    """
    POST: /repository-provider/{provider}/activity

    Acts upon a code repository provider's activity.
    """
    return cla.controllers.repository_service.received_activity(provider,
                                                                body)
