# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Hug types.
"""

from email.utils import parseaddr
from urllib.parse import urlparse
import hug


def valid_email(value):
    """
    Simple function to validate an email address.
    This implementation is NOT perfect.

    :param value: The email address to validate.
    :type value: string
    :return: Whether or not the value is a valid email address.
    :rtype: boolean
    """
    return '@' in parseaddr(value)[1]


def valid_url(value):
    """
    Simple function to validate a URL.
    This implementation is NOT perfect.

    :param value: The URL to validate.
    :type value: string
    :return: Whether or not the value is a valid URL.
    :rtype: boolean
    """
    parsed_url = urlparse(value)
    return len(parsed_url.scheme) > 0 and len(parsed_url.netloc) > 0


class Email(hug.types.Text):
    """Simple hug type for email address validation."""
    def __call__(self, value):
        value = super().__call__(value)
        if not valid_email(value):
            raise ValueError('Invalid email address specified')
        return value
email = Email()


class URL(hug.types.Text):
    """Simple hug type for URL validation."""
    def __call__(self, value):
        value = super().__call__(value)
        if not valid_url(value):
            raise ValueError('Invalid URL specified')
        return value
url = URL()
