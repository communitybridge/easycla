# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
The entry point for the CLA service. Lays out all routes and controller functions.
"""

import hug
from falcon import HTTP_401, HTTP_400
from hug.middleware import LogMiddleware

import cla
import cla.auth
import cla.controllers.company
import cla.controllers.event
import cla.controllers.gerrit
import cla.controllers.github
import cla.controllers.project
import cla.controllers.project_logo
import cla.controllers.repository
import cla.controllers.repository_service
import cla.controllers.signature
import cla.controllers.signing
import cla.controllers.user
import cla.hug_types
import cla.salesforce
from cla.utils import (
    get_supported_repository_providers,
    get_supported_document_content_types,
    get_session_middleware,
)


#
# Middleware
#

# Session Middleware
# hug.API('cla/routes').http.add_middleware(get_session_middleware())

# CORS Middleware
@hug.response_middleware()
def process_data(request, response, resource):
    # response.set_header('Access-Control-Allow-Origin', cla.conf['ALLOW_ORIGIN'])
    response.set_header("Access-Control-Allow-Origin", "*")
    response.set_header("Access-Control-Allow-Credentials", "true")
    response.set_header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    response.set_header("Access-Control-Allow-Headers", "Content-Type, Authorization")


# Here we comment out the custom 404. Make it back to Hug default 404 behaviour
# Custom 404.
#
# @hug.not_found()
# def not_found():
#     """Custom 404 handler to hide the default hug behaviour of displaying all available routes."""
#     return {'error': {'status': status.HTTP_NOT_FOUND,
#                       'description': 'URL is invalid.'
#                      }
#            }


@hug.directive()
def check_auth(request=None, **kwargs):
    """Returns the authenticated user"""
    return request and cla.auth.authenticate_user(request.headers)


@hug.exception(cla.auth.AuthError)
def handle_auth_error(exception, response=None, **kwargs):
    """Handles authentication errors"""
    response.status = HTTP_401
    return exception.response


#
# Health check route.
#
@hug.get("/health", versions=2)
def get_health(request):
    """
    GET: /health

    Returns a basic health check on the CLA system.
    """
    cla.salesforce.get_projects(request, "")
    request.context["session"]["health"] = "up"
    return request.headers


#
# User routes.
#
# @hug.get('/user', versions=1)
# def get_users():
#     """
#     GET: /user

#     Returns all CLA users.
#     """
#     # staff_verify(user)
#     return cla.controllers.user.get_users()


@hug.get("/user/{user_id}", versions=2)
def get_user(request, user_id: hug.types.uuid):
    """
    GET: /user/{user_id}

    Returns the requested user data based on ID.
    """
    try:
        auth_user = check_auth(request)
    except cla.auth.AuthError as auth_err:
        if auth_err.response == "missing authorization header":
            cla.log.info("getting github user: {}".format(user_id))
        else:
            raise auth_err

    return cla.controllers.user.get_user(user_id=user_id)


# @hug.get('/user/email/{user_email}', versions=1)
# def get_user_email(user_email: cla.hug_types.email, auth_user: check_auth):
#     """
#     GET: /user/email/{user_email}

#     Returns the requested user data based on user email.

#     TODO: Need to look into whether this has to be locked down more (by staff maybe?). Would that
#     break the user flow from GitHub?
#     """
#     return cla.controllers.user.get_user(user_email=user_email)


@hug.post("/user/gerrit", versions=1)
def post_or_get_user_gerrit(auth_user: check_auth):
    """
    GET: /user/gerrit

    For a Gerrit user, there is a case where a user with an lfid may be a user in the db.
    An endpoint to get a userId for gerrit, or create and retrieve the userId if not existent.
    """
    return cla.controllers.user.get_or_create_user(auth_user).to_dict()


# @hug.get('/user/github/{user_github_id}', versions=1)
# def get_user_github(user_github_id: hug.types.number, user: cla_user):
#     """
#     GET: /user/github/{user_github_id}

#     Returns the requested user data based on user GitHub ID.

#     TODO: Should this be locked down more? Staff only?
#     """
#     return cla.controllers.user.get_user(user_github_id=user_github_id)


# @hug.post('/user', versions=1,
#           examples=" - {'user_email': 'user@email.com', 'user_name': 'User Name', \
#                    'user_company_id': '<org-id>', 'user_github_id': 12345)")
# def post_user(user: cla_user, user_email: cla.hug_types.email, user_name=None,
#               user_company_id=None, user_github_id=None):
#     """
#     POST: /user

#     DATA: {'user_email': 'user@email.com', 'user_name': 'User Name',
#            'user_company_id': '<org-id>', 'user_github_id': 12345}

#     Returns the data of the newly created user.
#     """
#     # staff_verify(user) # Only staff can create users.
#     return cla.controllers.user.create_user(user_email=user_email,
#                                             user_name=user_name,
#                                             user_company_id=user_company_id,
#                                             user_github_id=user_github_id)


# @hug.put('/user', versions=1,
#          examples=" - {'user_id': '<user-id>', 'user_github_id': 23456)")
# def put_user(user: cla_user, user_id: hug.types.uuid, user_email=None, user_name=None,
#              user_company_id=None, user_github_id=None):
#     """
#     PUT: /user

#     DATA: {'user_id': '<user-id>', 'user_github_id': 23456}

#     Supports all the same fields as the POST equivalent.

#     Returns the data of the updated user.

#     TODO: Should the user be able to update their own CLA data?
#     """
#     return cla.controllers.user.update_user(user_id,
#                                             user_email=user_email,
#                                             user_name=user_name,
#                                             user_company_id=user_company_id,
#                                             user_github_id=user_github_id)


# @hug.delete('/user/{user_id}', versions=1)
# def delete_user(user: cla_user, user_id: hug.types.uuid):
#     """
#     DELETE: /user/{user_id}

#     Deletes the specified user.
#     """
#     # staff_verify(user)
#     return cla.controllers.user.delete_user(user_id)


@hug.get("/user/{user_id}/signatures", versions=1)
def get_user_signatures(auth_user: check_auth, user_id: hug.types.uuid):
    """
    GET: /user/{user_id}/signatures

    Returns a list of signatures associated with a user.
    """
    return cla.controllers.user.get_user_signatures(user_id)


@hug.get("/users/company/{user_company_id}", versions=1)
def get_users_company(auth_user: check_auth, user_company_id: hug.types.uuid):
    """
    GET: /users/company/{user_company_id}

    Returns a list of users associated with an company.

    TODO: Should probably not simply be auth only - need some role check?
    """
    return cla.controllers.user.get_users_company(user_company_id)


