# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Holds various exception classes.
"""

# Defined exceptions.
class InvalidParameters(Exception):
    """Exception raised when invalid parameters were supplied for a query."""
    pass
class DoesNotExist(Exception):
    """Exception called when queried values don't exist."""
    pass
class MultipleResults(Exception):
    """
    Exception raised when multiple results were returned from a query that
    should only have one matching result.
    """
    pass
