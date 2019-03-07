"""
Controller related to repository operations.
"""

import uuid
import cla.hug_types
from cla.utils import get_repository_instance, get_supported_repository_providers
from cla.models.dynamo_models import Project, Repository, UserPermissions, GitHubOrg
from cla.models import DoesNotExist
from cla.auth import AuthUser

def get_repositories():
    """
    Returns a list of repositories in the CLA system.

    :return: List of repositories in dict format.
    :rtype: [dict]
    """
    return [repository.to_dict() for repository in get_repository_instance().all()]


def get_repository(repository_id):
    """
    Returns the CLA repository requested by ID.

    :param repository_id: The repository ID.
    :type repository_id: ID
    :return: dict representation of the repository object.
    :rtype: dict
    """
    repository = get_repository_instance()
    try:
        repository.load(str(repository_id))
    except DoesNotExist as err:
        return {'errors': {'repository_id': str(err)}}
    return repository.to_dict()


def create_repository(auth_user: AuthUser, # pylint: disable=too-many-arguments
                      repository_project_id,
                      repository_name,
                      repository_organization_name, 
                      repository_type,
                      repository_url,
                      repository_external_id=None):
    """
    Creates a repository and returns the newly created repository in dict format.

    :param repository_project_id: The ID of the repository project.
    :type repository_project_id: string
    :param repository_name: The new repository name.
    :type repository_name: string
    :param repository_type: The new repository type ('github', 'gerrit', etc).
    :type repository_type: string
    :param repository_url: The new repository URL.
    :type repository_url: string
    :param repository_external_id: The ID of the repository from the repository provider.
    :type repository_external_id: string
    :return: dict representation of the new repository object.
    :rtype: dict
    """

    # Check that organization exists 
    github_organization = GitHubOrg()
    try:
        github_organization.load(str(repository_project_id))
    except DoesNotExist as err:
        return {'errors': {'organization_name': str(err)}}

    # Check that project is valid. 
    project = Project()
    try:
        project.load(str(repository_project_id))
    except DoesNotExist as err:
        return {'errors': {'repository_project_id': str(err)}}

    # Get SFDC project identifier
    sfdc_id = project.get_project_external_id()

    # Validate user is authorized for this project
    can_access = cla.controllers.project.check_user_authorization(auth_user, sfdc_id)
    if not can_access['valid']:
      return can_access['errors']

    repository = get_repository_instance()
    repository.set_repository_id(str(uuid.uuid4()))
    repository.set_repository_project_id(str(repository_project_id))
    repository.set_repository_sfdc_id(str(sfdc_id))
    repository.set_repository_name(repository_name)
    repository.set_repository_organization_name(repository_organization_name)
    repository.set_repository_type(repository_type)
    repository.set_repository_url(repository_url)
    if repository_external_id is not None:
        repository.set_repository_external_id(repository_external_id)
    repository.save()
    return repository.to_dict()


def update_repository(repository_id, # pylint: disable=too-many-arguments
                      repository_project_id=None,
                      repository_type=None,
                      repository_name=None,
                      repository_url=None,
                      repository_external_id=None):
    """
    Updates a repository and returns the newly updated repository in dict format.
    Values of None means the field will not be updated.

    :param repository_id: ID of the repository to update.
    :type repository_id: ID
    :param repository_project_id: ID of the repository project.
    :type repository_project_id: string
    :param repository_name: New name for the repository.
    :type repository_name: string | None
    :param repository_type: New type for repository ('github', 'gerrit', etc).
    :type repository_type: string | None
    :param repository_url: New URL for the repository.
    :type repository_url: string | None
    :param repository_external_id: ID of the repository from the service provider.
    :type repository_external_id: string
    :return: dict representation of the repository object.
    :rtype: dict
    """
    repository = Repository()
    try:
        repository.load(str(repository_id))
    except DoesNotExist as err:
        return {'errors': {'repository_id': str(err)}}
    # TODO: Ensure project_id exists.
    if repository_project_id is not None:
        repository.set_repository_project_id(str(repository_project_id))
    if repository_external_id is not None:
        repository.set_repository_external_id(repository_external_id)
    if repository_name is not None:
        repository.set_repository_name(repository_name)
    if repository_type is not None:
        supported_repo_types = get_supported_repository_providers().keys()
        if repository_type in supported_repo_types:
            repository.set_repository_type(repository_type)
        else:
            return {'errors': {'repository_type':
                               'Invalid value passed. The accepted values are: (%s)' \
                               %'|'.join(supported_repo_types)}}
    if repository_url is not None:
        try:
            val = cla.hug_types.url(repository_url)
            repository.set_repository_url(val)
        except ValueError as err:
            return {'errors': {'repository_url': 'Invalid URL specified'}}
    repository.save()
    return repository.to_dict()


def delete_repository(repository_id):
    """
    Deletes a repository based on ID.

    :param repository_id: The ID of the repository.
    :type repository_id: ID
    """
    repository = get_repository_instance()
    try:
        repository.load(str(repository_id))
    except DoesNotExist as err:
        return {'errors': {'repository_id': str(err)}}
    repository.delete()
    return {'success': True}
