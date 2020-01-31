# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT
import logging
import unittest
from unittest.mock import Mock

import json
import cla
from cla.controllers.signature import notify_whitelist_change
from cla.models.sns_email_models import MockSNS
from cla.models.dynamo_models import User, Signature, Project
from cla.user import CLAUser


class TestSignatureController(unittest.TestCase):
    def test_notify_whitelist_change(self):
        old_sig = Signature()
        new_sig = Signature()
        new_sig.set_signature_reference_name('Company')
        new_sig.set_signature_project_id('projectID')
        cla_manager = CLAUser({'name':'CLA Manager'})
        old_sig.set_domain_whitelist(['a.com','b.com'])
        new_sig.set_domain_whitelist(['b.com','d.com'])

        old_sig.set_github_whitelist([])
        new_sig.set_github_whitelist(['githubuser'])

        old_sig.set_email_whitelist(['mymail@gmail.com'])
        new_sig.set_email_whitelist([])

        old_sig.set_github_org_whitelist(['githuborg'])
        new_sig.set_github_org_whitelist(['githuborg'])

        snsClient = MockSNS()
        cla.controllers.signature.get_email_service = Mock()
        cla.controllers.signature.get_email_service.return_value = snsClient
        new_sig.get_managers = Mock(side_effect=mock_get_managers)

        cla.models.dynamo_models.Project.load = Mock(side_effect=mock_project)
        notify_whitelist_change(cla_manager, old_sig,new_sig)
        msg1 = snsClient.emails_sent[0]
        msg1 = json.loads(msg1)
        self.assertEqual(msg1['data']['subject'],'EasyCLA whitelist modified')
        self.assertEqual(msg1['data']['recipients'],['cla_manager1@gmail.com','cla_manager2@gmail.com'])
        body = msg1['data']['body']
        self.assertIn('a.com',body)
        self.assertNotIn('b.com',body)
        self.assertIn('d.com',body)
        self.assertIn('githubuser',body)
        self.assertIn('mymail@gmail.com',body)
        self.assertNotIn('githuborg',body)

def mock_get_managers():
    u1 = User()
    u1.set_lf_email('cla_manager1@gmail.com')
    u2 = User()
    u2.set_lf_email('cla_manager2@gmail.com')
    return [u1,u2]

def mock_project(project_id):
    p = Project()
    p.set_project_name('A project')

if __name__ == '__main__':
    unittest.main()
