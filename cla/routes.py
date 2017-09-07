"""
The entry point for the CLA service. Lays out all routes and controller functions.
"""

import hug
import cla.hug_types
import cla.controllers.user
import cla.controllers.project
import cla.controllers.signing
import cla.controllers.signature
import cla.controllers.repository
import cla.controllers.company
import cla.controllers.repository_service
import cla.controllers.github
from cla.utils import get_supported_repository_providers, \
                      get_supported_document_content_types, \
                      get_session_middleware
from falcon import HTTP_403

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
                   'user_company_id': '<org-id>', 'user_github_id': 12345)")
def post_user(user_email: cla.hug_types.email, user_name=None,
              user_company_id=None, user_github_id=None):
    """
    POST: /user

    DATA: {'user_email': 'user@email.com', 'user_name': 'User Name',
           'user_company_id': '<org-id>', 'user_github_id': 12345}

    Returns the data of the newly created user.
    """
    return cla.controllers.user.create_user(user_email=user_email,
                                            user_name=user_name,
                                            user_company_id=user_company_id,
                                            user_github_id=user_github_id)


@hug.put('/user', versions=1,
         examples=" - {'user_id': '<user-id>', 'user_github_id': 23456)")
def put_user(user_id: hug.types.uuid, user_email=None, user_name=None,
             user_company_id=None, user_github_id=None):
    """
    PUT: /user

    DATA: {'user_id': '<user-id>', 'user_github_id': 23456}

    Supports all the same fields as the POST equivalent.

    Returns the data of the updated user.
    """
    return cla.controllers.user.update_user(user_id,
                                            user_email=user_email,
                                            user_name=user_name,
                                            user_company_id=user_company_id,
                                            user_github_id=user_github_id)


@hug.delete('/user/{user_id}', versions=1)
def delete_user(user_id: hug.types.uuid):
    """
    DELETE: /user/{user_id}

    Deletes the specified user.
    """
    return cla.controllers.user.delete_user(user_id)


@hug.get('/user/{user_id}/signatures', versions=1)
def get_user_signatures(user_id: hug.types.uuid):
    """
    GET: /user/{user_id}/signatures

    Returns a list of signatures associated with a user.
    """
    return cla.controllers.user.get_user_signatures(user_id)


@hug.get('/users/company/{user_company_id}', versions=1)
def get_users_company(user_company_id: hug.types.uuid):
    """
    GET: /users/company/{user_company_id}

    Returns a list of users associated with an company.
    """
    return cla.controllers.user.get_users_company(user_company_id)


#
# Signature Routes.
#
@hug.get('/signature', versions=1)
def get_signatures():
    """
    GET: /signature

    Returns all CLA signatures.
    """
    return cla.controllers.signature.get_signatures()


@hug.get('/signature/{signature_id}', versions=1)
def get_signature(signature_id: hug.types.uuid):
    """
    GET: /signature/{signature_id}

    Returns the CLA signature requested by UUID.
    """
    return cla.controllers.signature.get_signature(signature_id)


@hug.post('/signature', versions=1,
          examples=" - {'signature_type': 'cla', 'signature_signed': true, \
                        'signature_approved': true, 'signature_sign_url': 'http://sign.com/here', \
                        'signature_return_url': 'http://cla-system.com/signed', \
                        'signature_project_id': '<project-id>', \
                        'signature_reference_id': '<ref-id>', \
                        'signature_reference_type': 'individual'}")
def post_signature(signature_project_id: hug.types.text, # pylint: disable=too-many-arguments
                   signature_reference_id: hug.types.text,
                   signature_reference_type: hug.types.one_of(['company', 'user']),
                   signature_type: hug.types.one_of(['cla', 'dco']),
                   signature_signed: hug.types.smart_boolean,
                   signature_approved: hug.types.smart_boolean,
                   signature_return_url: cla.hug_types.url,
                   signature_sign_url: cla.hug_types.url):
    """
    POST: /signature

    DATA: {'signature_type': 'cla',
           'signature_signed': true,
           'signature_approved': true,
           'signature_sign_url': 'http://sign.com/here',
           'signature_return_url': 'http://cla-system.com/signed',
           'signature_project_id': '<project-id>',
           'signature_reference_id': '<ref-id>',
           'signature_reference_type': 'individual'}

    signature_reference_type is either 'individual' or 'corporate', depending on the CLA type.
    signature_reference_id needs to reflect the user or company tied to this signature.

    Returns a CLA signatures that was created.
    """
    return cla.controllers.signature.create_signature(signature_project_id,
                                                      signature_reference_id,
                                                      signature_reference_type,
                                                      signature_type=signature_type,
                                                      signature_signed=signature_signed,
                                                      signature_approved=signature_approved,
                                                      signature_return_url=signature_return_url,
                                                      signature_sign_url=signature_sign_url)


