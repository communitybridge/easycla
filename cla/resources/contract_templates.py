"""
Holds various HTML contract templates.
"""

class ContractTemplate(object):
    def __init__(self, body=None):
        self.body = body

    def get_html_contract(self, legal_entity_name, preamble):
        html = self.body
        if html is not None:
            html = html.replace('{{legal_entity_name}}', legal_entity_name)
            html = html.replace('{{preamble}}', preamble)
        return html

class TestTemplate(ContractTemplate):
    def __init__(self, body=None):
        super().__init__(body)
        if self.body is None:
            self.body = """
<html>
    <body>
        <h3 class="legal-entity-name">
            {{legal_entity_name}}
        </h3>
        <div class="preamble">
            {{preamble}}
        </div>
    </body>
</html>"""
