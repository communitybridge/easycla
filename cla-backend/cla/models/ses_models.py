# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Holds the AWS SES email service that can be used to send emails.
"""

import boto3
import os
import cla
from cla.models import email_service_interface

region = os.environ.get('REGION', '')
sender_email_address = os.environ.get('SES_SENDER_EMAIL_ADDRESS', '')


class SES(email_service_interface.EmailService):
    """
    AWS SES email client model.
    """
    def __init__(self):
        self.sender_email = None
        self.region = None

    def initialize(self, config):
        self.region = region
        self.sender_email = sender_email_address

    def send(self, subject, body, recipient, attachment=None):
        msg = self.get_email_message(subject, body, self.sender_email, recipient, attachment)
        # Connect to SES.
        connection = self._get_connection()
        # Send the email.
        try:
            self._send(connection, msg)
        except Exception as err:
            cla.log.error('Error while sending AWS SES email to %s: %s', recipient, str(err))

    def _get_connection(self):
        """
        Mockable method to get a connection to the SES service.
        """
        return boto3.client('ses', region_name=self.region)

    def _send(self, connection, msg): # pylint: disable=no-self-use
        """
        Mockable send method.
        """
        connection.send_raw_email(RawMessage={'Data': msg.as_string()},
                                  Source=msg['From'],
                                  Destinations=[msg['To']])


class MockSES(SES):
    """
    Mockable AWS SES email client.
    """
    def __init__(self):
        super().__init__()
        self.emails_sent = []

    def _get_connection(self):
        return None

    def _send(self, connection, msg):
        self.emails_sent.append(msg)