@hug.put('/signature', versions=1,
         examples=" - {'signature_id': '01620259-d202-4350-8264-ef42a861922d', \
                       'signature_type': 'cla', 'signature_signed': true}")
def put_signature(signature_id: hug.types.uuid, # pylint: disable=too-many-arguments
                  signature_project_id=None,
                  signature_reference_id=None,
                  signature_reference_type=None,
                  signature_type=None,
                  signature_signed=None,
                  signature_approved=None,
                  signature_return_url=None,
                  signature_sign_url=None):
    """
    PUT: /signature

    DATA: {'signature_id': '<signature-id>',
           'signature_type': 'cla', 'signature_signed': true}

    Supports all the fields as the POST equivalent.

    Returns the CLA signature that was just updated.
    """
    return cla.controllers.signature.update_signature(
        signature_id,
        signature_project_id=signature_project_id,
        signature_reference_id=signature_reference_id,
        signature_reference_type=signature_reference_type,
        signature_type=signature_type,
        signature_signed=signature_signed,
        signature_approved=signature_approved,
        signature_return_url=signature_return_url,
        signature_sign_url=signature_sign_url)


@hug.delete('/signature/{signature_id}', versions=1)
def delete_signature(signature_id: hug.types.uuid):
    """
    DELETE: /signature/{signature_id}

    Deletes the specified signature.
    """
    return cla.controllers.signature.delete_signature(signature_id)


@hug.get('/signatures/user/{user_id}', versions=1)
def get_signatures_user(user_id: hug.types.uuid):
    """
    GET: /signatures/user/{user_id}

    Get all signatures for user specified.
    """
    return cla.controllers.signature.get_user_signatures(user_id)


@hug.get('/signatures/company/{company_id}', versions=1)
def get_signatures_company(company_id: hug.types.uuid):
    """
    GET: /signatures/company/{company_id}

    Get all signatures for company specified.
    """
    return cla.controllers.signature.get_company_signatures(company_id)


@hug.get('/signatures/project/{project_id}', versions=1)
def get_signatures_project(project_id: hug.types.text):
    """
    GET: /signatures/project/{project_id}

    Get all signatures for project specified.
    """
    return cla.controllers.signature.get_project_signatures(project_id)


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
# Company Routes.
#
@hug.get('/company', versions=1)
def get_companies():
    """
    GET: /company

    Returns all CLA companies.
    """
    return cla.controllers.company.get_companies()


@hug.get('/company/{company_id}', versions=1)
def get_company(company_id: hug.types.text):
    """
    GET: /company/{company_id}

    Returns the CLA company requested by UUID.
    """
    return cla.controllers.company.get_company(company_id)


@hug.post('/company', versions=1,
          examples=" - {'company_name': 'Company Name', \
                        'company_whitelist': ['user@safe.org'], \
                        'company_whitelist_patterns': ['*@safe.org']}")
def post_company(company_name: hug.types.text,
                 company_whitelist: hug.types.multiple,
                 company_whitelist_patterns: hug.types.multiple):
    """
    POST: /company

    DATA: {'company_name': 'Org Name',
           'company_whitelist': ['safe@email.org'],
           'company_whitelist': ['*@email.org']}

    Returns the CLA company that was just created.
    """
    return cla.controllers.company.create_company(
        company_name=company_name,
        company_whitelist=company_whitelist,
        company_whitelist_patterns=company_whitelist_patterns)


@hug.put('/company', versions=1,
         examples=" - {'company_id': '<company-id>', \
                       'company_name': 'New Company Name'}")
def put_company(company_id: hug.types.uuid, # pylint: disable=too-many-arguments
                company_name=None,
                company_whitelist=None,
                company_whitelist_patterns=None):
    """
    PUT: /company

    DATA: {'company_id': '<company-id>',
           'company_name': 'New Company Name'}

    Returns the CLA company that was just updated.
    """
    return cla.controllers.company.update_company(
        company_id,
        company_name=company_name,
        company_whitelist=company_whitelist,
        company_whitelist_patterns=company_whitelist_patterns)


