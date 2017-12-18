"""
Holds the AWS SES email service that can be used to send emails.
"""

import boto3
import cla
from cla.models import email_service_interface

class SES(email_service_interface.EmailService):
    """
    AWS SES email client model.
    """
    def __init__(self):
        self.sender_email = None
        self.region = None
        self.access_key = None
        self.secret_key = None

    def initialize(self, config):
        self.region = config['SES_REGION']
        self.access_key = config['SES_ACCESS_KEY']
        self.secret_key = config['SES_SECRET_KEY']
        self.sender_email = config['SES_SENDER_EMAIL_ADDRESS']

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
        return boto3.client('ses',
                            aws_access_key_id=self.access_key,
                            aws_secret_access_key=self.secret_key)

    def _send(self, connection, msg): # pylint: disable=no-self-use
        """
        Mockable send method.
        """
        connection.send_raw_email(msg.as_string(),
                                  source=msg['From'],
                                  destinations=[msg['To']])

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