@hug.post("/user/{user_id}/request-company-whitelist/{company_id}", versions=2)
def request_company_whitelist(
    user_id: hug.types.uuid,
    company_id: hug.types.uuid,
    user_email: cla.hug_types.email,
    project_id: hug.types.uuid,
    message=None,
    recipient_name: hug.types.text = None,
    recipient_email: cla.hug_types.email = None,
):
    """
    POST: /user/{user_id}/request-company-whitelist/{company_id}

    DATA: {'user_email': <email-selection>, 'message': 'custom message to manager'}

    Performs the necessary actions (ie: send email to manager) when the specified user requests to
    be added the the specified company's whitelist.
    """
    return cla.controllers.user.request_company_whitelist(
        user_id, str(company_id), str(user_email), str(project_id), message, str(recipient_name), str(recipient_email),
    )


@hug.post("/user/{contributor_id}/invite-company-admin", versions=2)
def invite_company_admin(
    contributor_id: hug.types.uuid,
    contributor_name: cla.hug_types.text,
    contributor_email: cla.hug_types.email,
    cla_manager_name: hug.types.text,
    cla_manager_email: cla.hug_types.email,
    project_name: hug.types.text,
    company_name: hug.types.text,
):
    """
    POST: /user/{user_id}/invite-company-admin

    DATA: {
            'contributor_id': 'uuid-13434-234234-234234',
            'contributor_email': 'Sally Field',
            'contributor_email': 'user@example.com',
            'cla_manager_name': 'John Doe',
            'cla_manager_email': 'admin@example.com',
            'project_name': 'Project Name',
            'company_name': 'Company Name'
        }

    Sends an Email to the prospective CLA Manager to sign up through the ccla console.
    """
    return cla.controllers.user.invite_cla_manager(
        contributor_id, contributor_name, str(contributor_email),
        str(cla_manager_name), str(cla_manager_email),
        project_name, company_name
    )


@hug.post("/user/{user_id}/request-company-ccla", versions=2)
def request_company_ccla(
    user_id: hug.types.uuid, user_email: cla.hug_types.email, company_id: hug.types.uuid, project_id: hug.types.uuid,
):
    """
    POST: /user/{user_id}/request_company_ccla

    Sends an Email to an admin of an existing company to sign a CCLA.
    """
    return cla.controllers.user.request_company_ccla(str(user_id), str(user_email), str(company_id), str(project_id))


@hug.post("/user/{user_id}/company/{company_id}/request-access", versions=2)
def request_company_admin_access(user_id: hug.types.uuid, company_id: hug.types.uuid):
    """
    POST: /user/{user_id}/company/{company_id}/request-access

    Sends an Email for a user requesting access to be on Company ACL.
    """
    return cla.controllers.user.request_company_admin_access(str(user_id), str(company_id))


@hug.get("/user/{user_id}/active-signature", versions=2)
def get_user_active_signature(user_id: hug.types.uuid):
    """
    GET: /user/{user_id}/active-signature

    Returns all metadata associated with a user's active signature.

    {'user_id': <user-id>,
     'project_id': <project-id>,
     'repository_id': <repository-id>,
     'pull_request_id': <PR>,
     'return_url': <url-where-user-initiated-signature-from>'}

    Returns null if the user does not have an active signature.
    """
    return cla.controllers.user.get_active_signature(user_id)


@hug.get("/user/{user_id}/project/{project_id}/last-signature", versions=2)
def get_user_project_last_signature(user_id: hug.types.uuid, project_id: hug.types.uuid):
    """
    GET: /user/{user_id}/project/{project_id}/last-signature

    Returns the user's latest ICLA signature for the project specified.
    """
    return cla.controllers.user.get_user_project_last_signature(user_id, project_id)


@hug.get("/user/{user_id}/project/{project_id}/last-signature/{company_id}", versions=1)
def get_user_project_company_last_signature(
    user_id: hug.types.uuid, project_id: hug.types.uuid, company_id: hug.types.uuid
):
    """
    GET: /user/{user_id}/project/{project_id}/last-signature/{company_id}

    Returns the user's latest employee signature for the project and company specified.
    """
    return cla.controllers.user.get_user_project_company_last_signature(user_id, project_id, company_id)


# #
# # Signature Routes.
# #
# @hug.get('/signature', versions=1)
# def get_signatures(auth_user: check_auth):
#     """
#     GET: /signature

#     Returns all CLA signatures.
#     """
#     # staff_verify(user)
#     return cla.controllers.signature.get_signatures()


@hug.get("/signature/{signature_id}", versions=1)
def get_signature(auth_user: check_auth, signature_id: hug.types.uuid):
    """
    GET: /signature/{signature_id}

    Returns the CLA signature requested by UUID.
    """
    return cla.controllers.signature.get_signature(signature_id)


@hug.post(
    "/signature",
    versions=1,
    examples=" - {'signature_type': 'cla', 'signature_signed': true, \
                        'signature_approved': true, 'signature_sign_url': 'http://sign.com/here', \
                        'signature_return_url': 'http://cla-system.com/signed', \
                        'signature_project_id': '<project-id>', \
                        'signature_reference_id': '<ref-id>', \
                        'signature_reference_type': 'individual'}",
)
def post_signature(
    auth_user: check_auth,  # pylint: disable=too-many-arguments
    signature_project_id: hug.types.uuid,
    signature_reference_id: hug.types.text,
    signature_reference_type: hug.types.one_of(["company", "user"]),
    signature_type: hug.types.one_of(["cla", "dco"]),
    signature_signed: hug.types.smart_boolean,
    signature_approved: hug.types.smart_boolean,
    signature_return_url: cla.hug_types.url,
    signature_sign_url: cla.hug_types.url,
    signature_user_ccla_company_id=None,
):
    """
    POST: /signature

    DATA: {'signature_type': 'cla',
           'signature_signed': true,
           'signature_approved': true,
           'signature_sign_url': 'http://sign.com/here',
           'signature_return_url': 'http://cla-system.com/signed',
           'signature_project_id': '<project-id>',
           'signature_user_ccla_company_id': '<company-id>',
           'signature_reference_id': '<ref-id>',
           'signature_reference_type': 'individual'}

    signature_reference_type is either 'individual' or 'corporate', depending on the CLA type.
    signature_reference_id needs to reflect the user or company tied to this signature.

    Returns a CLA signatures that was created.
    """
    return cla.controllers.signature.create_signature(
        signature_project_id,
        signature_reference_id,
        signature_reference_type,
        signature_type=signature_type,
        signature_user_ccla_company_id=signature_user_ccla_company_id,
        signature_signed=signature_signed,
        signature_approved=signature_approved,
        signature_return_url=signature_return_url,
        signature_sign_url=signature_sign_url,
    )


