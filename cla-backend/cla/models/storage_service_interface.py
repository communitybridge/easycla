# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Holds the storage service interfaces that all storage mechanisms must implement.
"""

class StorageService(object):
    """
    Interface to the storage services.
    """

    def initialize(self, config):
        """
        This method gets called once when starting the service.

        Make use of the CLA system config as needed.

        :param config: Dictionary of all data/configuration needed to initialize the service.
        :type config: dict
        """
        raise NotImplementedError()

    def store(self, filename, data):
        """
        Used to save file content to the storage provider.

        :param filename: The filename to save to the storage provider.
        :type filename: string
        :param data: The filename content binary data to store.
        :type data: binary data
        """
        raise NotImplementedError()

    def retrieve(self, filename):
        """
        Given a filanem, will retrieve the associated content.

        ;param filename: The filename in question.
        :type filename: string
        :return: The file content retrieved from the storage provider.
        :rtype: binary data
        """
        raise NotImplementedError()

    def delete(self, filename):
        """
        Given a filename, will delete the associated content.

        ;param filename: The filename in question.
        :type filename: string
        """
        raise NotImplementedError()
