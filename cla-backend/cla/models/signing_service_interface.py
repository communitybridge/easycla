# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Holds the signing service interfaces that all signing service models must implement.
"""

class SigningService(object):
    """
    Interface to the signing service.
    """

    def initialize(self, config):
        """
        This method gets called once when starting the service.

        Make use of the CLA system config as needed.

        :param config: Dictionary of all data/configuration needed to initialize the service.
        :type config: dict
        """
        raise NotImplementedError()

    def request_individual_signature(self, project_id, user_id, return_url_type, return_url, callback_url=None):
        """
        Method that will request a new signature from the user.

        Should return a dict of {'user_id': <user_id,
                                 'project_id': <project_id>,
                                 'signature_id': <signature_id>,
                                 'sign_url': <sign_url>}

        :param project_id: The ID of the project for this signature.
        :type project_id: string
        :param user_id: The ID of the user for this signature.
        :type user_id: string
        :param return_url: The URL the user will be sent to after signing.
        :type return_url: string
        :param callback_url: The URL that will be hit by the signing provider after successful
            signature from the user.
        :type callback_url: string
        :return: All data necessary to notify the user of the signing URL.
            Should return a dict of:

                {'user_id': <user_id,
                'project_id': <project_id>,
                'signature_id': <signature_id>,
                'sign_url': <sign_url>}

        :rtype: dict
        """
        raise NotImplementedError()

    def populate_sign_url(self, signature, callback_url=None):
        """
        Method used to populate the sign_url field in the signature object provided.

        Should perform all the necessary steps so that the user simply needs to be sent to
        signature.get_signature_sign_url() and the signing process should begin. Modifies the
        signature object in place.

        Should NOT save the signature object.

        :param signature: The Signature object in question.
        :type signature: cla.models.model_interfaces.Signature
        :param callback_url: The URL that will be hit by the signing provider upon successful
            signature. Will be used to update pull requests, merge requests, etc.
        :type callback_url: string
        """
        raise NotImplementedError()

    def signed_callback(self, content, repository_id, change_request_id):
        """
        Method that will handle the data that has been POSTed by the signing service
        as a callback from a successful signature.

        Should handle things like updating the signature object to 'signed'.

        :param content: The POST body of the callback.
        :type content: string
        :param repository_id: The ID of the repository that users are signing the signature for.
        :type repository_id: string
        :param change_request_id: The ID of the change request that initiated the signing prompt.
        :type change_request_id: string
        :param change_request_id: An identifier for the change request/pull request
            that this signature originated from.
        :type content: string
        """
        raise NotImplementedError()
