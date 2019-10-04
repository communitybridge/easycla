# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Controller related to the github application (CLA GitHub App).
"""
import hmac
import os
from pprint import pprint
from typing import Optional

import requests

import cla
from cla.auth import AuthUser
from cla.controllers.github_application import GitHubInstallation
from cla.controllers.project import check_user_authorization
from cla.models import DoesNotExist
from cla.models.dynamo_models import UserPermissions, Repository
from cla.utils import get_github_organization_instance, get_repository_service, get_oauth_client


def get_organizations():
    """
    Returns a list of github organizations in the CLA system.

    :return: List of github organizations in dict format.
    :rtype: [dict]
    """
    return [github_organization.to_dict() for github_organization in get_github_organization_instance().all()]


def get_organization(organization_name):
    """
    Returns the CLA github organization requested by Name.

    :param organization_name: The github organization Name.
    :type organization_name: Name
    :return: dict representation of the github organization object.
    :rtype: dict
    """
    github_organization = get_github_organization_instance()
    try:
        cla.log.debug(f'Loading GitHub by organization name: {organization_name}..')
        github_organization.load(str(organization_name))
        cla.log.debug(f'Loaded GitHub by organization name: {github_organization}')
    except DoesNotExist as err:
        cla.log.warning(f'organization name {organization_name} does not exist')
        return {'errors': {'organization_name': str(err)}}
    return github_organization.to_dict()


def create_organization(auth_user,
                        organization_name,
                        organization_sfid):
    """
    Creates a github organization and returns the newly created github organization in dict format.

    :param auth_user: authorization for this user.
    :type auth_user: AuthUser
    :param organization_name: The github organization name.
    :type organization_name: string 
    :param organization_sfid: The SFDC ID for the github organization. 
    :type organization_sfid: string/None
    :return: dict representation of the new github organization object.
    :rtype: dict
    """
    # Validate user is authorized for this SFDC ID. 
    can_access = check_user_authorization(auth_user, organization_sfid)
    if not can_access['valid']:
        return can_access['errors']

    github_organization = get_github_organization_instance()
    try:
        github_organization.load(str(organization_name))
    except DoesNotExist as err:
        cla.log.debug('creating organization: {} with sfid: {}'.format(organization_name, organization_sfid))
        github_organization.set_organization_name(str(organization_name))
        github_organization.set_organization_sfid(str(organization_sfid))
        github_organization.save()
        return github_organization.to_dict()

    cla.log.warning('organization already exists: {} - unable to create'.format(organization_name))
    return {'errors': {'organization_name': 'This organization already exists'}}


def update_organization(organization_name,  # pylint: disable=too-many-arguments
                        organization_sfid=None,
                        organization_installation_id=None):
    """
    Updates a github organization and returns the newly updated org in dict format.
    Values of None means the field will not be updated.

    :param organization_name: The github organization name.
    :type organization_name: string
    :param organization_sfid: The SFDC identifier ID for the organization.
    :type organization_sfid: string/None
    :param organization_installation_id: The github app installation id.
    :type organization_installation_id: string/None
    :return: dict representation of the new github organization object.
    :rtype: dict
    """

    github_organization = get_github_organization_instance()
    try:
        github_organization.load(str(organization_name))
    except DoesNotExist as err:
        cla.log.warning('organization does not exist: {} - unable to update'.format(organization_name))
        return {'errors': {'repository_id': str(err)}}

    github_organization.set_organization_name(organization_name)
    if organization_installation_id:
        github_organization.set_organization_installation_id(organization_installation_id)
    if organization_sfid:
        github_organization.set_organization_sfid(organization_sfid)

    github_organization.save()
    cla.log.debug('updated organization: {}'.format(organization_name))
    return github_organization.to_dict()


def delete_organization(auth_user, organization_name):
    """
    Deletes a github organization based on Name.

    :param organization_name: The Name of the github organization.
    :type organization_name: Name
    """
    # Retrieve SFDC ID for this organization 
    github_organization = get_github_organization_instance()
    try:
        github_organization.load(str(organization_name))
    except DoesNotExist as err:
        cla.log.warning('organization does not exist: {} - unable to delete'.format(organization_name))
        return {'errors': {'organization_name': str(err)}}

    organization_sfid = github_organization.get_organization_sfid()

    # Validate user is authorized for this SFDC ID. 
    can_access = check_user_authorization(auth_user, organization_sfid)
    if not can_access['valid']:
        return can_access['errors']

    # Find all repositories that are under this organization 
    repositories = Repository().get_repositories_by_organization(organization_name)
    for repository in repositories:
        repository.delete()
    github_organization.delete()
    return {'success': True}


def user_oauth2_callback(code, state, request):
    github = get_repository_service('github')
    return github.oauth2_redirect(state, code, request)


def user_authorization_callback(body):
    return {'status': 'nothing to do here.'}


def get_org_name_from_event(body: dict) -> Optional[str]:
    """
    Attempts to extract the organization name from the installation created event.

    :param body: the github installation created event body
    :return: returns either the organization name or None
    """
    try:
        # Webhook event payload
        # see: https://developer.github.com/v3/activity/events/types/#webhook-payload-example-12
        return body['installation']['account']['login']
    except KeyError:
        return None


def activity(body):
    cla.log.debug('github.activity - processing github activity callback...')

    # GitHub Application
    if 'installation' in body:
        cla.log.debug('github.activity - processing github installation activity callback...')

        # New Installations
        if 'action' in body and body['action'] == 'created':
            org_name = get_org_name_from_event(body)
            if org_name is None:
                cla.log.warning(f'Unable to determine organization name from the github create event: {body}')
                return {'status', 'GitHub installation created event malformed.'}

            cla.log.debug(f'Locating organization using name: {org_name}')
            existing = get_organization(org_name)
            if 'errors' in existing:
                cla.log.warning(f'Received github installation created event for organization: {org_name}, but '
                                'the organization is not configured in EasyCLA')
                # TODO: Need a way of keeping track of new organizations that don't have projects yet.
                return {'status': 'Github Organization must be created through the Project Management Console.'}
            elif not existing['organization_installation_id']:
                update_organization(
                    existing['organization_name'],
                    existing['organization_sfid'],
                    body['installation']['id'],
                )
                cla.log.info('github.activity - Organization enrollment completed: %s', existing['organization_name'])
                return {'status': 'Organization Enrollment Completed. CLA System is operational'}
            else:
                cla.log.info('github.activity - Organization already enrolled: %s', existing['organization_name'])
                return {'status': 'Organization already enrolled in the CLA system'}
        else:  # TODO: Handle action == 'deleted'
            pass

    # Pull Requests
    if 'pull_request' in body:
        cla.log.debug('github.activity - processing github pull_request activity callback...')

        # New PR opened
        if body['action'] == 'opened' or body['action'] == 'reopened' or body['action'] == 'synchronize':
            # Copied from repository_service.py
            provider = 'github'
            service = cla.utils.get_repository_service(provider)
            result = service.received_activity(body)
            return result


def get_organization_repositories(organization_name):
    github_organization = get_github_organization_instance()
    try:
        github_organization.load(str(organization_name))
        if github_organization.get_organization_installation_id() is not None:
            cla.log.debug('GitHub Organization ID: {}'.format(github_organization.get_organization_installation_id()))
            try:
                installation = GitHubInstallation(github_organization.get_organization_installation_id())
            except Exception as e:
                msg = ('Unable to load repositories from organization: {} ({}) due to GitHub '
                       'installation permission problem or other issue, error: {} - returning error response'.
                       format(organization_name, github_organization.get_organization_installation_id(), e))
                cla.log.warn(msg)
                return {'errors': {'organization_name': organization_name, 'error': msg}}

            if installation.repos:
                repos = []
                for repo in installation.repos:
                    repos.append(repo.full_name)
                return repos
            else:
                cla.log.debug('No repositories found for Github installation id: {}'.
                              format(github_organization.get_organization_installation_id()))
                return []
    except DoesNotExist as err:
        cla.log.warning('organization name {} does not exist, error: {}'.format(organization_name, err))
        return {'errors': {'organization_name': organization_name, 'error': str(err)}}


def get_organization_by_sfid(auth_user: AuthUser, sfid):
    # Check if user has permissions
    user_permissions = UserPermissions()
    try:
        user_permissions.load(auth_user.username)
    except DoesNotExist as err:
        cla.log.warning('user {} does not exist, error: {}'.format(auth_user.username, err))
        return {'errors': {'user does not exist': str(err)}}

    user_permissions_json = user_permissions.to_dict()

    authorized_projects = user_permissions_json.get('projects')
    if sfid not in authorized_projects:
        cla.log.warning('user {} is not authorized for this Salesforce ID: {}'.
                        format(auth_user.username, sfid))
        return {'errors': {'user is not authorized for this Salesforce ID.': str(sfid)}}

    # Get all organizations under an SFDC ID
    try:
        organizations = get_github_organization_instance().get_organization_by_sfid(sfid)
    except DoesNotExist as err:
        cla.log.warning('sfid {} does not exist, error: {}'.format(sfid, err))
        return {'errors': {'sfid': str(err)}}
    return [organization.to_dict() for organization in organizations]


def org_is_covered_by_cla(owner):
    orgs = get_organizations()
    for org in orgs:
        # Org urls have to match and full enrollment has to be completed.
        if org['organization_name'] == owner and \
                org['organization_project_id'] and \
                org['organization_installation_id']:
            cla.log.debug('org: {} with project id: {} is covered by cla'.
                          format(org['organization_name'], org['organization_project_id']))
            return True

    cla.log.debug('org: {} is not covered by cla'.format(owner))
    return False


def validate_organization(body):
    if 'endpoint' in body and body['endpoint']:
        endpoint = body['endpoint']
        r = requests.get(endpoint)

        if r.status_code == 200:
            if "http://schema.org/Organization" in r.content.decode('utf-8'):
                return {"status": "ok"}
            else:
                return {"status": "invalid"}
        elif r.status_code == 404:
            return {"status": "not found"}
        else:
            return {"status": "error"}


def webhook_secret_validation(webhook_signature, data):
    if not webhook_signature:
        return False

    sha_name, signature = webhook_signature.split('=')

    if not sha_name == 'sha1':
        return False

    mac = hmac.new(os.environ.get('GH_APP_WEBHOOK_SECRET', '').encode('utf-8'), msg=data, digestmod='sha1')
    pprint(str(mac.hexdigest()))
    pprint(str(signature))
    pprint(data)

    return True if hmac.compare_digest(mac.hexdigest(), signature) else False


def check_namespace(namespace):
    """
    Checks if the namespace provided is a valid GitHub organization.

    :param namespace: The namespace to check.
    :type namespace: string
    :return: Whether or not the namespace is valid.
    :rtype: bool
    """
    oauth = get_oauth_client()
    response = oauth.get('https://api.github.com/users/' + namespace)
    return response.ok


def get_namespace(namespace):
    """
    Gets info on the GitHub account/organization provided.

    :param namespace: The namespace to get.
    :type namespace: string
    :return: Dict of info on the account in question.
    :rtype: dict
    """
    oauth = get_oauth_client()
    response = oauth.get('https://api.github.com/users/' + namespace)
    if response.ok:
        return response.json()
    else:
        return {'errors': {'namespace': 'Invalid GitHub account namespace'}}
