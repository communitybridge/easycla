"""
DocRaptor PDF generator.
"""

import docraptor
from cla.models.pdf_service_interface import PDFService

class DocRaptor(PDFService):
    """
    Implementation of the DocRaptor PDF Service.
    """
    def __init__(self):
        self.api_key = None
        self.test_mode = False
        self.javascript = False

    def initialize(self, config):
        self.api_key = "zIHdSVf0zOh9bQPBGxo"
        docraptor.configuration.username = self.api_key
        self.debug_mode = True
        docraptor.configuration.debug = self.debug_mode
        self.test_mode = True
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
