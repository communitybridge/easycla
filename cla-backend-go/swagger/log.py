#!/usr/bin/python

# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MI

"""
Common logger routines.
"""
import datetime
import logging
import os

loggers = {}


def setup_custom_logger(name, log_dir: str = '.', prefix: str = 'migration'):
    """
    Sets up a custom logger using the specified name.

    :param name: the name of the logger
    :type name str
    :param log_dir: the output log directory - default is the current directory
    :type log_dir str
    :param prefix: the output log filename prefix
    :type prefix str
    :return: a custom logger using the specified name
    """

    global loggers

    if loggers.get(name):
        return loggers.get(name)
    else:
        # Returns a logger with the specified name, creating it if necessary.
        logger = logging.getLogger(name)
        logger.setLevel(logging.DEBUG)
        if logger.hasHandlers():
            # print('clearing log handlers for logger: {}'.format(name))
            logger.handlers.clear()
            if logger.parent is not None:
                # print('clearing parent log handlers for logger: {}'.format(name))
                logger.parent.handlers.clear()

        # Setup the log handlers
        formatter = logging.Formatter(
            fmt='[%(asctime)s][%(levelname)4s] - %(message)s', datefmt='%Y%m%dT%H:%M:%S%z')

        stream_handler = logging.StreamHandler()
        stream_handler.setFormatter(formatter)
        stream_handler.setLevel(logging.DEBUG)
        # stream_handler.propagate = False
        logger.addHandler(stream_handler)

        logger_filename = (f'{log_dir}{os.sep}{prefix}-'
                           f'{datetime.datetime.now().strftime("%Y%m%d-%H%M%S")}.log')
        fh = logging.FileHandler(logger_filename)
        fh.setFormatter(formatter)
        fh.setLevel(logging.DEBUG)
        # fh.propagate = False
        logger.addHandler(fh)

        loggers[name] = logger
        return logger