@hug.put(
    "/signature",
    versions=1,
    examples=" - {'signature_id': '01620259-d202-4350-8264-ef42a861922d', \
                       'signature_type': 'cla', 'signature_signed': true}",
)
def put_signature(
    auth_user: check_auth,  # pylint: disable=too-many-arguments
    signature_id: hug.types.uuid,
    signature_project_id=None,
    signature_reference_id=None,
    signature_reference_type=None,
    signature_type=None,
    signature_signed=None,
    signature_approved=None,
    signature_return_url=None,
    signature_sign_url=None,
    domain_whitelist=None,
    email_whitelist=None,
    github_whitelist=None,
    github_org_whitelist=None,
):
    """
    PUT: /signature

    DATA: {'signature_id': '<signature-id>',
           'signature_type': 'cla', 'signature_signed': true}

    Supports all the fields as the POST equivalent.

    Returns the CLA signature that was just updated.
    """
    return cla.controllers.signature.update_signature(
        signature_id,
        auth_user=auth_user,
        signature_project_id=signature_project_id,
        signature_reference_id=signature_reference_id,
        signature_reference_type=signature_reference_type,
        signature_type=signature_type,
        signature_signed=signature_signed,
        signature_approved=signature_approved,
        signature_return_url=signature_return_url,
        signature_sign_url=signature_sign_url,
        domain_whitelist=domain_whitelist,
        email_whitelist=email_whitelist,
        github_whitelist=github_whitelist,
        github_org_whitelist=github_org_whitelist,
    )


@hug.delete("/signature/{signature_id}", versions=1)
def delete_signature(auth_user: check_auth, signature_id: hug.types.uuid):
    """
    DELETE: /signature/{signature_id}

    Deletes the specified signature.
    """
    # staff_verify(user)
    return cla.controllers.signature.delete_signature(signature_id)


@hug.get("/signatures/user/{user_id}", versions=1)
def get_signatures_user(auth_user: check_auth, user_id: hug.types.uuid):
    """
    GET: /signatures/user/{user_id}

    Get all signatures for user specified.
    """
    return cla.controllers.signature.get_user_signatures(user_id)


@hug.get("/signatures/user/{user_id}/project/{project_id}", versions=1)
def get_signatures_user_project(auth_user: check_auth, user_id: hug.types.uuid, project_id: hug.types.uuid):
    """
    GET: /signatures/user/{user_id}/project/{project_id}

    Get all signatures for user, filtered by project_id specified.
    """
    return cla.controllers.signature.get_user_project_signatures(user_id, project_id)


@hug.get("/signatures/user/{user_id}/project/{project_id}/type/{signature_type}", versions=1)
def get_signatures_user_project(
    auth_user: check_auth,
    user_id: hug.types.uuid,
    project_id: hug.types.uuid,
    signature_type: hug.types.one_of(["individual", "employee"]),
):
    """
    GET: /signatures/user/{user_id}/project/{project_id}/type/[individual|corporate|employee]

    Get all signatures for user, filtered by project_id and signature type specified.
    """
    return cla.controllers.signature.get_user_project_signatures(user_id, project_id, signature_type)


@hug.get("/signatures/company/{company_id}", versions=1)
def get_signatures_company(auth_user: check_auth, company_id: hug.types.uuid):
    """
    GET: /signatures/company/{company_id}

    Get all signatures for company specified.
    """
    return cla.controllers.signature.get_company_signatures_by_acl(auth_user.username, company_id)


@hug.get("/signatures/project/{project_id}", versions=1)
def get_signatures_project(auth_user: check_auth, project_id: hug.types.uuid):
    """
    GET: /signatures/project/{project_id}

    Get all signatures for project specified.
    """
    return cla.controllers.signature.get_project_signatures(project_id)


@hug.get("/signatures/company/{company_id}/project/{project_id}", versions=1)
def get_signatures_project_company(company_id: hug.types.uuid, project_id: hug.types.uuid):
    """
     GET: /signatures/company/{company_id}/project/{project_id}

     Get all signatures for project specified and a company specified
     """
    return cla.controllers.signature.get_project_company_signatures(company_id, project_id)


@hug.get("/signatures/company/{company_id}/project/{project_id}/employee", versions=1)
def get_project_employee_signatures(company_id: hug.types.uuid, project_id: hug.types.uuid):
    """
     GET: /signatures/company/{company_id}/project/{project_id}

     Get all employee signatures for project specified and a company specified
     """
    return cla.controllers.signature.get_project_employee_signatures(company_id, project_id)


@hug.get("/signature/{signature_id}/manager", versions=1)
def get_cla_managers(auth_user: check_auth, signature_id: hug.types.uuid):
    """
    GET: /project/{project_id}/managers

    Returns the CLA Managers from a CCLA's signature ACL.
    """
    return cla.controllers.signature.get_cla_managers(auth_user.username, signature_id)


@hug.post("/signature/{signature_id}/manager", versions=1)
def add_cla_manager(auth_user: check_auth, signature_id: hug.types.uuid, lfid: hug.types.text):
    """
    POST: /project/{project_id}/manager

    Adds CLA Manager to a CCLA's signature ACL and returns the new list of CLA managers.
    """
    return cla.controllers.signature.add_cla_manager(auth_user, signature_id, lfid)


@hug.delete("/signature/{signature_id}/manager/{lfid}", versions=1)
def remove_cla_manager(auth_user: check_auth, signature_id: hug.types.uuid, lfid: hug.types.text):
    """
    DELETE: /signature/{signature_id}/manager/{lfid}

    Removes a CLA Manager from a CCLA's signature ACL and returns the modified list of CLA Managers.
    """
    return cla.controllers.signature.remove_cla_manager(auth_user.username, signature_id, lfid)


#
# Repository Routes.
#
# @hug.get('/repository', versions=1)
# def get_repositories(auth_user: check_auth):
#     """
#     GET: /repository

#     Returns all CLA repositories.
#     """
#     # staff_verify(user)
#     return cla.controllers.repository.get_repositories()


@hug.get("/repository/{repository_id}", versions=1)
def get_repository(auth_user: check_auth, repository_id: hug.types.text):
    """
    GET: /repository/{repository_id}

    Returns the CLA repository requested by UUID.
    """
    return cla.controllers.repository.get_repository(repository_id)


@hug.post(
    "/repository",
    versions=1,
    examples=" - {'repository_project_id': '<project-id>', \
                        'repository_external_id': 'repo1', \
                        'repository_name': 'Repo Name', \
                        'repository_organization_name': 'Organization Name', \
                        'repository_type': 'github', \
                        'repository_url': 'http://url-to-repo.com'}",
)
def post_repository(
    auth_user: check_auth,  # pylint: disable=too-many-arguments
    repository_project_id: hug.types.uuid,
    repository_name: hug.types.text,
    repository_organization_name: hug.types.text,
    repository_type: hug.types.one_of(get_supported_repository_providers().keys()),
    repository_url: cla.hug_types.url,
    repository_external_id=None,
):
    """
    POST: /repository

    DATA: {'repository_project_id': '<project-id>',
           'repository_external_id': 'repo1',
           'repository_name': 'Repo Name',
           'repository_organization_name': 'Organization Name',
           'repository_type': 'github',
           'repository_url': 'http://url-to-repo.com'}

    repository_external_id is the ID of the repository given by the repository service provider.
    It is used to redirect the user back to the appropriate location once signing is complete.

    Returns the CLA repository that was just created.
    """
    return cla.controllers.repository.create_repository(
        auth_user,
        repository_project_id,
        repository_name,
        repository_organization_name,
        repository_type,
        repository_url,
        repository_external_id,
    )


