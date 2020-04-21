# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Holds the AWS SNS email service that can be used to send emails.
"""

import boto3
import os
import cla
import uuid
import json
import datetime
from cla.models import email_service_interface

region = os.environ.get('REGION', '')
sender_email_address = os.environ.get('SES_SENDER_EMAIL_ADDRESS', '')
topic_arn = os.environ.get('SNS_EVENT_TOPIC_ARN', '')


class SNS(email_service_interface.EmailService):
    """
    AWS SNS email client model.
    """

    def __init__(self):
        self.region = None
        self.sender_email = None
        self.topic_arn = None

    def initialize(self, config):
        self.region = region
        self.sender_email = sender_email_address
        self.topic_arn = topic_arn

    def send(self, subject, body, recipient, attachment=None):
        msg = self.get_email_message(subject, body, self.sender_email, recipient, attachment)
        # Connect to SNS.
        connection = self._get_connection()
        # Send the email.
        try:
            self._send(connection, msg)
        except Exception as err:
            cla.log.error('Error while sending AWS SNS email to %s: %s', recipient, str(err))

    def _get_connection(self):
        """
        Mockable method to get a connection to the SNS service.
        """
        return boto3.client('sns', region_name=self.region)

    def _send(self, connection, msg):  # pylint: disable=no-self-use
        """
        Mockable send method.
        """
        connection.publish(
            TopicArn=self.topic_arn,
            Message=msg,
        )

    def get_email_message(self, subject, body, sender, recipients, attachment=None):  # pylint: disable=too-many-arguments
        """
        Helper method to get a prepared email message given the subject,
        body, and recipient provided.

        :param subject: The email subject
        :type subject: string
        :param body: The email body
        :type body: string
        :param sender: The sender email
        :type sender: string
        :param recipients: An array of recipient email addresses
        :type recipient: string
        :param attachment: The attachment dict (see EmailService.send() documentation).
        :type: attachment: dict
        :return: The json message
        :rtype: string
        """
        msg = {}
        source = {}
        data = {}

        data["body"] = body
        data["from"] = sender
        data["subject"] = subject
        data["type"] = "cla-email-event"
        if isinstance(recipients, str):
            data["recipients"] = [recipients]
        else:
            data["recipients"] = recipients

        msg["data"] = data

        source["client_id"] = "easycla-service"
        source["description"] = "EasyCLA Service"
        source["name"] = "EasyCLA Service"
        msg["source_id"] = source

        msg["id"] = str(uuid.uuid4())
        msg["type"] = "cla-email-event"
        msg["version"] = "0.1.0"
        json_string = json.dumps(msg)
        # cla.log.debug(f'Email JSON: {json_string}')
        return json_string


class MockSNS(SNS):
    """
    Mockable AWS SNS email client.
    """

    def __init__(self):
        super().__init__()
        self.emails_sent = []

    def _get_connection(self):
        return None

    def _send(self, connection, msg):
        self.emails_sent.append(msg)
