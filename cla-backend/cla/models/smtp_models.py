# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Holds the STMP email service that can be used to send emails.
"""

import smtplib
import cla
from cla.models import email_service_interface


class SMTP(email_service_interface.EmailService):
    """
    Simple SMTP email client.
    """

    def __init__(self):
        self.sender_email = None
        self.host = None
        self.port = None

    def initialize(self, config):
        self.sender_email = config['SMTP_SENDER_EMAIL_ADDRESS']
        self.host = config['SMTP_HOST']
        self.port = config['SMTP_PORT']

    def send(self, subject, body, recipient, attachment=None):
        msg = self.get_email_message(subject, body, self.sender_email, [recipient], attachment)
        try:
            self._send(msg)
        except Exception as err:
            cla.log.error('Error while sending STMP email to %s: %s', recipient, str(err))

    def _send(self, msg):
        """
        Mockable send method.
        """
        smtp_client = smtplib.SMTP()
        smtp_client.connect(self.host, self.port)
        smtp_client.send_message(msg)
        smtp_client.quit()


class MockSMTP(SMTP):
    """
    Mockable simple SMTP email client.
    """

    def __init__(self):
        super().__init__()
        self.emails_sent = []

    def _send(self, msg):
        self.emails_sent.append(msg)