@hug.put(
    "/repository",
    versions=1,
    examples=" - {'repository_id': '<repo-id>', \
                       'repository_id': 'http://new-url-to-repository.com'}",
)
def put_repository(
    auth_user: check_auth,  # pylint: disable=too-many-arguments
    repository_id: hug.types.text,
    repository_project_id=None,
    repository_name=None,
    repository_type=None,
    repository_url=None,
    repository_external_id=None,
):
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
        repository_external_id=repository_external_id,
    )


@hug.delete("/repository/{repository_id}", versions=1)
def delete_repository(auth_user: check_auth, repository_id: hug.types.text):
    """
    DELETE: /repository/{repository_id}

    Deletes the specified repository.
    """
    # staff_verify(user)
    return cla.controllers.repository.delete_repository(repository_id)


# #
# # Company Routes.
# #
@hug.get("/company", versions=1)
def get_companies(auth_user: check_auth):
    """
    GET: /company

    Returns all CLA companies associated with user.
    """
    cla.controllers.user.get_or_create_user(auth_user)  # Find or Create user -- For first login
    return cla.controllers.company.get_companies_by_user(auth_user.username)


@hug.get("/company", versions=2)
def get_all_companies():
    """
    GET: /company

    Returns all CLA companies.
    """
    return cla.controllers.company.get_companies()


@hug.get("/company/{company_id}", versions=2)
def get_company(company_id: hug.types.text):
    """
    GET: /company/{company_id}

    Returns the CLA company requested by UUID.
    """
    return cla.controllers.company.get_company(company_id)


@hug.get("/company/{company_id}/project/unsigned", versions=1)
def get_unsigned_projects_for_company(company_id: hug.types.text):
    """
    GET: /company/{company_id}/project/unsigned

    Returns a list of projects that the company has not signed CCLAs for.
    """
    return cla.controllers.project.get_unsigned_projects_for_company(company_id)


@hug.post(
    "/company",
    versions=1,
    examples=" - {'company_name': 'Company Name', \
                        'company_manager_id': 'user-id'}",
)
def post_company(
    auth_user: check_auth,
    company_name: hug.types.text,
    company_manager_user_name=None,
    company_manager_user_email=None,
    company_manager_id=None,
    response=None
):
    """
    POST: /company

    DATA: {'company_name': 'Org Name',
           'company_manager_id': <user-id>}

    Returns the CLA company that was just created.
    """

    create_resp = cla.controllers.company.create_company(
        auth_user,
        company_name=company_name,
        company_manager_id=company_manager_id,
        company_manager_user_name=company_manager_user_name,
        company_manager_user_email=company_manager_user_email,
        response=response,
    )

    response.status = create_resp.get("status_code")

    return create_resp.get("data")


@hug.put(
    "/company",
    versions=1,
    examples=" - {'company_id': '<company-id>', \
                       'company_name': 'New Company Name'}",
)
def put_company(
    auth_user: check_auth,  # pylint: disable=too-many-arguments
    company_id: hug.types.uuid,
    company_name=None,
    company_manager_id=None,
):
    """
    PUT: /company

    DATA: {'company_id': '<company-id>',
           'company_name': 'New Company Name'}

    Returns the CLA company that was just updated.
    """

    return cla.controllers.company.update_company(
        company_id, company_name=company_name, company_manager_id=company_manager_id, username=auth_user.username,
    )


@hug.delete("/company/{company_id}", versions=1)
def delete_company(auth_user: check_auth, company_id: hug.types.text):
    """
    DELETE: /company/{company_id}

    Deletes the specified company.
    """
    # staff_verify(user)
    return cla.controllers.company.delete_company(company_id, username=auth_user.username)


@hug.put("/company/{company_id}/import/whitelist/csv", versions=1)
def put_company_whitelist_csv(body, auth_user: check_auth, company_id: hug.types.uuid):
    """
    PUT: /company/{company_id}/import/whitelist/csv

    Imports a CSV file of whitelisted user emails.
    Expects the first column to have a header in the first row and contain email addresses.
    """
    # staff_verify(user) or company_manager_verify(user, company_id)
    content = body.read().decode()
    return cla.controllers.company.update_company_whitelist_csv(content, company_id, username=auth_user.username)


@hug.get("/companies/{manager_id}", version=1)
def get_manager_companies(manager_id: hug.types.uuid):
    """
    GET: /companies/{manager_id}

    Returns a list of companies a manager is associated with
    """
    return cla.controllers.company.get_manager_companies(manager_id)


# #
# # Project Routes.
# #
@hug.get("/project", versions=1)
def get_projects(auth_user: check_auth):
    """
    GET: /project

    Returns all CLA projects.
    """
    # staff_verify(user)
    projects = cla.controllers.project.get_projects()
    # For public endpoint, don't show the project_external_id.
    for project in projects:
        if "project_external_id" in project:
            del project["project_external_id"]
    return projects


@hug.get("/project/{project_id}", versions=2)
def get_project(project_id: hug.types.uuid):
    """
    GET: /project/{project_id}

    Returns the CLA project requested by ID.
    """
    project = cla.controllers.project.get_project(project_id)
    # For public endpoint, don't show the project_external_id.
    if "project_external_id" in project:
        del project["project_external_id"]
    return project


@hug.get("/project/{project_id}/manager", versions=1)
def get_project_managers(auth_user: check_auth, project_id: hug.types.uuid):
    """
    GET: /project/{project_id}/managers
    Returns the CLA project managers.
    """
    return cla.controllers.project.get_project_managers(auth_user.username, project_id, enable_auth=True)


@hug.post("/project/{project_id}/manager", versions=1)
def add_project_manager(auth_user: check_auth, project_id: hug.types.text, lfid: hug.types.text):
    """
    POST: /project/{project_id}/manager
    Returns the new list of project managers
    """
    return cla.controllers.project.add_project_manager(auth_user.username, project_id, lfid)


@hug.delete("/project/{project_id}/manager/{lfid}", versions=1)
def remove_project_manager(auth_user: check_auth, project_id: hug.types.text, lfid: hug.types.text):
    """
    DELETE: /project/{project_id}/project/{lfid}
    Returns a success message if it was deleted
    """
    return cla.controllers.project.remove_project_manager(auth_user.username, project_id, lfid)


