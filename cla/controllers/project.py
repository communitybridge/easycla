"""
Controller related to project operations.
"""

import uuid
import cla
from cla.utils import get_project_instance, get_document_instance, get_signature_instance, \
                      get_company_instance
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
    try:
        p_instance = get_project_instance()
        project = p_instance.get_project_by_external_id(str(project_external_id))
    except DoesNotExist as err:
        return {'errors': {'project_external_id': str(err)}}
    return project.to_dict()


def create_project(project_external_id, project_name, project_ccla_requires_icla_signature):
    """
    Creates a project and returns the newly created project in dict format.

    :param project_external_id: The project's external ID.
    :type project_external_id: string
    :param project_name: The project's name.
    :type project_name: string
    :param project_ccla_requires_icla_signature: Whether or not the project requires ICLA with CCLA.
    :type project_ccla_requires_icla_signature: bool
    :return: dict representation of the project object.
    :rtype: dict
    """
    project = get_project_instance()
    project.set_project_id(str(uuid.uuid4()))
    project.set_project_external_id(str(project_external_id))
    project.set_project_name(project_name)
    project.set_project_ccla_requires_icla_signature(project_ccla_requires_icla_signature)
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
    # Get all reference_ids of signatures that match the project_id AND are of type 'CCLA' AND have
    # reference_type of 'company'. Return all the companies matching those reference_ids.
    signature = get_signature_instance()
    signatures = signature.get_signatures_by_project(str(project_id),
                                                     signature_signed=True,
                                                     signature_approved=True,
                                                     signature_type='ccla',
                                                     signature_reference_type='company')
    company_ids = [signature.get_reference_id() for signature in signatures]
    company = get_company_instance()
    return [comp.to_dict() for comp in company.all(company_ids)]


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
    return document.to_dict()


def post_project_document(project_id,
                          document_type,
                          document_name,
                          document_content_type,
                          document_content,
                          new_major_version=None):
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
    :param new_major_version: Whether or not to bump up the major version.
    :type new_major_version: boolean
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
    project.save()
    return project.to_dict()


def delete_project_document(project_id, document_type, major_version, minor_version):
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
    document = cla.utils.get_project_document(project, document_type, major_version, minor_version)
    if document is None:
        return {'errors': {'document': 'Document version not found'}}
    if document_type == 'individual':
        project.remove_project_individual_document(document)
    else:
        project.remove_project_corporate_document(document)
    project.save()
    return {'success': True}
