# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

USER_TABLE_DATA = {
    "Table": {
        "AttributeDefinitions": [
            {"AttributeName": "user_id", "AttributeType": "S"},
            {"AttributeName": "user_external_id", "AttributeType": "S"},
        ],
        "ItemCount": 0,
        "KeySchema": [{"AttributeName": "user_id", "KeyType": "HASH"}],
        "GlobalSecondaryIndexes": [
            {
                "IndexName": "github-user-external-id-index",
                "KeySchema": [{"AttributeName": "user_external_id", "KeyType": "HASH"}],
                "Projection": {"ProjectionType": "ALL"},
                "ProvisionedThroughput": {"ReadCapacityUnits": 1, "WriteCapacityUnits": 1},
            }
        ],
        "ProvisionedThroughput": {"ReadCapacityUnits": 1, "WriteCapacityUnits": 1},
    }
}

COMPANY_TABLE_DATA = {
    "Table": {
        "AttributeDefinitions": [{"AttributeName": "company_id", "AttributeType": "S"}],
        "KeySchema": [{"AttributeName": "company_id", "KeyType": "HASH"}],
    }
}

SIGNATURE_TABLE_DATA = {
    "Table": {
        "AttributeDefinitions": [
            {"AttributeName": "signature_id", "AttributeType": "S",},
            {"AttributeName": "signature_project_external_id", "AttrubuteType": "S"},
            {"AttributeName": "signature_company_project_external_id", "AttributeType": "S"},
            {"AttributeName": "signature_company_initial_manager_id", "AttributeType": "S"},
            {"AttributeName": "signature_company_secondary_manager_list", "AttributeType": "M"},
        ],
        "KeySchema": [{"AttributeName": "signature_id", "KeyType": "HASH"}],
        "GlobalSecondaryIndexes": [
            {
                "IndexName": "project-signature-external-id-index",
                "KeySchema": [{"AttributeName": "signature_project_external_id", "KeyType": "HASH"}],
                "Projection": {"ProjectionType": "ALL"},
                "ProvisionedThroughput": {"ReadCapacityUnits": 1, "WriteCapacityUnits": 1},
            },
            {
                "IndexName": "signature-company-signatory-index",
                "KeySchema": [{"AttributeName": "signature_company_project_external_id", "KeyType": "HASH"}],
                "Projection": {"ProjectionType": "ALL"},
                "ProvisionedThroughput": {"ReadCapacityUnits": 1, "WriteCapcityUnits": 1},
            },
            {
                "IndexName": "signature-company-initial-manager-index",
                "KeySchema": [{"AttributeName": "signature_company_initial_manager_id", "KeyType": "HASH"}],
                "Projection": {"ProjectionType": "ALL"},
                "ProvisionedThroughput": {"ReadCapacityUnits": 1, "WriteCapacityUnits": 1},
            },
        ],
    }
}
