"""
Controller related to user operations.
"""

import uuid
import hug
from cla.utils import get_user_instance
from cla.models import DoesNotExist
import cla

def get_users():
    """
    Returns a list of users in the CLA system.

    :return: List of users in dict format.
    :rtype: [dict]
    """
    return [user.to_dict() for user in get_user_instance().all()]

def get_user(user_id=None, user_email=None, user_github_id=None):
    """
    Returns the CLA user requested by ID or email.

    :param user_id: The user's ID.
    :type user_id: string
    :param user_email: The user's email address.
    :type user_email: string
    :param user_github_id: The user's github ID.
    :type user_github_id: integer
    :return: dict representation of the user object.
    :rtype: dict
    """
    if user_id is not None:
        user = get_user_instance()
        try:
            user.load(user_id)
        except DoesNotExist as err:
            return {'errors': {'user_id': str(err)}}
    elif user_email is not None:
        user = get_user_instance().get_user_by_email(str(user_email).lower())
        if user is None:
            return {'errors': {'user_email': 'User not found'}}
    elif user_github_id is not None:
        user = get_user_instance().get_user_by_github_id(user_github_id)
        if user is None:
            return {'errors': {'user_github_id': 'User not found'}}
    return user.to_dict()

def create_user(user_email, user_name=None, user_organization_id=None, user_github_id=None):
    """
    Creates a user and returns it in dict format.

    :param user_email: The email address of the new user.
    :type user_email: string
    :param user_name: The name of the new user.
    :type user_name: string
    :param user_organization_id: The organization ID the user belongs to.
    :type user_organization_id: string
    :param user_github_id: The GitHub ID of the user (optional).
    :type user_github_id: integer | None
    :return: dict object containing user data.
    :rtype: dict
    """
    user = get_user_instance()
    user.set_user_id(str(uuid.uuid4()))
    user.set_user_email(str(user_email).lower())
    user.set_user_name(user_name)
    user.set_user_organization_id(user_organization_id)
    user.set_user_github_id(user_github_id)
    user.save()
    return user.to_dict()

def update_user(user_id, user_email=None, user_name=None,
                user_organization_id=None, user_github_id=None):
    """
    Updates a user and returns it in dict format.

    :param user_id: The user ID of the user to update.
    :type user_id: string
    :param user_email: The new email address for the user.
    :type user_email: string
    :param user_name: The new name for the user.
    :type user_name: string
    :param user_organization_id: The new organization ID for the user.
    :type user_organization_id: string
    :param user_github_id: The new GitHub ID of the user (optional).
    :type user_github_id: integer | None
    :return: dict object containing the updated user data.
    :rtype: dict
    """
    user = get_user_instance()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    if user_email is not None:
        try:
            val = cla.hug_types.email(user_email)
            user.set_user_email(val)
        except ValueError as err:
            return {'errors': {'user_email': 'Invalid email specified'}}
    if user_name is not None:
        user.set_user_name(str(user_name))
    if user_organization_id is not None:
        # TODO: Ensure organization_id exists.
        user.set_user_organization_id(user_organization_id)
    if user_github_id is not None:
        try:
            val = hug.types.number(user_github_id)
            user.set_user_github_id(val)
        except ValueError as err:
            return {'errors': {'user_github_id': 'Invalid GitHub ID specified'}}
    user.save()
    return user.to_dict()

def delete_user(user_id):
    """
    Deletes a user based on their ID.

    :param user_id: The ID of the user to delete.
    :type user_id: string
    """
    user = get_user_instance()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    user.delete()
    return {'success': True}

def get_user_agreements(user_id):
    """
    Given a user ID, returns the user's agreements.

    :param user_id: The user's ID.
    :type user_id: string
    :return: list of agreement data for this user.
    :rtype: [dict]
    """
    user = get_user_instance()
    try:
        user.load(user_id)
    except DoesNotExist as err:
        return {'errors': {'user_id': str(err)}}
    agreements = user.get_user_agreements()
    return [agr.to_dict() for agr in agreements]

def get_users_organization(user_organization_id):
    """
    Fetches all users that are associated with the organization specified.

    :param user_organization_id: The ID of the organization in question.
    :type user_organization_id: string
    :return: A list of user data in dict format.
    :rtype: [dict]
    """
    users = get_user_instance().get_users_by_organization(user_organization_id)
    return [user.to_dict() for user in users]
