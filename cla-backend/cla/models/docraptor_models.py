# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

"""
DocRaptor PDF generator.
"""

import docraptor
import os
from cla.models.pdf_service_interface import PDFService

docraptor_key = os.environ['DOCRAPTOR_API_KEY']
docraptor_test_mode = os.environ.get('DOCRAPTOR_TEST_MODE', '').lower() == 'true'

class DocRaptor(PDFService):
    """
    Implementation of the DocRaptor PDF Service.
    """
    def __init__(self):
        self.api_key = None
        self.test_mode = False
        self.javascript = False

    def initialize(self, config):
        self.api_key = docraptor_key
        docraptor.configuration.username = self.api_key
        self.debug_mode = False
        docraptor.configuration.debug = self.debug_mode
        self.test_mode = docraptor_test_mode
        self.javascript = True

    def generate(self, content, external_resource=False):
        doc_api = docraptor.DocApi()
        data = {'test': self.test_mode,
                'name': 'docraptor-python.pdf', # help you find a document later
                'document_type': 'pdf',
                'javascript': self.javascript}
        if external_resource:
            data['document_url'] = content
        else:
            data['document_content'] = content
        return doc_api.create_doc(data)

class MockDocRaptor(DocRaptor):
    """
    Mock version of the DocRaptor service.
    """
    def generate(self, content, external_resource=False):
        f = open('tests/resources/test.pdf', 'rb')
        data = f.read()
        f.close()
        return data
