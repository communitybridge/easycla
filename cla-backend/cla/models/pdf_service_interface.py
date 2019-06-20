# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: AGPL-3.0-or-later

"""
Holds the PDF generator service interfaces that all PDF generators implement.
"""

class PDFService(object):
    """
    Interface to the PDF generator services.
    """

    def initialize(self, config):
        """
        This method gets called once when starting the service.

        Make use of the CLA system config as needed.

        :param config: Dictionary of all data/configuration needed to initialize the service.
        :type config: dict
        """
        raise NotImplementedError()

    def generate(self, content, external_resource=False):
        """
        Method used to generate a PDF document from HTML content.

        :param content: The HTML content (or URL) to turn into a PDF document.
        :type subject: string
        :param external_resource: Whether or not the content is a URL to an external resource.
        :type recipients: boolean
        :return: The resulting PDF binary content.
        :rtype: binary data
        """
        raise NotImplementedError()
