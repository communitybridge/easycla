"""
This test assumes that the local instance of DynamoDB is currently running at
http://localhost:8000 (the default).

For local development with DynamoDB, you may want to check out this link on
how to setup the service locally:
http://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.html

TLDR for Ubuntu:

    sudo apt-get install default-jre (if you don't already have java)
    wget https://s3-us-west-2.amazonaws.com/dynamodb-local/dynamodb_local_latest.tar.gz
    tar xf dynamodb_local_latest.tar.gz
    java -Djava.library.path=./DynamoDBLocal_lib -jar DynamoDBLocal.jar -sharedDb
"""

import unittest
# Importing to setup proper python path and DB for tests.
from test_cla import CLATestCase
import cla
# Needs to be imported AFTER cla as the DB configuration needs to be in place
# before the models are loaded.
from test_databases import DatabaseTestCase # pylint: disable=wrong-import-position
from test_key_value_stores import KeyValueStoreTestCase # pylint: disable=wrong-import-position

class DynamoDBTestCase(DatabaseTestCase):
    """Database tests for the DynamoDB storage engine."""
    pass

class DynamoDBStoreTestCase(KeyValueStoreTestCase):
    """Key-value tests for the DynamoDB storage engine."""
    pass

if __name__ == '__main__':
    unittest.main()
