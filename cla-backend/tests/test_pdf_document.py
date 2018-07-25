"""
Tests having to do with the PDF document generation.
"""

import unittest
import uuid
import hug

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla

class PDFTestCase(CLATestCase):
    """PDF test cases."""
    def test_pdf_generation(self):
        """Tests for generating PDF documents."""

        return # Disabled by default as we don't want to call out to DocRaptor every time.

        project = self.create_project('Name Here')
        project_id = project['project_id']
        self.create_document(project_id)
        template = cla.resources.contract_templates.CNCFTemplate('Individual', 1, 0 )
        content = template.get_html_contract('Cloud Native Computing Foundation', '')
        pdf_generator = cla.utils.get_pdf_service()
        pdf_content = pdf_generator.generate(content)
        f = open('/tmp/test.pdf', 'wb')
        f.write(pdf_content)
        f.close()

if __name__ == '__main__':
    unittest.main()
