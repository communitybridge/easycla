# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Controller related to project CLA Group mapping operations.
"""

from cla.models import DoesNotExist
from cla.utils import (get_project_cla_group_instance)


def get_project_cla_groups():
    """
    Returns a list of projects CLA Group mappings in the CLA system.

    :return: List of projects in dict format.
    :rtype: [dict]
    """
    return [project.to_dict() for project in get_project_cla_group_instance().all()]


def get_project_cla_group(cla_group_id):
    """
    Returns the Projects associated with the CLA Group

    :param cla_group_id: The CLA Group ID
    :type cla_group_id: string
    :return: dict representation of the project CLA Group mappings
    :rtype: dict
    """
    project = get_project_cla_group_instance()
    try:
        return project.get_by_cla_group_id(cla_group_id=str(cla_group_id))
    except DoesNotExist as err:
        return {'errors': {'cla_group_id': str(err)}}
