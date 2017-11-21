"""
Holds various HTML contract templates.
"""

class ContractTemplate(object):
    def __init__(self, document_type='Individual', major_version=1, minor_version=0, body=None):
        self.document_type = 'Individual'
        self.major_version = 1
        self.minor_version = 0
        self.body = body

    def get_html_contract(self, legal_entity_name, preamble):
        html = self.body
        if html is not None:
            html = html.replace('{{document_type}}', self.document_type)
            html = html.replace('{{major_version}}', str(self.major_version))
            html = html.replace('{{minor_version}}', str(self.minor_version))
            html = html.replace('{{legal_entity_name}}', legal_entity_name)
            html = html.replace('{{preamble}}', preamble)
        return html

class TestTemplate(ContractTemplate):
    def __init__(self, document_type='Individual', major_version=1, minor_version=0, body=None):
        super().__init__(body)
        if self.body is None:
            self.body = """
<html>
    <body>
        <h3 class="legal-entity-name" style="text-align: center">
            {{legal_entity_name}}<br />
            {{document_type}} Contributor License Agreement ("Agreement") v{{major_version}}.{{minor_version}}
        </h3>
        <div class="preamble">
            {{preamble}}
        </div>
        <p>If you have not already done so, please complete and sign, then scan and email a pdf file of this Agreement to cla@cncf.io.<br />If necessary, send an original signed Agreement to The Linux Foundation: 1 Letterman Drive, Building D, Suite D4700, San Francisco CA 94129, U.S.A.<br />Please read this document carefully before signing and keep a copy for your records.
        </p>
        <p>You accept and agree to the following terms and conditions for Your present and future Contributions submitted to the Foundation. In return, the Foundation shall not use Your Contributions in a way that is contrary to the public benefit or inconsistent with its nonprofit status and bylaws in effect at the time of the Contribution. Except for the license granted herein to the Foundation and recipients of software distributed by the Foundation, You reserve all right, title, and interest in and to Your Contributions
        </p>
    </body>
</html>"""
