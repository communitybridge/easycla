# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
CLA-specific global variables and configuration.
"""

import os
import sys
import logging
import importlib
from cla import config

# Current version.
__version__ = '0.2.4'


class Config(dict):
    """
    A simple configuration object with dictionary-like properties.
    """
    def __init__(self, instance_config='cla_config'):  # pylint: disable=super-init-not-called
        """
        Initialize config object and load up default configuration file.
        """
        self.from_module(config)
        # Attempt to load the instance-specific configuration file.
        try:
            i = importlib.import_module(instance_config)
            self.from_module(i)
        except ImportError:
            logging.info('Could not load instance configuration from file: %s.py', instance_config)

    def from_module(self, mod):
        """
        Load up attributes from a module as configuration items.

        Will ignore all attributes that are not all uppercase.
        """
        for key in dir(mod):
            # Only load up capitalized attributes.
            if key.isupper():
                self[key] = getattr(mod, key)


def get_logger(configuration):
    """
    Returns a configured logger object for the CLA.
    """
    logger = logging.getLogger('cla')
    handler = logging.StreamHandler(sys.stdout)
    handler.setFormatter(configuration['LOG_FORMAT'])
    logger.addHandler(handler)
    logger.setLevel(configuration['LOG_LEVEL'])
    return logger


# The global configuration singleton.
conf = Config()
# The global logger singleton.
log = get_logger(conf)
