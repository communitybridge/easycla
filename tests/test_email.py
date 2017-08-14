"""
Tests having to do with the agreement endpoints.
"""

import unittest
import copy

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
from cla.models.smtp_models import MockSMTP
from cla.models.ses_models import MockSES

class EmailTestCase(CLATestCase):
    """Email test cases."""
    def test_send_email(self):
        """Tests for sending email through email models."""
        models = [MockSMTP, MockSES]
        for model in models:
            email_service = model()
            email_service.initialize(cla.conf)
            email_service.emails_sent = []
            email_service.send('Test Subject1', 'Test body1', 'cla-test@mailinator.com')
            document_path = 'resources/test.pdf'
            attachment = {'type': 'file',
                          'file': document_path,
                          'filename': 'document.pdf'}
            email_service.send('Test Subject2',
                               'Test body2',
                               'cla-test@mailinator.com',
                               attachment)
            attachment = {'type': 'content',
                          'content': open(document_path, 'rb').read(),
                          'content-type': 'application/pdf',
                          'filename': 'document2.pdf'}
            email_service.send('Test Subject3',
                               'Test body3',
                               'cla-test@mailinator.com',
                               attachment)
            self.assertTrue(len(email_service.emails_sent) == 3)
            email_service.emails_sent = []

if __name__ == '__main__':
    unittest.main()
