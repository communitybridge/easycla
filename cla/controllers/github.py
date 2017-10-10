"""
Controller related to the github application (CLA GitHub App).
"""
import requests
import hmac
import cla
from pprint import pprint
from cla.utils import get_github_organization_instance, get_repository_service
from cla.models import DoesNotExist
from cla.controllers.github_application import GitHubInstallation


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
        github_organization.load(str(organization_name))
    except DoesNotExist as err:
        return {'errors': {'organization_name': str(err)}}
    return github_organization.to_dict()


def create_organization(organization_name, # pylint: disable=too-many-arguments
                        organization_project_id=None,
                        organization_installation_id=None):
    """
    Creates a github organization and returns the newly created github organization in dict format.

    :param organization_project_id: The ID of the github organization project.
    :type organization_project_id: string/None
    :param organization_name: The github organization name.
    :type organization_name: string
    :param organization_installation_id: The github app installation id.
    :type organization_installation_id: string/None
    :return: dict representation of the new github organization object.
    :rtype: dict
    """
    github_organization = get_github_organization_instance()
    github_organization.set_organization_name(organization_name)
    if organization_project_id:
        github_organization.set_organization_company_id(organization_project_id)
    if organization_installation_id:
        github_organization.set_organization_installation_id(organization_installation_id)
    github_organization.save()
    return github_organization.to_dict()


def update_organization(organization_name, # pylint: disable=too-many-arguments
                        organization_project_id=None,
                        organization_installation_id=None):
    """
    Updates a github organization and returns the newly updated org in dict format.
    Values of None means the field will not be updated.

    :param organization_project_id: The ID of the github organization project.
    :type organization_project_id: string/None
    :param organization_name: The github organization name.
    :type organization_name: string
    :param organization_installation_id: The github app installation id.
    :type organization_installation_id: string/None
    :return: dict representation of the new github organization object.
    :rtype: dict
    """
    github_organization = get_github_organization_instance()
    try:
        github_organization.load(str(organization_name))
    except DoesNotExist as err:
        return {'errors': {'repository_id': str(err)}}
    github_organization.set_organization_name(organization_name)
    if organization_project_id:
        github_organization.set_organization_company_id(organization_project_id)
    if organization_installation_id:
        github_organization.set_organization_installation_id(organization_installation_id)
    github_organization.save()
    return github_organization.to_dict()


def delete_organization(organization_name):
    """
    Deletes a github organization based on Name.

    :param organization_name: The Name of the github organization.
    :type organization_name: Name
    """
    github_organization = get_github_organization_instance()
    try:
        github_organization.load(str(organization_name))
    except DoesNotExist as err:
        return {'errors': {'organization_name': str(err)}}
    github_organization.delete()
    return {'success': True}

def user_oauth2_callback(code, state, request):
    github = get_repository_service('github')
    return github.oauth2_redirect(state, code, request)

def user_authorization_callback(body):
    return {'status': 'nothing to do here.'}


def activity(body):
    # GitHub Application
    if 'installation' in body:

        # New Installations
        if 'action' in body and body['action'] == 'created':
            existing = get_organization(body['installation']['account']['login'])
            if 'errors' in existing:
                org = create_organization(
                    body['installation']['account']['login'],
                    None,
                    body['installation']['id']
                )
                return org
            elif not existing['organization_installation_id']:
                update_organization(
                    existing['organization_name'],
                    existing['organization_company_id'],
                    body['installation']['id']
                )
                return {'status': 'Organization Enrollment Completed. CLA System is operational.'}
            else:
                return {'status': 'Organization already exists in our system. Enrollment not completed.'}

    # Pull Requests
    if 'pull_request' in body:

        # Makes sure that the repo is known to us
        if org_is_covered_by_cla(body['pull_request']['head']['repo']['owner']['login']):

            # New PR opened
            if body['action'] == 'opened' or body['action'] == 'reopened':

                # Copied from repository_service.py
                provider = 'github'
                service = cla.utils.get_repository_service(provider)
                result = service.received_activity(body)
                return result

        # If the repo is not covered, post an annoying message on the Pull Request
        else:
            return {'status': 'Repo not covered under CLA System.'}


def get_organization_repositories(organization_name):
    github_organization = get_github_organization_instance()
    try:
        org = github_organization.load(str(organization_name))
        if org['organization_installation_id']:
            installation = GitHubInstallation(org['organization_installation_id'])
            if installation.repos:
                return installation.repos
            else:
                return []
    except DoesNotExist as err:
        return {'errors': {'organization_name': str(err)}}


def org_is_covered_by_cla(owner):
    orgs = get_organizations()
    for org in orgs:
        # Org urls have to match and full enrollment has to be completed.
        if org['organization_name'] == owner and \
                org['organization_project_id'] and \
                org['organization_installation_id']:
            return True

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

    mac = hmac.new(cla.conf['GITHUB_APP_WEBHOOK_SECRET'].encode('utf-8'), msg=data, digestmod='sha1')

    pprint(str(mac.hexdigest()))
    pprint(str(signature))
    pprint(data)

    return True if str(mac.hexdigest()) == str(signature) else False

def check_namespace(namespace):
    """
    Checks if the namespace provided is a valid GitHub organization.

    :param namespace: The namespace to check.
    :type namespace: string
    :return: Whether or not the namespace is valid.
    :rtype: bool
    """
    main_installation = GitHubInstallation()
    return main_installation.namespace_exists(namespace)
