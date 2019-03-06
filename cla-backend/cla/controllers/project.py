"""
Controller related to project operations.
"""

import uuid
import urllib
import io
import cla
import cla.resources.contract_templates
from cla.auth import AuthUser, admin_list
from cla.utils import get_project_instance, get_document_instance, get_signature_instance, \
                      get_company_instance, get_pdf_service, get_github_organization_instance
from cla.models import DoesNotExist
from cla.models.dynamo_models import Signature, Project, Company, UserPermissions, Repository, \
                                     GitHubOrg
from falcon import HTTPForbidden
from cla.controllers.github_application import GitHubInstallation

def check_user_authorization(auth_user: AuthUser, sfid):
    # Check if user has permissions on this project
    user_permissions = UserPermissions()
    try: 
        user_permissions.load(auth_user.username)
    except DoesNotExist as err:
        return {'valid': False, 'errors': {'errors': {'user does not exist': str(err)}}}

    user_permissions_json = user_permissions.to_dict()

    authorized_projects = user_permissions_json.get('projects')
    if sfid not in authorized_projects:
        return {'valid': False, 'errors': {'errors': {'user is not authorized for this Salesforce ID.': str(sfid)}}}

    return {'valid': True}

def get_projects():
    """
    Returns a list of projects in the CLA system.

    :return: List of projects in dict format.
    :rtype: [dict]
    """
    return [project.to_dict() for project in get_project_instance().all()]

def project_acl_verify(username, project_obj):
    if username in project_obj.get_project_acl():
        return True

    raise HTTPForbidden('Unauthorized',
        'Provided Token credentials does not have sufficient permissions to access resource')

