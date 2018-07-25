"""
Tests having to do with the configuration object.
"""

import os
import unittest
# Importing to setup proper python path and DB for tests.
import test_cla
import cla
from cla.utils import get_database_models, get_signing_service

class ConfigTestCase(unittest.TestCase):
    """Tests for CLA configuration files."""
    def test_config_import(self): # pylint: disable=no-self-use
        """Test instance configuration import."""
        # Ensure this doesn't give an import error.
        cla.Config('fake_config')

    def test_db_import(self):
        """Test invalid DB selected."""
        conf = {'DATABASE': 'FakeDB'}
        with self.assertRaises(Exception) as context:
            get_database_models(conf)
        self.assertTrue('Invalid database selection in configuration: FakeDB' in str(context.exception))
        # Ensure no issues with valid database and document signing service.
        get_database_models(cla.conf)

    def test_signing_service_import(self):
        """Test invalid signing service selected."""
        conf = {'SIGNING_SERVICE': 'FakeSigningService'}
        with self.assertRaises(Exception) as context:
            get_signing_service(conf)
        self.assertTrue('Invalid signing service selected in configuration: FakeSigningService' in str(context.exception))

    def test_environment_variable(self):
        """Test environment values for configuration options."""
        os.environ['CLA_TEST'] = 'testing'
        os.environ['CLA_DATABASE'] = 'TestDB'
        test_config = cla.Config()
        self.assertEqual(test_config['TEST'], 'testing')
        self.assertEqual(test_config['DATABASE'], 'TestDB')

if __name__ == '__main__':
    unittest.main()
