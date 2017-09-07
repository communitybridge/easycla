"""
Tests to be run for every supported storage engine.

Ensure you import cla and set the appropriate configurations for the database
to be tested before importing this module. DatabaseTestCase should be used
as a base class.
"""

import uuid
import unittest
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
# TODO: Add document tests.
from cla.utils import get_user_instance, get_project_instance, get_document_instance, \
                      get_signature_instance, get_repository_instance, get_company_instance, \
                      create_database, delete_database
from cla.models import InvalidParameters, MultipleResults, DoesNotExist

class DatabaseTestCase(CLATestCase):
    """Database test cases."""
    def test_user_model(self):
        """Test user creation/querying/assertions."""
        user = get_user_instance()
        with self.assertRaises(DoesNotExist) as context:
            user.load(str(uuid.uuid4()))
        self.assertTrue('User not found' in str(context.exception))
        user_id = str(uuid.uuid4())
        user_email = 'test@user.com'
        user_name = 'User Name'
        user_github_id = 88888
        user = get_user_instance()
        user.set_user_id(user_id)
        user.set_user_email(user_email)
        user.set_user_name(user_name)
        user.set_user_github_id(user_github_id)
        user.save()
        user = get_user_instance()
        user.load(user_id)
        self.assertEqual(user.get_user_email(), user_email)
        self.assertEqual(user.get_user_name(), user_name)
        self.assertEqual(user.get_user_github_id(), user_github_id)

    def test_repository_model(self):
        """Test repository creation/querying/assertions."""
        project = get_project_instance()
        project.set_project_id(str(uuid.uuid4()))
        project.set_project_name('Project Name')
        repo_id = str(uuid.uuid4())
        repo_name = 'Repo Name'
        repo_type = 'mock_github'
        repo_url = 'https://some-github-url.com/repo-name'
        repo_signature_type = 'url+pdf'
        repo_signature_content = 'https://some-github-url.com/cla.pdf'
        repo = get_repository_instance()
        repo.set_repository_id(repo_id)
        repo.set_repository_project_id(project.get_project_id())
        repo.set_repository_name(repo_name)
        repo.set_repository_type(repo_type)
        repo.set_repository_url(repo_url)
        repo.save()

        # Re-load repo by ID.
        repo2 = get_repository_instance()
        repo2.load(repo_id)
        self.assertEqual(repo.get_repository_id(), repo2.get_repository_id())
        self.assertEqual(repo.get_repository_project_id(), repo2.get_repository_project_id())
        self.assertEqual(repo.get_repository_name(), repo2.get_repository_name())
        self.assertEqual(repo.get_repository_type(), repo2.get_repository_type())
        self.assertEqual(repo.get_repository_url(), repo2.get_repository_url())

        # Test repository not found query.
        repo3 = get_repository_instance()
        with self.assertRaises(DoesNotExist) as context:
            repo3.load(str(uuid.uuid4()))
        self.assertTrue('Repository not found' in str(context.exception))

        # Test loading all repositories.
        repo4 = get_repository_instance()
        repo4.set_repository_id(str(uuid.uuid4()))
        repo4.set_repository_project_id(project.get_project_id())
        repo4.set_repository_name('Repo Name 2')
        repo4.set_repository_type('github')
        repo4.set_repository_url('https://github.com/repo-name')
        repo4.save()
        repos = get_repository_instance().all()
        expected = [repo.get_repository_name() for repo in repos]
        self.assertEqual(len(expected), 2)
        self.assertTrue('Repo Name' in expected)
        self.assertTrue('Repo Name 2' in expected)
        repos = get_repository_instance().all(ids=[repo4.get_repository_id()])
        self.assertEqual(len(repos), 1)
        self.assertEqual(repos[0].get_repository_id(), repo4.get_repository_id())

        # Test JSON serialization.
        repo_data = repo.to_dict()
        self.assertEqual(repo_data['repository_id'], repo.get_repository_id())
        self.assertEqual(repo_data['repository_project_id'], repo.get_repository_project_id())
        self.assertEqual(repo_data['repository_name'], repo_name)
        self.assertEqual(repo_data['repository_type'], repo_type)
        self.assertEqual(repo_data['repository_url'], repo_url)

    def test_signature_model(self):
        """Test signature creation/querying/assertions."""
        signature_id = str(uuid.uuid4())
        project_id = self.create_project()['project_id']
        user_id = str(uuid.uuid4())
        signature_type = 'cla'
        signature_corporate = True
        signature_signed = True
        signature_approved = True
        signature = get_signature_instance()
        signature.set_signature_id(signature_id)
        signature.set_signature_project_id(project_id)
        signature.set_signature_document_major_version(1)
        signature.set_signature_document_minor_version(0)
        signature.set_signature_reference_id(user_id)
        signature.set_signature_reference_type('user')
        signature.set_signature_type(signature_type)
        signature.set_signature_signed(signature_signed)
        signature.set_signature_approved(signature_approved)
        signature.save()

        # Re-load signature by ID.
        signature2 = get_signature_instance()
        signature2.load(signature_id)
        self.assertEqual(signature.get_signature_id(), signature2.get_signature_id())
        self.assertEqual(signature.get_signature_project_id(), signature2.get_signature_project_id())
        self.assertEqual(signature.get_signature_reference_id(),
                         signature2.get_signature_reference_id())
        self.assertEqual(signature.get_signature_reference_type(),
                         signature2.get_signature_reference_type())
        self.assertEqual(signature.get_signature_type(), signature2.get_signature_type())
        self.assertEqual(signature.get_signature_signed(), signature2.get_signature_signed())
        self.assertEqual(signature.get_signature_approved(), signature2.get_signature_approved())

        # Test signature not found query.
        signature3 = get_signature_instance()
        with self.assertRaises(DoesNotExist) as context:
            signature3.load(str(uuid.uuid4()))
        self.assertTrue('Signature not found' in str(context.exception))

        # Test loading all signatures.
        signatures = get_signature_instance().all()
        expected = [signature.get_signature_id() for signature in signatures]
        self.assertEqual(len(expected), 1)
        self.assertTrue(signature.get_signature_id() in expected)
        signatures = get_signature_instance().all(ids=[signature.get_signature_id()])
        self.assertEqual(len(signatures), 1)
        self.assertEqual(signatures[0].get_signature_id(), signature.get_signature_id())

        # Test JSON serialization.
        signature_data = signature.to_dict()
        self.assertEqual(signature_data['signature_id'], signature_id)
        self.assertEqual(signature_data['signature_project_id'], project_id)
        self.assertEqual(signature_data['signature_reference_id'], user_id)
        self.assertEqual(signature_data['signature_reference_type'], 'user')
        self.assertEqual(signature_data['signature_type'], signature_type)
        self.assertEqual(signature_data['signature_signed'], signature_signed)
        self.assertEqual(signature_data['signature_approved'], signature_approved)

    def test_company_model(self):
        """Test company creation/querying/assertions."""
        # Create signature for tests.
        signature_id = str(uuid.uuid4())
        project_id = str(uuid.uuid4())
        user_id = str(uuid.uuid4())
        user_email = 'test@user.com'
        signature_type = 'cla'
        signature_corporate = True
        signature_signed = True
        signature_approved = True
        signature = get_signature_instance()
        signature.set_signature_id(signature_id)
        signature.set_signature_project_id(project_id)
        signature.set_signature_document_major_version(1)
        signature.set_signature_document_minor_version(0)
        signature.set_signature_reference_id(user_id)
        signature.set_signature_reference_type('user')
        signature.set_signature_type(signature_type)
        signature.set_signature_signed(signature_signed)
        signature.set_signature_approved(signature_approved)
        signature.save()
        # Create company for tests.
        company_id = str(uuid.uuid4())
        name = 'Org name'
        whitelist = ['whitelist.org', 'okdomain.com']
        whitelist_patterns = ['^info@.*', '.*admin.*']
        company = get_company_instance()
        company.set_company_id(company_id + '-FAKE')
        company.set_company_name(name + ' Bad')
        company.set_company_id(company_id)
        company.set_company_name(name)
        for wl_item in whitelist:
            company.add_company_whitelist(wl_item)
        for excp in whitelist_patterns:
            company.add_company_whitelist_pattern(excp)
        company.save()

        # Re-load company by ID.
        company2 = get_company_instance()
        company2.load(company_id)
        self.assertEqual(company.get_company_id(), company2.get_company_id())
        self.assertEqual(company.get_company_name(),
                         company2.get_company_name())
        self.assertEqual(company.get_company_whitelist(),
                         company2.get_company_whitelist())
        self.assertEqual(company.get_company_whitelist_patterns(),
                         company2.get_company_whitelist_patterns())

        # Test company not found query.
        company3 = get_company_instance()
        with self.assertRaises(DoesNotExist) as context:
            company3.load(str(uuid.uuid4()))
        self.assertTrue('Company not found' in str(context.exception))

        # Test add/remove whitelist.
        company.set_company_whitelist(['safe.org'])
        company.add_company_whitelist('another.org')
        self.assertTrue('safe.org' in company.get_company_whitelist())
        self.assertTrue('another.org' in company.get_company_whitelist())
        company.save()
        company4 = get_company_instance()
        company4.load(company.get_company_id())
        self.assertTrue('safe.org' in company4.get_company_whitelist())
        self.assertTrue('another.org' in company4.get_company_whitelist())
        company4.remove_company_whitelist('safe.org')
        company4.save()
        company4 = get_company_instance()
        company4.load(company.get_company_id())
        self.assertTrue('safe.org' not in company4.get_company_whitelist())
        self.assertTrue('another.org' in company4.get_company_whitelist())

        # Test add/remove whitelist patterns.
        company.set_company_whitelist_patterns(['^admin@.*$'])
        company.add_company_whitelist_pattern('.*@blacklist.org')
        self.assertTrue('^admin@.*$' in company.get_company_whitelist_patterns())
        self.assertTrue('.*@blacklist.org' in company.get_company_whitelist_patterns())
        company.save()
        company5 = get_company_instance()
        company5.load(company.get_company_id())
        self.assertTrue('^admin@.*$' in company5.get_company_whitelist_patterns())
        self.assertTrue('.*@blacklist.org' in company5.get_company_whitelist_patterns())
        company5.remove_company_whitelist_pattern('.*@blacklist.org')
        company5.save()
        company5 = get_company_instance()
        company5.load(company.get_company_id())
        self.assertTrue('^admin@.*$' in company5.get_company_whitelist_patterns())
        self.assertTrue('.*@blacklist.org' not in company5.get_company_whitelist_patterns())

        # Test JSON serialization.
        company_data = company.to_dict()
        self.assertEqual(company_data['company_id'], company_id)
        self.assertEqual(company_data['company_name'], name)
        self.assertEqual(company_data['company_whitelist'],
                         ['safe.org', 'another.org'])
        self.assertEqual(company_data['company_whitelist_patterns'],
                         ['^admin@.*$', '.*@blacklist.org'])

        # Test loading all companys.
        companys = get_company_instance().all()
        expected = [company.get_company_id() for company in companys]
        self.assertEqual(len(expected), 1)
        self.assertTrue(company.get_company_id() in expected)
        companys = get_company_instance().all(ids=[company.get_company_id()])
        self.assertEqual(len(companys), 1)
        self.assertEqual(companys[0].get_company_id(), company.get_company_id())
