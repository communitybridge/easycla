# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Storage service that stores files locally on disk.
"""

import os
import cla
from cla.models import storage_service_interface

class LocalStorage(storage_service_interface.StorageService):
    """
    Store documents locally.
    """
    def __init__(self):
        self.folder = None

    def initialize(self, config):
        self.folder = config['LOCAL_STORAGE_FOLDER']
        if not os.path.exists(self.folder):
            cla.log.info('Local storage folder does not exist, creating: %s', self.folder)
            os.makedirs(self.folder)

    def store(self, filename, data):
        cla.log.info('Storing filename content locally: %s', filename)
        path = self.folder + '/' + filename
        try:
            fhandle = open(path, 'wb')
            fhandle.write(data)
            fhandle.close()
        except Exception as err:
            cla.log.error('Could not save filename %s at %s: %s', filename, path, str(err))

    def retrieve(self, filename):
        cla.log.info('Retrieving filename content from local disk: %s', filename)
        path = self.folder + '/' + filename
        try:
            fhandle = open(path, 'rb')
            data = fhandle.read()
            fhandle.close()
        except FileNotFoundError:
            cla.log.error('Could not find filename content for %s: %s', filename, path)
            return None
        except Exception as err:
            cla.log.error('Could not load filename content for %s (%s): %s',
                          filename, path, str(err))
            return None
        return data

    def delete(self, filename):
        cla.log.info('Deleting filename content from local disk: %s', filename)
        path = self.folder + '/' + filename
        try:
            os.remove(path)
        except FileNotFoundError:
            cla.log.error('Could not delete filename content for %s: %s', filename, path)
        except Exception as err:
            cla.log.error('Error while deleting filename content for %s (%s): %s',
                          filename, path, str(err))
        return None
