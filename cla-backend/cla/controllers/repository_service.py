# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Controller related to repository service provider activity.
"""

import cla


def received_activity(provider, data):
    """
    Handles receiving webhook activity from the repository provider.

    Forwards the data to the appropriate provider.
    """
    service = cla.utils.get_repository_service(provider)
    result = service.received_activity(data)
    return result


def change_icon(provider, signed=False):
    """
    Properly triages the change icon request to the appropriate provider.
    """
    return cla.utils.change_icon(provider, signed)


def oauth2_redirect(provider, state, code, installation_id, github_repository_id, change_request_id, request): # pylint: disable=too-many-arguments
    """
    Properly triages the OAuth2 redirect to the appropriate provider.
    """
    service = cla.utils.get_repository_service(provider)
    return service.oauth2_redirect(state, code, installation_id, github_repository_id, change_request_id, request)


def sign_request(provider, installation_id, github_repository_id, change_request_id, request):
    """
    Properly triage the sign request to the appropriate provider.
    """
    service = cla.utils.get_repository_service(provider)
    return service.sign_request(installation_id, github_repository_id, change_request_id, request)