@hug.get("/project/external/{project_external_id}", version=1)
def get_external_project(auth_user: check_auth, project_external_id: hug.types.text):
    """
    GET: /project/external/{project_external_id}

    Returns the list of CLA projects marching the requested external ID.
    """
    return cla.controllers.project.get_projects_by_external_id(project_external_id, auth_user.username)


@hug.post("/project", versions=1, examples=" - {'project_name': 'Project Name'}")
def post_project(
    auth_user: check_auth,
    project_external_id: hug.types.text,
    project_name: hug.types.text,
    project_icla_enabled: hug.types.boolean,
    project_ccla_enabled: hug.types.boolean,
    project_ccla_requires_icla_signature: hug.types.boolean,
):
    """
    POST: /project

    DATA: {'project_external_id': '<proj-external-id>', 'project_name': 'Project Name',
           'project_icla_enabled': True, 'project_ccla_enabled': True,
           'project_ccla_requires_icla_signature': True}

    Returns the CLA project that was just created.
    """
    # staff_verify(user) or pm_verify_external_id(user, project_external_id)

    return cla.controllers.project.create_project(
        project_external_id,
        project_name,
        project_icla_enabled,
        project_ccla_enabled,
        project_ccla_requires_icla_signature,
        auth_user.username,
    )


@hug.put(
    "/project",
    versions=1,
    examples=" - {'project_id': '<proj-id>', \
                       'project_name': 'New Project Name'}",
)
def put_project(
    auth_user: check_auth,
    project_id: hug.types.uuid,
    project_name=None,
    project_icla_enabled=None,
    project_ccla_enabled=None,
    project_ccla_requires_icla_signature=None,
):
    """
    PUT: /project

    DATA: {'project_id': '<project-id>',
           'project_name': 'New Project Name'}

    Returns the CLA project that was just updated.
    """
    # staff_verify(user) or pm_verify(user, project_id)
    return cla.controllers.project.update_project(
        project_id,
        project_name=project_name,
        project_icla_enabled=project_icla_enabled,
        project_ccla_enabled=project_ccla_enabled,
        project_ccla_requires_icla_signature=project_ccla_requires_icla_signature,
        username=auth_user.username,
    )


@hug.delete("/project/{project_id}", versions=1)
def delete_project(auth_user: check_auth, project_id: hug.types.uuid):
    """
    DELETE: /project/{project_id}

    Deletes the specified project.
    """
    # staff_verify(user)
    return cla.controllers.project.delete_project(project_id, username=auth_user.username)


@hug.get("/project/{project_id}/repositories", versions=1)
def get_project_repositories(auth_user: check_auth, project_id: hug.types.uuid):
    """
    GET: /project/{project_id}/repositories

    Gets the specified project's repositories.
    """
    return cla.controllers.project.get_project_repositories(auth_user, project_id)


@hug.get("/project/{project_id}/repositories_group_by_organization", versions=1)
def get_project_repositories_group_by_organization(auth_user: check_auth, project_id: hug.types.uuid):
    """
    GET: /project/{project_id}/repositories_by_org

    Gets the specified project's repositories. grouped by organization name
    """
    return cla.controllers.project.get_project_repositories_group_by_organization(auth_user, project_id)


@hug.get("/project/{project_id}/configuration_orgs_and_repos", versions=1)
def get_project_configuration_orgs_and_repos(auth_user: check_auth, project_id: hug.types.uuid):
    """
    GET: /project/{project_id}/configuration_orgs_and_repos

    Gets the repositories from github api
    Gets all repositories for from an sfdc project ID
    """
    return cla.controllers.project.get_project_configuration_orgs_and_repos(auth_user, project_id)


@hug.get("/project/{project_id}/document/{document_type}", versions=2)
def get_project_document(
    project_id: hug.types.uuid, document_type: hug.types.one_of(["individual", "corporate"]),
):
    """
    GET: /project/{project_id}/document/{document_type}

    Fetch a project's signature document.
    """
    return cla.controllers.project.get_project_document(project_id, document_type)


@hug.get("/project/{project_id}/document/{document_type}/pdf", version=2)
def get_project_document_raw(
    response,
    auth_user: check_auth,
    project_id: hug.types.uuid,
    document_type: hug.types.one_of(["individual", "corporate"]),
):
    """
    GET: /project/{project_id}/document/{document_type}/pdf

    Returns the PDF document matching the latest individual or corporate contract for that project.
    """
    response.set_header("Content-Type", "application/pdf")
    return cla.controllers.project.get_project_document_raw(project_id, document_type)


@hug.get(
    "/project/{project_id}/document/{document_type}/pdf/{document_major_version}/{document_minor_version}", version=1,
)
def get_project_document_matching_version(
    response,
    auth_user: check_auth,
    project_id: hug.types.uuid,
    document_type: hug.types.one_of(["individual", "corporate"]),
    document_major_version: hug.types.number,
    document_minor_version: hug.types.number,
):
    """
    GET: /project/{project_id}/document/{document_type}/pdf/{document_major_version}/{document_minor_version}

    Returns the PDF document version matching the individual or corporate contract for that project.
    """
    response.set_header("Content-Type", "application/pdf")
    return cla.controllers.project.get_project_document_raw(
        project_id,
        document_type,
        document_major_version=document_major_version,
        document_minor_version=document_minor_version,
    )


@hug.get("/project/{project_id}/companies", versions=2)
def get_project_companies(project_id: hug.types.uuid):
    """
    GET: /project/{project_id}/companies
s
    Check if project exists and retrieves all companies
    """
    return cla.controllers.project.get_project_companies(project_id)


@hug.post(
    "/project/{project_id}/document/{document_type}",
    versions=1,
    examples=" - {'document_name': 'doc_name.pdf', \
                        'document_content_type': 'url+pdf', \
                        'document_content': 'http://url.com/doc.pdf', \
                        'new_major_version': true}",
)
def post_project_document(
    auth_user: check_auth,
    project_id: hug.types.uuid,
    document_type: hug.types.one_of(["individual", "corporate"]),
    document_name: hug.types.text,
    document_content_type: hug.types.one_of(get_supported_document_content_types()),
    document_content: hug.types.text,
    document_preamble=None,
    document_legal_entity_name=None,
    new_major_version=None,
):
    """
    POST: /project/{project_id}/document/{document_type}

    DATA: {'document_name': 'doc_name.pdf',
           'document_content_type': 'url+pdf',
           'document_content': 'http://url.com/doc.pdf',
           'document_preamble': 'Preamble here',
           'document_legal_entity_name': 'Legal entity name',
           'new_major_version': false}

    Creates a new CLA document for a specified project.

    Will create a new revision of the individual or corporate document. if new_major_version is set,
    the document will have a new major version and this will force users to re-sign.

    If document_content_type starts with 'storage+', the document_content is assumed to be base64
    encoded binary data that will be saved in the CLA system's configured storage service.
    """
    # staff_verify(user) or pm_verify(user, project_id)
    return cla.controllers.project.post_project_document(
        project_id=project_id,
        document_type=document_type,
        document_name=document_name,
        document_content_type=document_content_type,
        document_content=document_content,
        document_preamble=document_preamble,
        document_legal_entity_name=document_legal_entity_name,
        new_major_version=new_major_version,
        username=auth_user.username,
    )


