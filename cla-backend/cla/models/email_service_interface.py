# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
Holds the email service interfaces that all email models must implement.
"""

from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.mime.application import MIMEApplication

class EmailService(object):
    """
    Interface to the email services.
    """

    def initialize(self, config):
        """
        This method gets called once when starting the service.

        Make use of the CLA system config as needed.

        :param config: Dictionary of all data/configuration needed to initialize the service.
        :type config: dict
        """
        raise NotImplementedError()

    def send(self, subject, body, recipient, attachment=None):
        """
        Method used to send out an email from the CLA system.

        :param subject: The subject of this email.
        :type subject: string
        :param body: The body of this email.
        :type body: string
        :param recipient: The email addresse of the recipients.
        :type recipients: string
        :param attachment: Dictionary containing the contents and content type of
            an attachment to include in the email. Example:

                {'type': 'content', 'content': <blob>,
                'content-type': 'application/pdf', 'filename': 'cla.pdf'}
                {'type': 'file', 'file': '/tmp/test.pdf',
                'content-type': 'application/pdf', 'filename': 'cla.pdf'}

            Specifying a content type and filename is optional.
        :type attachment: dict
        """
        raise NotImplementedError()

    def get_email_message(self, subject, body, sender, recipients, attachment=None): # pylint: disable=too-many-arguments
        """
        Helper method to get a prepared MIMEMultipart email message given the subject,
        body, and recipient provided.

        :param subject: The email subject
        :type subject: string
        :param body: The email body
        :type body: string
        :param sender: The sender email
        :type sender: string
        :param recipients: An array of recipient email addresses
        :type recipients: string
        :param attachment: The attachment dict (see EmailService.send() documentation).
        :type: attachment: dict
        :return: The compiled MIMEMultipart message
        :rtype: MIMEMultipart
        """
        msg = MIMEMultipart()
        msg['Subject'] = subject
        msg['From'] = sender
        if isinstance(recipients, str):
            msg['To'] = [recipients]
        else:
            msg['To'] = recipients
        # Add message body.
        part = MIMEText(body)
        msg.attach(part)
        # Add attachment.
        self.handle_email_attachment(msg, attachment)
        return msg

    def handle_email_attachment(self, msg, attachment): # pylint: disable=no-self-use
        """
        Helper method to parse the attachment and add it to the email message.

        :param msg: The email message object.
        :type msg: email.message.EmailMessage
        :param attachment: The attachment dict (see EmailService.send() documentation).
        :type: attachment: dict
        """
        if attachment is None:
            return
        content = None
        if attachment['type'] == 'content':
            content = attachment['content']
        else: # attachment['type'] == 'file':
            content = open(attachment['file'], 'rb').read()
        name = 'document.pdf'
        if 'filename' in attachment:
            name = attachment['filename']
        part = MIMEApplication(content)
        part.add_header('Content-Disposition', 'attachment', filename=name)
        msg.attach(part)