@hug.delete('/company/{company_id}', versions=1)
def delete_company(company_id: hug.types.text):
    """
    DELETE: /company/{company_id}

    Deletes the specified company.
    """
    return cla.controllers.company.delete_company(company_id)


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
def post_project(project_external_id: hug.types.text, project_name: hug.types.text, project_ccla_requires_icla_signature: hug.types.boolean):
    """
    POST: /project

    DATA: {'project_external_id': '<proj-external-id>', 'project_name': 'Project Name', 'project_ccla_requires_icla_signature': True}

    Returns the CLA project that was just created.
    """
    return cla.controllers.project.create_project(project_external_id, project_name, project_ccla_requires_icla_signature)


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

    Fetch a project's signature document.
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


@hug.delete('/project/{project_id}/document/{document_type}/{major_version}/{minor_version}', versions=1)
def delete_project_document(project_id: hug.types.text,
                            document_type: hug.types.one_of(['individual', 'corporate']),
                            major_version: hug.types.number,
                            minor_version: hug.types.number):
    """
    DELETE: /project/{project_id}/document/{document_type}/{revision}

    Delete a project's signature document by revision.
    """
    return cla.controllers.project.delete_project_document(project_id,
                                                           document_type,
                                                           major_version,
                                                           minor_version)


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

    Creates a new signature given project and user IDs. The user will be redirected to the
    return_url once signature is complete. If the optional callback_url is provided, the signing
    service provider will hit that URL once user signature is confirmed (typically used to update
    the pull request/merge request/etc.

    Returns a dict of the format:

        {'user_id': <user_id>,
         'signature_id': <signature_id>,
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


@hug.get('/return-url/{signature_id}', versions=1)
def get_return_url(signature_id: hug.types.uuid, event=None):
    """
    GET: /return-url/{signature_id}

    The endpoint the user will be redirected to upon completing signature. Will utilize the
    signature's "signature_return_url" field to redirect the user to the appropriate location.

    Will also capture the signing service provider's return GET parameters, such as DocuSign's
    'event' flag that describes the redirect reason.
    """
    return cla.controllers.signing.return_url(signature_id, event)


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


#
# GitHub Routes.
#
@hug.get('/github/organizations', versions=1)
def get_github_organizations():
    """
    GET: /github/organizations

    Returns all CLA Github Organizations.
    """
    return cla.controllers.github.get_organizations()


@hug.get('/github/organizations/{organization_name}', versions=1)
def get_github_organization(organization_name: hug.types.text):
    """
    GET: /github/organizations/{organization_name}

    Returns the CLA Github Organization requested by Name.
    """
    return cla.controllers.github.get_organization(organization_name)


@hug.get('/github/organizations/{organization_name}/repositories', versions=1)
def get_github_organization_repos(organization_name: hug.types.text):
    """
    GET: /github/organizations/{organization_name}/repositories

    Returns a list of Repositories selected under this organization.
    """
    return cla.controllers.github.get_organization_repositories(organization_name)


@hug.post('/github/organizations', versions=1,
          examples=" - {'organization_project_id': '<project-id>', \
                        'organization_name': 'org-name'}")
def post_github_organization(organization_project_id: hug.types.text, # pylint: disable=too-many-arguments
                             organization_name: hug.types.text):
    """
    POST: /github/organizations

    DATA: {'organization_project_id': '<project-id>',
           'organization_name': 'org-name'}

    Returns the CLA GitHub Organization that was just created.
    """
    return cla.controllers.github.create_organization(organization_project_id,
                                                      organization_name)


@hug.delete('/github/organizations/{organization_name}', versions=1)
def delete_repository(organization_name: hug.types.text):
    """
    DELETE: /github/organizations/{organization_name}

    Deletes the specified Github Organization.
    """
    return cla.controllers.github.delete_organization(organization_name)


@hug.post('/github/installation', versions=1)
def github_app_installation(body):
    """
    POST: /github/installation

    GitHub will fire off this webhook when new installation of our CLA app occurs.
    """
    return cla.controllers.github.user_authorization_callback(body)


@hug.post('/github/activity', versions=1)
def github_app_activity(body, request, response):
    """
    POST: /github/activity

    Acts upon any events triggered by our app installed in someone's organisation.
    """
    # Verify that Webhook Signature is valid
    # if cla.controllers.github.webhook_secret_validation(request.headers.get('X-HUB-SIGNATURE'), request._wrap_stream()):
    return cla.controllers.github.activity(body)
    # else:
    #     response.status = HTTP_403
    #     return {'status': 'Not Authorized'}


@hug.post('/github/validate', versions=1)
def github_organization_validation(body):
    """
    POST: /github/activity

    Acts upon any events triggered by our app installed in someone's organisation.
    """
    return cla.controllers.github.validate_organization(body)