@hug.post(
    "/project/{project_id}/document/template/{document_type}",
    versions=1,
    examples=" - {'document_name': 'doc_name.pdf', \
                        'document_preamble': 'Preamble here', \
                        'document_legal_entity_name': 'Legal entity name', \
                        'template_name': 'CNCFTemplate', \
                        'new_major_version': true}",
)
def post_project_document_template(
    auth_user: check_auth,
    project_id: hug.types.uuid,
    document_type: hug.types.one_of(["individual", "corporate"]),
    document_name: hug.types.text,
    document_preamble: hug.types.text,
    document_legal_entity_name: hug.types.text,
    template_name: hug.types.one_of(
        [
            "CNCFTemplate",
            "OpenBMCTemplate",
            "TungstenFabricTemplate",
            "OpenColorIOTemplate",
            "OpenVDBTemplate",
            "ONAPTemplate",
            "TektonTemplate",
        ]
    ),
    new_major_version=None,
):
    """
    POST: /project/{project_id}/document/template/{document_type}

#     DATA: {'document_name': 'doc_name.pdf',
#            'document_preamble': 'Preamble here',
#            'document_legal_entity_name': 'Legal entity name',
#            'template_name': 'CNCFTemplate',
#            'new_major_version': false}

#     Creates a new CLA document from a template for a specified project.

#     Will create a new revision of the individual or corporate document. if new_major_version is set,
#     the document will have a new major version and this will force users to re-sign.

#     The document_content_type is assumed to be 'storage+pdf', which means the document content will
#     be saved in the CLA system's configured storage service.
#     """
    # staff_verify(user) or pm_verify(user, project_id)
    return cla.controllers.project.post_project_document_template(
        project_id=project_id,
        document_type=document_type,
        document_name=document_name,
        document_preamble=document_preamble,
        document_legal_entity_name=document_legal_entity_name,
        template_name=template_name,
        new_major_version=new_major_version,
        username=auth_user.username,
    )


@hug.delete(
    "/project/{project_id}/document/{document_type}/{major_version}/{minor_version}", versions=1,
)
def delete_project_document(
    auth_user: check_auth,
    project_id: hug.types.uuid,
    document_type: hug.types.one_of(["individual", "corporate"]),
    major_version: hug.types.number,
    minor_version: hug.types.number,
):
    #     """
    #     DELETE: /project/{project_id}/document/{document_type}/{revision}

    #     Delete a project's signature document by revision.
    #     """
    #     # staff_verify(user)
    return cla.controllers.project.delete_project_document(
        project_id, document_type, major_version, minor_version, username=auth_user.username,
    )


# #
# # Document Signing Routes.
# #
@hug.post(
    "/request-individual-signature",
    versions=2,
    examples=" - {'project_id': 'some-proj-id', \
                        'user_id': 'some-user-uuid'}",
)
def request_individual_signature(
    project_id: hug.types.uuid, user_id: hug.types.uuid, return_url_type=None, return_url=None,
):
    """
    POST: /request-individual-signature

    DATA: {'project_id': 'some-project-id',
           'user_id': 'some-user-id',
           'return_url_type': Gerrit/Github. Optional depending on presence of return_url
           'return_url': <optional>}

    Creates a new signature given project and user IDs. The user will be redirected to the
    return_url once signature is complete.

    Returns a dict of the format:

        {'user_id': <user_id>,
         'signature_id': <signature_id>,
         'project_id': <project_id>,
         'sign_url': <sign_url>}

    User should hit the provided URL to initiate the signing process through the
    signing service provider.
    """
    return cla.controllers.signing.request_individual_signature(project_id, user_id, return_url_type, return_url)


@hug.post(
    "/request-corporate-signature",
    versions=1,
    examples=" - {'project_id': 'some-proj-id', \
                        'company_id': 'some-company-uuid'}",
)
def request_corporate_signature(
    auth_user: check_auth,
    project_id: hug.types.uuid,
    company_id: hug.types.uuid,
    send_as_email=False,
    authority_name=None,
    authority_email=None,
    return_url_type=None,
    return_url=None,
):
    """
    POST: /request-corporate-signature

    DATA: {'project_id': 'some-project-id',
           'company_id': 'some-company-id',
           'send_as_email': 'boolean',
           'authority_name': 'string',
           'authority_email': 'string',
           'return_url': <optional>}

    Creates a new signature given project and company IDs. The manager will be redirected to the
    return_url once signature is complete.

    The send_as_email flag determines whether to send the signing document because the CLA signatory/signer
    may not necessarily be a corporate/company manager/authority with signing privileges (e.g. may be the
    company manager, but not responsible for signing the CLAs).

    Returns a dict of the format:

        {'company_id': <user_id>,
         'signature_id': <signature_id>,
         'project_id': <project_id>,
         'sign_url': <sign_url>}

    Manager should hit the provided URL to initiate the signing process through the
    signing service provider.
    """
    # staff_verify(user) or company_manager_verify(user, company_id)
    return cla.controllers.signing.request_corporate_signature(
        auth_user, project_id, company_id, send_as_email, authority_name, authority_email, return_url_type, return_url,
    )


@hug.post("/request-employee-signature", versions=2)
def request_employee_signature(
    project_id: hug.types.uuid,
    company_id: hug.types.uuid,
    user_id: hug.types.uuid,
    return_url_type: hug.types.text,
    return_url=None,
):
    """
    POST: /request-employee-signature

    DATA: {'project_id': <project-id>,
           'company_id': <company-id>,
           'user_id': <user-id>,
           'return_url': <optional>}

    Creates a placeholder signature object that represents an employee of a company having confirmed
    that they indeed work for company X which already has a CCLA with the project. This does not
    require a full DocuSign signature process, which means the sign/callback URLs and document
    versions may not be populated or reliable.
    """
    return cla.controllers.signing.request_employee_signature(
        project_id, company_id, user_id, return_url_type, return_url
    )


@hug.post("/check-prepare-employee-signature", versions=2)
def check_and_prepare_employee_signature(
    project_id: hug.types.uuid, company_id: hug.types.uuid, user_id: hug.types.uuid
):
    """
    POST: /check-prepare-employee-signature

    DATA: {'project_id': <project-id>,
           'company_id': <company-id>,
           'user_id': <user-id>
           }

    Checks if an employee is ready to sign a CCLA for a company.
    """
    return cla.controllers.signing.check_and_prepare_employee_signature(project_id, company_id, user_id)


