# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Holds the model interfaces that all key-value store models must implement.
"""

class KeyValueStore(object):
    """Interface to a persistent thread-safe key-value store."""

    def get(self, key):
        """
        Abstract method to retrieve a value from the store.

        :param key: The key of the value to get.
        :type key: string
        :return: The value requested.
        :rtype: string
        """
        raise NotImplementedError()

    def exists(self, key):
        """
        Abstract method to check if a key exists in the store.

        :param key: The key to check for.
        :type key: string
        :return: Whether or not the key exists in the store.
        :rtype: boolean
        """
        raise NotImplementedError()

    def set(self, key, value):
        """
        Abstract method to set a value in the store.

        :param key: The key of the value to set.
        :type key: string
        :param value: The value to set.
        :type value: string
        """
        raise NotImplementedError()

    def delete(self, key):
        """
        Abstract method to delete a value in the store.

        Should also take care of deleting the document content from the storage service
        if the content_type field starts with 'storage+'.

        :param key: The key of the value to delete.
        :type key: string
        """
        raise NotImplementedError()
