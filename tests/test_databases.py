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
                      get_agreement_instance, get_repository_instance, get_organization_instance, \
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
        repo_agreement_type = 'url+pdf'
        repo_agreement_content = 'https://some-github-url.com/cla.pdf'
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

    def test_agreement_model(self):
        """Test agreement creation/querying/assertions."""
        agreement_id = str(uuid.uuid4())
        project_id = self.create_project()['project_id']
        user_id = str(uuid.uuid4())
        agreement_type = 'cla'
        agreement_corporate = True
        agreement_signed = True
        agreement_approved = True
        agreement = get_agreement_instance()
        agreement.set_agreement_id(agreement_id)
        agreement.set_agreement_project_id(project_id)
        agreement.set_agreement_document_revision(1)
        agreement.set_agreement_reference_id(user_id)
        agreement.set_agreement_reference_type('user')
        agreement.set_agreement_type(agreement_type)
        agreement.set_agreement_signed(agreement_signed)
        agreement.set_agreement_approved(agreement_approved)
        agreement.save()

        # Re-load agreement by ID.
        agreement2 = get_agreement_instance()
        agreement2.load(agreement_id)
        self.assertEqual(agreement.get_agreement_id(), agreement2.get_agreement_id())
        self.assertEqual(agreement.get_agreement_project_id(), agreement2.get_agreement_project_id())
        self.assertEqual(agreement.get_agreement_reference_id(),
                         agreement2.get_agreement_reference_id())
        self.assertEqual(agreement.get_agreement_reference_type(),
                         agreement2.get_agreement_reference_type())
        self.assertEqual(agreement.get_agreement_type(), agreement2.get_agreement_type())
        self.assertEqual(agreement.get_agreement_signed(), agreement2.get_agreement_signed())
        self.assertEqual(agreement.get_agreement_approved(), agreement2.get_agreement_approved())

        # Test agreement not found query.
        agreement3 = get_agreement_instance()
        with self.assertRaises(DoesNotExist) as context:
            agreement3.load(str(uuid.uuid4()))
        self.assertTrue('Agreement not found' in str(context.exception))

        # Test loading all agreements.
        agreements = get_agreement_instance().all()
        expected = [agreement.get_agreement_id() for agreement in agreements]
        self.assertEqual(len(expected), 1)
        self.assertTrue(agreement.get_agreement_id() in expected)
        agreements = get_agreement_instance().all(ids=[agreement.get_agreement_id()])
        self.assertEqual(len(agreements), 1)
        self.assertEqual(agreements[0].get_agreement_id(), agreement.get_agreement_id())

        # Test JSON serialization.
        agreement_data = agreement.to_dict()
        self.assertEqual(agreement_data['agreement_id'], agreement_id)
        self.assertEqual(agreement_data['agreement_project_id'], project_id)
        self.assertEqual(agreement_data['agreement_reference_id'], user_id)
        self.assertEqual(agreement_data['agreement_reference_type'], 'user')
        self.assertEqual(agreement_data['agreement_type'], agreement_type)
        self.assertEqual(agreement_data['agreement_signed'], agreement_signed)
        self.assertEqual(agreement_data['agreement_approved'], agreement_approved)

    def test_organization_model(self):
        """Test organization creation/querying/assertions."""
        # Create agreement for tests.
        agreement_id = str(uuid.uuid4())
        project_id = str(uuid.uuid4())
        user_id = str(uuid.uuid4())
        user_email = 'test@user.com'
        agreement_type = 'cla'
        agreement_corporate = True
        agreement_signed = True
        agreement_approved = True
        agreement = get_agreement_instance()
        agreement.set_agreement_id(agreement_id)
        agreement.set_agreement_project_id(project_id)
        agreement.set_agreement_document_revision(1)
        agreement.set_agreement_reference_id(user_id)
        agreement.set_agreement_reference_type('user')
        agreement.set_agreement_type(agreement_type)
        agreement.set_agreement_signed(agreement_signed)
        agreement.set_agreement_approved(agreement_approved)
        agreement.save()
        # Create organization for tests.
        organization_id = str(uuid.uuid4())
        name = 'Org name'
        whitelist = ['whitelist.org', 'okdomain.com']
        exclude_patterns = ['^info@.*', '.*admin.*']
        organization = get_organization_instance()
        organization.set_organization_id(organization_id + '-FAKE')
        organization.set_organization_name(name + ' Bad')
        organization.set_organization_id(organization_id)
        organization.set_organization_name(name)
        for wl_item in whitelist:
            organization.add_organization_whitelist(wl_item)
        for excp in exclude_patterns:
            organization.add_organization_exclude_pattern(excp)
        organization.save()

        # Re-load organization by ID.
        organization2 = get_organization_instance()
        organization2.load(organization_id)
        self.assertEqual(organization.get_organization_id(), organization2.get_organization_id())
        self.assertEqual(organization.get_organization_name(),
                         organization2.get_organization_name())
        self.assertEqual(organization.get_organization_whitelist(),
                         organization2.get_organization_whitelist())
        self.assertEqual(organization.get_organization_exclude_patterns(),
                         organization2.get_organization_exclude_patterns())

        # Test organization not found query.
        organization3 = get_organization_instance()
        with self.assertRaises(DoesNotExist) as context:
            organization3.load(str(uuid.uuid4()))
        self.assertTrue('Organization not found' in str(context.exception))

        # Test add/remove whitelist.
        organization.set_organization_whitelist(['safe.org'])
        organization.add_organization_whitelist('another.org')
        self.assertTrue('safe.org' in organization.get_organization_whitelist())
        self.assertTrue('another.org' in organization.get_organization_whitelist())
        organization.save()
        organization4 = get_organization_instance()
        organization4.load(organization.get_organization_id())
        self.assertTrue('safe.org' in organization4.get_organization_whitelist())
        self.assertTrue('another.org' in organization4.get_organization_whitelist())
        organization4.remove_organization_whitelist('safe.org')
        organization4.save()
        organization4 = get_organization_instance()
        organization4.load(organization.get_organization_id())
        self.assertTrue('safe.org' not in organization4.get_organization_whitelist())
        self.assertTrue('another.org' in organization4.get_organization_whitelist())

        # Test add/remove exclude patterns.
        organization.set_organization_exclude_patterns(['^admin@.*$'])
        organization.add_organization_exclude_pattern('.*@blacklist.org')
        self.assertTrue('^admin@.*$' in organization.get_organization_exclude_patterns())
        self.assertTrue('.*@blacklist.org' in organization.get_organization_exclude_patterns())
        organization.save()
        organization5 = get_organization_instance()
        organization5.load(organization.get_organization_id())
        self.assertTrue('^admin@.*$' in organization5.get_organization_exclude_patterns())
        self.assertTrue('.*@blacklist.org' in organization5.get_organization_exclude_patterns())
        organization5.remove_organization_exclude_pattern('.*@blacklist.org')
        organization5.save()
        organization5 = get_organization_instance()
        organization5.load(organization.get_organization_id())
        self.assertTrue('^admin@.*$' in organization5.get_organization_exclude_patterns())
        self.assertTrue('.*@blacklist.org' not in organization5.get_organization_exclude_patterns())

        # Test JSON serialization.
        organization_data = organization.to_dict()
        self.assertEqual(organization_data['organization_id'], organization_id)
        self.assertEqual(organization_data['organization_name'], name)
        self.assertEqual(organization_data['organization_whitelist'],
                         ['safe.org', 'another.org'])
        self.assertEqual(organization_data['organization_exclude_patterns'],
                         ['^admin@.*$', '.*@blacklist.org'])

        # Test loading all organizations.
        organizations = get_organization_instance().all()
        expected = [organization.get_organization_id() for organization in organizations]
        self.assertEqual(len(expected), 1)
        self.assertTrue(organization.get_organization_id() in expected)
        organizations = get_organization_instance().all(ids=[organization.get_organization_id()])
        self.assertEqual(len(organizations), 1)
        self.assertEqual(organizations[0].get_organization_id(), organization.get_organization_id())