@hug.post(
    "/signed/individual/{installation_id}/{github_repository_id}/{change_request_id}", versions=2,
)
def post_individual_signed(
    body,
    installation_id: hug.types.number,
    github_repository_id: hug.types.number,
    change_request_id: hug.types.number,
):
    """
    POST: /signed/individual/{installation_id}/{github_repository_id}/{change_request_id}

    TODO: Need to protect this endpoint somehow - at the very least ensure it's coming from
    DocuSign and the data hasn't been tampered with.

    Callback URL from signing service upon ICLA signature.
    """
    content = body.read()
    return cla.controllers.signing.post_individual_signed(
        content, installation_id, github_repository_id, change_request_id
    )


@hug.post("/signed/gerrit/individual/{user_id}", versions=2)
def post_individual_signed_gerrit(body, user_id: hug.types.uuid):
    """
    POST: /signed/gerritindividual/{user_id}

    Callback URL from signing service upon ICLA signature for a Gerrit user.
    """
    content = body.read()
    return cla.controllers.signing.post_individual_signed_gerrit(content, user_id)


@hug.post("/signed/corporate/{project_id}/{company_id}", versions=2)
def post_corporate_signed(body, project_id: hug.types.uuid, company_id: hug.types.uuid):
    """
    POST: /signed/corporate/{project_id}/{company_id}

    TODO: Need to protect this endpoint somehow - at the very least ensure it's coming from
    DocuSign and the data hasn't been tampered with.

    Callback URL from signing service upon CCLA signature.
    """
    content = body.read()
    return cla.controllers.signing.post_corporate_signed(content, project_id, company_id)


@hug.get("/return-url/{signature_id}", versions=2)
def get_return_url(signature_id: hug.types.uuid, event=None):
    """
    GET: /return-url/{signature_id}

    The endpoint the user will be redirected to upon completing signature. Will utilize the
    signature's "signature_return_url" field to redirect the user to the appropriate location.

    Will also capture the signing service provider's return GET parameters, such as DocuSign's
    'event' flag that describes the redirect reason.
    """
    return cla.controllers.signing.return_url(signature_id, event)


@hug.post("/send-authority-email", versions=2)
def send_authority_email(
    auth_user: check_auth,
    company_name: hug.types.text,
    project_name: hug.types.text,
    authority_name: hug.types.text,
    authority_email: cla.hug_types.email,
):
    """
    POST: /send-authority-email

    DATA: {
            'authority_name': John Doe,
            'authority_email': authority@example.com,
            'company_id': <company_id>
            'project_id': <project_id>
        }
    """
    return cla.controllers.signing.send_authority_email(company_name, project_name, authority_name, authority_email)


# #
# # Repository Provider Routes.
# #
@hug.get(
    "/repository-provider/{provider}/sign/{installation_id}/{github_repository_id}/{change_request_id}", versions=2,
)
def sign_request(
    provider: hug.types.one_of(get_supported_repository_providers().keys()),
    installation_id: hug.types.text,
    github_repository_id: hug.types.text,
    change_request_id: hug.types.text,
    request,
):
    """
    GET: /repository-provider/{provider}/sign/{installation_id}/{repository_id}/{change_request_id}

    The endpoint that will initiate a CLA signature for the user.
    """
    return cla.controllers.repository_service.sign_request(
        provider, installation_id, github_repository_id, change_request_id, request
    )


@hug.get("/repository-provider/{provider}/oauth2_redirect", versions=2)
def oauth2_redirect(
    auth_user: check_auth,  # pylint: disable=too-many-arguments
    provider: hug.types.one_of(get_supported_repository_providers().keys()),
    state: hug.types.text,
    code: hug.types.text,
    repository_id: hug.types.text,
    change_request_id: hug.types.text,
    request=None,
):
    """
    GET: /repository-provider/{provider}/oauth2_redirect

    TODO: This has been deprecated in favor of GET:/github/installation for GitHub Apps.

    Handles the redirect from an OAuth2 provider when initiating a signature.
    """
    # staff_verify(user)
    return cla.controllers.repository_service.oauth2_redirect(
        provider, state, code, repository_id, change_request_id, request
    )


@hug.post("/repository-provider/{provider}/activity", versions=2)
def received_activity(body, provider: hug.types.one_of(get_supported_repository_providers().keys())):
    """
    POST: /repository-provider/{provider}/activity

    TODO: Need to secure this endpoint somehow - maybe use GitHub's Webhook secret option.

    Acts upon a code repository provider's activity.
    """
    return cla.controllers.repository_service.received_activity(provider, body)


#
# GitHub Routes.
#
@hug.get("/github/organizations", versions=1)
def get_github_organizations(auth_user: check_auth):
    """
    GET: /github/organizations

    Returns all CLA Github Organizations.
    """
    return cla.controllers.github.get_organizations()


@hug.get("/github/organizations/{organization_name}", versions=1)
def get_github_organization(auth_user: check_auth, organization_name: hug.types.text):
    """
    GET: /github/organizations/{organization_name}

    Returns the CLA Github Organization requested by Name.
    """
    return cla.controllers.github.get_organization(organization_name)


@hug.get("/github/organizations/{organization_name}/repositories", versions=1)
def get_github_organization_repos(auth_user: check_auth, organization_name: hug.types.text):
    """
    GET: /github/organizations/{organization_name}/repositories

    Returns a list of Repositories selected under this organization.
    """
    return cla.controllers.github.get_organization_repositories(organization_name)


@hug.get("/sfdc/{sfid}/github/organizations", versions=1)
def get_github_organization_by_sfid(auth_user: check_auth, sfid: hug.types.text):
    """
    GET: /github/organizations/sfdc/{sfid}

    Returns a list of Github Organizations under this SFDC ID.
    """
    return cla.controllers.github.get_organization_by_sfid(auth_user, sfid)


@hug.post(
    "/github/organizations",
    versions=1,
    examples=" - {'organization_sfid': '<organization-sfid>', \
                        'organization_name': 'org-name'}",
)
def post_github_organization(
    auth_user: check_auth,  # pylint: disable=too-many-arguments
    organization_name: hug.types.text,
    organization_sfid: hug.types.text,
):
    """
    POST: /github/organizations

    DATA: { 'auth_user' : AuthUser to verify user permissions
            'organization_sfid': '<sfid-id>',
            'organization_name': 'org-name'}

    Returns the CLA GitHub Organization that was just created.
    """
    return cla.controllers.github.create_organization(auth_user, organization_name, organization_sfid)


