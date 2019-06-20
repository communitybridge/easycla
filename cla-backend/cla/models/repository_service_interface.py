# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Holds the repository service interfaces that all repository models must implement.
"""

class RepositoryService(object):
    """
    Interface to the repository services.
    """

    def initialize(self, config):
        """
        This method gets called once when starting the service.

        Make use of the CLA system config as needed.

        :param config: Dictionary of all data/configuration needed to initialize the service.
        :type config: dict
        """
        raise NotImplementedError()

    def received_activity(self, data):
        """
        Method that will be called when the repository service fires a webhook.

        Will perform various tasks in order to ensure the CLA constraints are
        being applied for the user on the repository service.

        :param data: The data provided by the webhook.
        :type data: Depends on service
        :return: A response dictionary to the service provider.
        :rtype: dict
        """
        raise NotImplementedError()

    def sign_request(self, repository_id, change_request_id, request):
        """
        Method called when the user is requesting to sign a CLA.

        :param repository_id: The ID of the repository in question.
        :type repository: string
        :param change_request_id: The ID of the change request for this signature signature.
            For GitHub, this would be the pull request number that initiated the signature.
        :type change_request_id: string
        :param request: The hug request object.
        :type request: Request
        """
        raise NotImplementedError()

    def update_change_request(self, installation_id, github_repository_id, change_request_id):
        """
        This method should handle updating the pull request/change request in the
        repository provider in order to mirror the state of the signatures in the CLA DB.

        Will be called on change request creation/update, and after a user signs an signature.
        For GitHub, this should handle creating/updating the comment/status.

        :TODO: Update comments.

        :param repository: The Repository in question.
        :type repository: cla.models.model_interfaces.Repository
        :param change_request_id: Parameter to identify the change request/pull request in question.
        :type change_request_id: string
        """
        raise NotImplementedError()

    def get_return_url(self, repository_id, change_request_id):
        """
        Method meant to be overriden by the repository provider which will return
        the URL the user should be redirected to upon signature signed, if signature was initiated
        from a pull request/merge request/etc, specific to the repository service provider.

        :param repository_id: The ID of the repository in question.
        :type repository_id: string
        :param change_request_id: The ID of the change request/pull request in question.
        :type change_request_id: string
        """
        raise NotImplementedError()