def get_project(project_id, user_id=None):
    """
    Returns the CLA project requested by ID.

    :param project_id: The project's ID.
    :type project_id: string
    :return: dict representation of the project object.
    :rtype: dict
    """
    project = get_project_instance()
    try:
        project.load(project_id=str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    return project.to_dict()


def get_unsigned_projects_for_company(company_id):
    """
    Returns a list of projects that the company has not signed a CCLA for. 

    :param company_id: The company's ID.
    :type company_id: string
    :return: dict representation of the projects object.
    :rtype: [dict]
    """
    # Verify company is valid
    company = Company()
    try:
        company.load(company_id)
    except DoesNotExist as err:
        return {'errors': {'company_id': str(err)}}

    # get project ids that the company has signed the CCLAs for. 
    signature = Signature()
    signed_project_ids = signature.get_projects_by_company_signed(company_id)
    # from all projects, retrieve projects that are not in the signed project ids
    unsigned_projects = [project.to_dict() for project in Project().all() if project.get_project_id() not in signed_project_ids]
    return unsigned_projects
    

def get_projects_by_external_id(project_external_id, username):
    """
    Returns the CLA projects requested by External ID.

    :param project_external_id: The project's External ID.
    :type project_external_id: string
    :return: dict representation of the project object.
    :rtype: dict
    """
    try:
        p_instance = get_project_instance()
        projects = p_instance.get_projects_by_external_id(str(project_external_id), username)
    except DoesNotExist as err:
        return {'errors': {'project_external_id': str(err)}}
    return [project.to_dict() for project in projects]


def create_project(project_external_id, project_name, project_icla_enabled, project_ccla_enabled,
                   project_ccla_requires_icla_signature, project_acl_username):
    """
    Creates a project and returns the newly created project in dict format.

    :param project_external_id: The project's external ID.
    :type project_external_id: string
    :param project_name: The project's name.
    :type project_name: string
    :param project_icla_enabled: Whether or not the project supports ICLAs.
    :type project_icla_enabled: bool
    :param project_ccla_enabled: Whether or not the project supports CCLAs.
    :type project_ccla_enabled: bool
    :param project_ccla_requires_icla_signature: Whether or not the project requires ICLA with CCLA.
    :type project_ccla_requires_icla_signature: bool
    :return: dict representation of the project object.
    :rtype: dict
    """
    project = get_project_instance()
    project.set_project_id(str(uuid.uuid4()))
    project.set_project_external_id(str(project_external_id))
    project.set_project_name(project_name)
    project.set_project_icla_enabled(project_icla_enabled)
    project.set_project_ccla_enabled(project_ccla_enabled)
    project.set_project_ccla_requires_icla_signature(project_ccla_requires_icla_signature)
    project.set_project_acl(project_acl_username)
    project.save()

    return project.to_dict()


def update_project(project_id, project_name=None, project_icla_enabled=None,
                   project_ccla_enabled=None, project_ccla_requires_icla_signature=None, username=None):
    """
    Updates a project and returns the newly updated project in dict format.
    A value of None means the field should not be updated.

    :param project_id: ID of the project to update.
    :type project_id: string
    :param project_name: New project name.
    :type project_name: string | None
    :param project_icla_enabled: Whether or not the project supports ICLAs.
    :type project_icla_enabled: bool | None
    :param project_ccla_enabled: Whether or not the project supports CCLAs.
    :type project_ccla_enabled: bool | None
    :param project_ccla_requires_icla_signature: Whether or not the project requires ICLA with CCLA.
    :type project_ccla_requires_icla_signature: bool | None
    :return: dict representation of the project object.
    :rtype: dict
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    project_acl_verify(username, project)
    if project_name is not None:
        project.set_project_name(project_name)
    if project_icla_enabled is not None:
        project.set_project_icla_enabled(project_icla_enabled)
    if project_ccla_enabled is not None:
        project.set_project_ccla_enabled(project_ccla_enabled)
    if project_ccla_requires_icla_signature is not None:
        project.set_project_ccla_requires_icla_signature(project_ccla_requires_icla_signature)
    project.save()
    return project.to_dict()


def delete_project(project_id, username=None):
    """
    Deletes an project based on ID.

    :TODO: Need to also delete the documents saved with the storage provider.

    :param project_id: The ID of the project.
    :type project_id: string
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    project_acl_verify(username, project)
    project.delete()
    return {'success': True}


def get_project_companies(project_id):
    """
    Get a project's associated companies (via CCLA link).

    :param project_id: The ID of the project.
    :type project_id: string
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    # Get all reference_ids of signatures that match project_id AND are of reference type 'company'.
    # Return all the companies matching those reference_ids.
    signature = get_signature_instance()
    signatures = signature.get_signatures_by_project(str(project_id),
                                                     signature_signed=True,
                                                     signature_approved=True,
                                                     signature_reference_type='company')
    company_ids = list(set([signature.get_signature_reference_id() for signature in signatures]))
    company = get_company_instance()
    return [comp.to_dict() for comp in company.all(company_ids)]

def _get_project_document(project_id, document_type, major_version=None, minor_version=None):
    """
    See documentation for get_project_document().
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    if document_type == 'individual':
        try:
            document = project.get_project_individual_document(major_version, minor_version)
        except DoesNotExist as err:
            return {'errors': {'document': str(err)}}
    else:
        try:
            document = project.get_project_corporate_document(major_version, minor_version)
        except DoesNotExist as err:
            return {'errors': {'document': str(err)}}
    return document

def get_project_document(project_id, document_type, major_version=None, minor_version=None):
    """
    Returns the specified project's document based on type (ICLA or CCLA) and version.

    :param project_id: The ID of the project to fetch the document from.
    :type project_id: string
    :param document_type: The type of document (individual or corporate).
    :type document_type: string
    :param major_version: The major version number.
    :type major_version: integer
    :param minor_version: The minor version number.
    :type minor_version: integer
    """
    document = _get_project_document(project_id, document_type, major_version, minor_version=None)
    if isinstance(document, dict):
        return document
    return document.to_dict()

def get_project_document_raw(project_id, document_type, document_major_version=None, document_minor_version=None):
    """
    Same as get_project_document() except that it returns the raw PDF document instead.
    """
    document = _get_project_document(project_id, document_type, document_major_version, document_minor_version)
    if isinstance(document, dict):
        return document
    content_type = document.get_document_content_type()
    if content_type.startswith('url+'):
        pdf_url = document.get_document_content()
        pdf = urllib.request.urlopen(pdf_url)
    else:
        content = document.get_document_content()
        pdf = io.BytesIO(content)
    return pdf

def post_project_document(project_id,
                          document_type,
                          document_name,
                          document_content_type,
                          document_content,
                          document_preamble,
                          document_legal_entity_name,
                          new_major_version=None,
                          username=None):
    """
    Will create a new document for the project specified.

    :param project_id: The ID of the project to add this document to.
    :type project_id: string
    :param document_type: The type of document (individual or corporate).
    :type document_type: string
    :param document_name: The name of this new document.
    :type document_name: string
    :param document_content_type: The content type of this document ('pdf', 'url+pdf',
        'storage+pdf', etc).
    :type document_content_type: string
    :param document_content: The content of the document (or URL to content if content type
        starts with 'url+'.
    :type document_content: string or binary data
    :param document_preamble: The document preamble.
    :type document_preamble: string
    :param document_legal_entity_name: The legal entity name on the document.
    :type document_legal_entity_name: string
    :param new_major_version: Whether or not to bump up the major version.
    :type new_major_version: boolean
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    project_acl_verify(username, project)
    document = get_document_instance()
    document.set_document_name(document_name)
    document.set_document_content_type(document_content_type)
    document.set_document_content(document_content)
    document.set_document_preamble(document_preamble)
    document.set_document_legal_entity_name(document_legal_entity_name)
    if document_type == 'individual':
        major, minor = cla.utils.get_last_version(project.get_project_individual_documents())
        if new_major_version:
            document.set_document_major_version(major + 1)
            document.set_document_minor_version(0)
        else:
            if major == 0:
                major = 1
            document.set_document_major_version(major)
            document.set_document_minor_version(minor + 1)
        project.add_project_individual_document(document)
    else:
        major, minor = cla.utils.get_last_version(project.get_project_corporate_documents())
        if new_major_version:
            document.set_document_major_version(major + 1)
            document.set_document_minor_version(0)
        else:
            if major == 0:
                major = 1
            document.set_document_major_version(major)
            document.set_document_minor_version(minor + 1)
        project.add_project_corporate_document(document)
    project.save()
    return project.to_dict()

def post_project_document_template(project_id,
                                   document_type,
                                   document_name,
                                   document_preamble,
                                   document_legal_entity_name,
                                   template_name,
                                   new_major_version=None,
                                   username=None):
    """
    Will create a new document for the project specified, using the existing template.

    :param project_id: The ID of the project to add this document to.
    :type project_id: string
    :param document_type: The type of document (individual or corporate).
    :type document_type: string
    :param document_name: The name of this new document.
    :type document_name: string
    :param document_preamble: The document preamble.
    :type document_preamble: string
    :param document_legal_entity_name: The legal entity name on the document.
    :type document_legal_entity_name: string
    :param template_name: The name of the template object to use.
    :type template_name: string
    :param new_major_version: Whether or not to bump up the major version.
    :type new_major_version: boolean
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    project_acl_verify(username, project)
    document = get_document_instance()
    document.set_document_name(document_name)
    document.set_document_preamble(document_preamble)
    document.set_document_legal_entity_name(document_legal_entity_name)
    if document_type == 'individual':
        major, minor = cla.utils.get_last_version(project.get_project_individual_documents())
        if new_major_version:
            document.set_document_major_version(major + 1)
            document.set_document_minor_version(0)
        else:
            document.set_document_minor_version(minor + 1)
        project.add_project_individual_document(document)
    else:
        major, minor = cla.utils.get_last_version(project.get_project_corporate_documents())
        if new_major_version:
            document.set_document_major_version(major + 1)
            document.set_document_minor_version(0)
        else:
            document.set_document_minor_version(minor + 1)
        project.add_project_corporate_document(document)
    # Need to take the template, inject the preamble and legal entity name, and add the tabs.
    tmplt = getattr(cla.resources.contract_templates, template_name)
    template = tmplt(document_type=document_type.capitalize(),
                     major_version=document.get_document_major_version(),
                     minor_version=document.get_document_minor_version())
    content = template.get_html_contract(document_legal_entity_name, document_preamble)
    pdf_generator = get_pdf_service()
    pdf_content = pdf_generator.generate(content)
    document.set_document_content_type('storage+pdf')
    document.set_document_content(pdf_content, b64_encoded=False)
    document.set_raw_document_tabs(template.get_tabs())
    project.save()
    return project.to_dict()

def delete_project_document(project_id, document_type, major_version, minor_version, username=None):
    """
    Deletes the document from the specified project.

    :param project_id: The ID of the project in question.
    :type project_id: string
    :param document_type: The type of document to remove (individual or corporate).
    :type document_type: string
    :param major_version: The document major version number to remove.
    :type major_version: integer
    :param minor_version: The document minor version number to remove.
    :type minor_version: integer
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    project_acl_verify(username, project)
    document = cla.utils.get_project_document(project, document_type, major_version, minor_version)
    if document is None:
        return {'errors': {'document': 'Document version not found'}}
    if document_type == 'individual':
        project.remove_project_individual_document(document)
    else:
        project.remove_project_corporate_document(document)
    project.save()
    return {'success': True}

def add_permission(auth_user: AuthUser, username: str, project_sfdc_id: str):
    if auth_user.username not in admin_list:
        return {'error': 'unauthorized'}

    cla.log.info('project ({}) added for user ({}) by {}'.format(project_sfdc_id, username, auth_user.username))

    user_permission = UserPermissions()
    try:
        user_permission.load(username)
    except Exception as err:
        print('user not found. creating new user: {}'.format(err))
        # create new user
        user_permission = UserPermissions(username=username)

    user_permission.add_project(project_sfdc_id)
    user_permission.save()

def remove_permission(auth_user: AuthUser, username: str, project_sfdc_id: str):
    if auth_user.username not in admin_list:
        return {'error': 'unauthorized'}

    cla.log.info('project ({}) removed for ({}) by {}'.format(project_sfdc_id, username, auth_user.username))

    user_permission = UserPermissions()
    try:
        user_permission.load(username)
    except Exception as err:
        print('Unable to update user permission: {}'.format(err))
        return {'error': err}

    user_permission.remove_project(project_sfdc_id)
    user_permission.save()

def get_project_repositories(auth_user: AuthUser, project_id):
    """
    Get a project's repositories.

    :param project_id: The ID of the project.
    :type project_id: string
    """
    # Load Project
    project = Project()
    try:
        project.load(project_id=str(project_id))
    except DoesNotExist as err:
        return {'valid': False, 'errors': {'errors': {'project_id': str(err)}}}
    
    # Get SFDC project identifier
    sfid = project.get_project_external_id()

    # Validate user is authorized for this project
    can_access = check_user_authorization(auth_user, sfid)
    if not can_access['valid']:
      return can_access['errors']

    # Obtain repositories
    repositories = project.get_project_repositories()
    return [repository.to_dict() for repository in repositories]

def get_project_repositories_group_by_organization(auth_user: AuthUser, project_id):
    """
    Get a project's repositories.

    :param project_id: The ID of the project.
    :type project_id: string
    """
    # Load Project
    project = Project()
    try:
        project.load(project_id=str(project_id))
    except DoesNotExist as err:
        return {'valid': False, 'errors': {'errors': {'project_id': str(err)}}}
    
    # Get SFDC project identifier
    sfid = project.get_project_external_id()

    # Validate user is authorized for this project
    can_access = check_user_authorization(auth_user, sfid)
    if not can_access['valid']:
      return can_access['errors']

    # Obtain repositories
    repositories = project.get_project_repositories()
    repositories = [repository.to_dict() for repository in repositories]

    # Group them by organization
    organizations_dict = {}
    for repository in repositories:
        org_name = repository['repository_organization_name']
        if org_name in organizations_dict:
            organizations_dict[org_name].append(repository)
        else:
            organizations_dict[org_name] = [repository]
    
    organizations = []
    for key, value in organizations_dict.items():
        organizations.append({'name': key, 'repositories': value})

    return organizations
    

def get_project_configuration_orgs_and_repos(auth_user: AuthUser, project_id):
    # Load Project
    project = Project()
    try:
        project.load(project_id=str(project_id))
    except DoesNotExist as err:
        return {'valid': False, 'errors': {'errors': {'project_id': str(err)}}}
    
    # Get SFDC project identifier
    sfid = project.get_project_external_id()

    # Validate user is authorized for this project
    can_access = check_user_authorization(auth_user, sfid)
    if not can_access['valid']:
      return can_access['errors']

    # Obtain information for this project
    orgs_and_repos = get_github_repositories_by_org(project)
    repositories = get_sfdc_project_repositories(project)
    return {
        'orgs_and_repos': orgs_and_repos,
        'repositories': repositories
    }

def get_github_repositories_by_org(project):
    """
    Gets organization with the project_id specified and all its repositories from Github API

    :param project: The Project object
    :type project: Project
    :return: [] of organizations and its repositories
    [{
        'organization_name': ..
        ...
        'repositories': [{
            'repository_github_id': ''
            'repository_name': ''
            'repository_type': ''
            'repository_url': ''
        }]
    }]
    :rtype: array
    """

    organization_dicts = []
    # Get all organizations connected to this project
    github_organizations = GitHubOrg().get_organization_by_sfid(project.get_project_external_id())
    # Iterate over each organization
    for github_organization in github_organizations:
        installation_id = github_organization.get_organization_installation_id()
        # Verify installation_id exist
        if installation_id is not None:
            installation = GitHubInstallation(installation_id)
            # Prepare organization in dict 
            organization_dict = github_organization.to_dict()
            organization_dict['repositories'] = []
            # Get repositories from Github API
            github_repos = installation.repos
            if github_repos is not None:
                for repo in github_repos:
                    # Convert repository entities from lib to a dict.
                    repo_dict = {
                        'repository_github_id': repo.id, 
                        'repository_name': repo.full_name, 
                        'repository_type': 'github', 
                        'repository_url': repo.html_url
                    }
                    # Add repository to organization repositories list
                    organization_dict['repositories'].append(repo_dict)
            # Add organization dict to list
            organization_dicts.append(organization_dict)
    return organization_dicts

def get_sfdc_project_repositories(project):
    """
    Gets all SFDC repositories and divide them for current contract group and other contract groups
    :param project: The Project object
    :type project: Project
    :return: array of all sfdc project repositories
    :rtype: dict
    """

    # Get all SFDC Project repositories
    sfdc_id = project.get_project_external_id()
    all_project_repositories = Repository().get_repository_by_sfdc_id(sfdc_id)
    return [repo.to_dict() for repo in all_project_repositories]