@hug.delete("/github/organizations/{organization_name}", versions=1)
def delete_organization(auth_user: check_auth, organization_name: hug.types.text):
    """
    DELETE: /github/organizations/{organization_name}

    Deletes the specified Github Organization.
    """
    # staff_verify(user)
    return cla.controllers.github.delete_organization(auth_user, organization_name)


@hug.get("/github/installation", versions=2)
def github_oauth2_callback(code, state, request):
    """
    GET: /github/installation

    TODO: Need to secure this endpoint - possibly with GitHub's Webhook secrets.

    GitHub will send the user to this endpoint when new OAuth2 handshake occurs.
    This needs to match the callback used when users install the app as well (below).
    """
    return cla.controllers.github.user_oauth2_callback(code, state, request)


@hug.post("/github/installation", versions=2)
def github_app_installation(body, request, response):
    """
    POST: /github/installation

    TODO: Need to secure this endpoint - possibly with GitHub's Webhook secret.

    GitHub will fire off this webhook when new installation of our CLA app occurs.
    """
    return cla.controllers.github.user_authorization_callback(body)


@hug.post("/github/activity", versions=2)
def github_app_activity(body, request, response):
    """
    POST: /github/activity

    TODO: Need to secure this endpoint with GitHub's Webhook secret.

    Acts upon any events triggered by our app installed in someone's organization.
    """
    # Verify that Webhook Signature is valid
    # valid_request = cla.controllers.github.webhook_secret_validation(request.headers.get('X-HUB-SIGNATURE'), request.stream.read())
    # cla.log.info(valid_request)
    # if valid_request:
    event_type = request.headers.get('X-GITHUB-EVENT')
    if event_type is None:
        response.status = HTTP_400
        return {'status': 'Invalid request'}

    return cla.controllers.github.activity(event_type, body)
    # else:
    #     response.status = HTTP_403
    #     return {'status': 'Not Authorized'}


@hug.post("/github/validate", versions=1)
def github_organization_validation(body):
    """
    POST: /github/validate

    TODO: Need to secure this endpoint with GitHub's Webhook secret.
    """
    return cla.controllers.github.validate_organization(body)


@hug.get("/github/check/namespace/{namespace}", versions=1)
def github_check_namespace(namespace):
    """
    GET: /github/check/namespace/{namespace}

    Returns True if the namespace provided is a valid GitHub account.
    """
    return cla.controllers.github.check_namespace(namespace)


@hug.get("/github/get/namespace/{namespace}", versions=1)
def github_get_namespace(namespace):
    """
    GET: /github/get/namespace/{namespace}

    Returns info on the GitHub account provided.
    """
    return cla.controllers.github.get_namespace(namespace)


#
# Gerrit instance routes
#
@hug.get("/project/{project_id}/gerrits", versions=1)
def get_project_gerrit_instance(project_id: hug.types.uuid):
    """
    GET: /project/{project_id}/gerrits

    Returns all CLA Gerrit instances for this project.
    """
    return cla.controllers.gerrit.get_gerrit_by_project_id(project_id)


@hug.get("/gerrit/{gerrit_id}", versions=2)
def get_gerrit_instance(gerrit_id: hug.types.uuid):
    """
    GET: /gerrit/gerrit_id

    Returns Gerrit instance with the given gerrit id.
    """
    return cla.controllers.gerrit.get_gerrit(gerrit_id)


@hug.post("/gerrit", versions=1)
def create_gerrit_instance(
    project_id: hug.types.uuid,
    gerrit_name: hug.types.text,
    gerrit_url: cla.hug_types.url,
    group_id_icla=None,
    group_id_ccla=None,
):
    """
    POST: /gerrit

    Creates a gerrit instance
    """
    return cla.controllers.gerrit.create_gerrit(project_id, gerrit_name, gerrit_url, group_id_icla, group_id_ccla)


@hug.delete("/gerrit/{gerrit_id}", versions=1)
def delete_gerrit_instance(gerrit_id: hug.types.uuid):
    """
    DELETE: /gerrit/{gerrit_id}

    Deletes the specified gerrit instance.
    """
    return cla.controllers.gerrit.delete_gerrit(gerrit_id)


@hug.get(
    "/gerrit/{gerrit_id}/{contract_type}/agreementUrl.html", versions=2, output=hug.output_format.html,
)
def get_agreement_html(gerrit_id: hug.types.uuid, contract_type: hug.types.text):
    """
    GET: /gerrit/{gerrit_id}/{contract_type}/agreementUrl.html

    Generates an appropriate HTML file for display in the Gerrit console.
    """
    return cla.controllers.gerrit.get_agreement_html(gerrit_id, contract_type)


# The following routes are only provided for project and cla manager
# permission management, and are not to be called by the UI Consoles.
@hug.get("/project/logo/{project_sfdc_id}", versions=1)
def upload_logo(auth_user: check_auth, project_sfdc_id: hug.types.text):
    return cla.controllers.project_logo.create_signed_logo_url(auth_user, project_sfdc_id)


@hug.post("/project/permission", versions=1)
def add_project_permission(auth_user: check_auth, username: hug.types.text, project_sfdc_id: hug.types.text):
    return cla.controllers.project.add_permission(auth_user, username, project_sfdc_id)


@hug.delete("/project/permission", versions=1)
def remove_project_permission(auth_user: check_auth, username: hug.types.text, project_sfdc_id: hug.types.text):
    return cla.controllers.project.remove_permission(auth_user, username, project_sfdc_id)


@hug.post("/company/permission", versions=1)
def add_company_permission(auth_user: check_auth, username: hug.types.text, company_id: hug.types.text):
    return cla.controllers.company.add_permission(auth_user, username, company_id)


@hug.delete("/company/permission", versions=1)
def remove_company_permission(auth_user: check_auth, username: hug.types.text, company_id: hug.types.text):
    return cla.controllers.company.remove_permission(auth_user, username, company_id)


@hug.get("/events", versions=1)
def search_events(request, response):
    return cla.controllers.event.events(request, response)


@hug.get("/events/{event_id}", versions=1)
def get_event(event_id: hug.types.text, response):
    """
    GET: /event/{event_id}
    Returns the requested event data based on ID
    """
    return cla.controllers.event.get_event(event_id=event_id, response=response)


@hug.post("/events", versions=1)
def create_event(
    event_data: hug.types.text,
    event_type: hug.types.text = None,
    user_id: hug.types.text = None,
    event_project_id: hug.types.text = None,
    event_company_id: hug.types.text = None,
    response=None,
):
    return cla.controllers.event.create_event(
        response=response,
        event_type=event_type,
        event_company_id=event_company_id,
        event_data=event_data,
        event_project_id=event_project_id,
        user_id=user_id,
    )


# Session Middleware
__hug__.http.add_middleware(get_session_middleware())
__hug__.http.add_middleware(LogMiddleware(logger=cla.log))
