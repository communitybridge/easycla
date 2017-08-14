"""
Tests to be run for all supported key-value stores.
"""

import unittest

# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
from cla.models import DoesNotExist
from cla.models.dynamo_models import Store as DynamoDBStore

class KeyValueStoreTestCase(CLATestCase):
    """Key-value store test cases."""
    def test_dynamodb_store(self):
        """Test using dynamodb as key-value store."""
        store = DynamoDBStore()
        self.assertEqual(store.delete('test'), None)
        with self.assertRaises(DoesNotExist) as context:
            store.get('test')
        self.assertTrue('Key not found' in str(context.exception))
        self.assertFalse(store.exists('test'))
        store.set('test', 'value')
        self.assertTrue(store.exists('test'))
        self.assertEqual(store.get('test'), 'value')
        self.assertEqual(store.delete('test'), None)

if __name__ == '__main__':
    unittest.main()
