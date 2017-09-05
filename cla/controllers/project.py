"""
Controller related to project operations.
"""

import cla
from cla.utils import get_project_instance, get_document_instance
from cla.models import DoesNotExist


def get_projects():
    """
    Returns a list of projects in the CLA system.

    :return: List of projects in dict format.
    :rtype: [dict]
    """
    return [project.to_dict() for project in get_project_instance().all()]


def get_project(project_id):
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


def get_project_by_external_id(project_external_id):
    """
    Returns the CLA project requested by External ID.

    :param project_external_id: The project's External ID.
    :type project_external_id: string
    :return: dict representation of the project object.
    :rtype: dict
    """
    project = get_project_instance()
    try:
        project.load(project_external_id=str(project_external_id))
    except DoesNotExist as err:
        return {'errors': {'project_external_id': str(err)}}
    return project.to_dict()


def create_project(project_id, project_name=None):
    """
    Creates a project and returns the newly created project in dict format.

    :param project_id: The project's given ID.
    :type project_id: string
    :param project_name: The project's name.
    :type project_name: string
    :return: dict representation of the project object.
    :rtype: dict
    """
    project = get_project_instance()
    project.set_project_id(str(project_id))
    project.set_project_name(project_name)
    project.save()
    return project.to_dict()


def update_project(project_id, project_name=None):
    """
    Updates a project and returns the newly updated project in dict format.
    A value of None means the field should not be updated.

    :param project_id: ID of the project to update.
    :type project_id: string
    :param project_name: New project name.
    :type project_name: string | None
    :return: dict representation of the project object.
    :rtype: dict
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    if project_name is not None:
        project.set_project_name(project_name)
    project.save()
    return project.to_dict()


def delete_project(project_id):
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
    project.delete()
    return {'success': True}


def get_project_repositories(project_id):
    """
    Get a project's repositories.

    :param project_id: The ID of the project.
    :type project_id: string
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    repositories = project.get_project_repositories()
    return [repository.to_dict() for repository in repositories]


def get_project_document(project_id, document_type, revision=None):
    """
    Returns the specified project's document based on type (ICLA or CCLA) and revision.

    :param project_id: The ID of the project to fetch the document from.
    :type project_id: string
    :param document_type: The type of document (individual or corporate).
    :type document_type: string
    :param revision: The revision number of the document to fetch.
    :type revision: integer
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    if document_type == 'individual':
        try:
            document = project.get_project_individual_document(revision)
        except DoesNotExist as err:
            return {'errors': {'document': str(err)}}
    else:
        try:
            document = project.get_project_corporate_document(revision)
        except DoesNotExist as err:
            return {'errors': {'document': str(err)}}
    return document.to_dict()


def post_project_document(project_id,
                          document_type,
                          document_name,
                          document_content_type,
                          document_content):
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
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    document = get_document_instance()
    document.set_document_name(document_name)
    document.set_document_content_type(document_content_type)
    document.set_document_content(document_content)
    if document_type == 'individual':
        revision = cla.utils.get_last_revision(project.get_project_individual_documents())
        document.set_document_revision(revision + 1)
        project.add_project_individual_document(document)
    else:
        revision = cla.utils.get_last_revision(project.get_project_corporate_documents())
        document.set_document_revision(revision + 1)
        project.add_project_corporate_document(document)
    project.save()
    return project.to_dict()


def delete_project_document(project_id, document_type, revision):
    """
    Deletes the document from the specified project.

    :param project_id: The ID of the project in question.
    :type project_id: string
    :param document_type: The type of document to remove (individual or corporate).
    :type document_type: string
    :param revision: The document revision number to remove.
    :type revision: integer
    """
    project = get_project_instance()
    try:
        project.load(str(project_id))
    except DoesNotExist as err:
        return {'errors': {'project_id': str(err)}}
    document = cla.utils.get_project_document(project, document_type, revision)
    if document is None:
        return {'errors': {'document': 'Document revision not found'}}
    if document_type == 'individual':
        project.remove_project_individual_document(document)
    else:
        project.remove_project_corporate_document(document)
    project.save()
    return {'success': True}
