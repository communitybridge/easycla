# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

from datetime import datetime, timedelta
import json
import os
import re

# To get updated botocore data files
os.environ['AWS_DATA_PATH'] = '.'

import boto3
import botocore
from botocore.exceptions import ClientError
from botocore.vendored import requests

CLIENT = boto3.client('dynamodb')
SLACK_WEBHOOK = os.environ.get('SLACK_WEBHOOK')
REGION = os.environ.get('AWS_DEFAULT_REGION')
CONSOLE_ENDPOINT = 'https://console.aws.amazon.com/dynamodb/home?region={region}#backups:'.format(region=REGION)


def main(event, context):
    tables = get_tables_to_backup()
    results = {
        "success": [],
        "failure": []
    }

    for table in tables:
        try:
            create_backup(table)
            results['success'].append(table)
        except Exception as e:
            print("Error creating backup for table {table}.\n. Error: {err}".format(table=table, err=str(e)))
            results['failure'].append(table)

    if os.environ.get('BACKUP_REMOVAL_ENABLED') == 'true':
        try:
            remove_stale_backups(tables)
        except Exception as e:
            print("Error removing stale backups. Error: {err}".format(err=str(e)))

    message = format_message(results)
    send_to_slack(message)


def create_backup(table):
    timestamp = datetime.now().strftime('%Y%m%d%H%M%S')
    backup_name = table + "_" + timestamp
    CLIENT.create_backup(
        TableName=table,
        BackupName=backup_name
    )

def remove_stale_backups(tables):
    paginator = CLIENT.get_paginator('list_backups')
    upper_bound = datetime.now() - timedelta(days=int(os.environ.get('BACKUP_RETENTION_DAYS')))

    print("Removing backups before the following date: {date}".format(date=upper_bound))

    for page in paginator.paginate(TimeRangeUpperBound=upper_bound):
        for table in page['BackupSummaries']:
            if table['TableName'] in tables:
                CLIENT.delete_backup(BackupArn=table['BackupArn'])


def format_message(results):
    msg = ""

    if not results['success'] and not results['failure']:
        return "Tried running DynamoDB backup, but no tables were specified.\nPlease check your configuration."

    msg += "Tried to backup {total} DynamoDB tables. {successes} succeeded, and {failures} failed. See all backups <{url}|here>.".format(
        total=(len(results['success']) + len(results['failure'])),
        successes=len(results['success']),
        failures=len(results['failure']),
        url=CONSOLE_ENDPOINT.format(region=REGION)
    )

    if results['success']:
        msg += "\nThe following tables were successful:\n - "
        msg += "\n - ".join(results['success'])

    if results['failure']:
        msg += "\nThe following tables failed:\n - "
        msg += "\n - ".join(results['failure'])

    return msg


def send_to_slack(message):
    if not SLACK_WEBHOOK:
        print('No SLACK_WEBHOOK provided. Not sending a message...')
        return
    data = {"text": message}
    resp = requests.post(SLACK_WEBHOOK, json=data)

    resp.raise_for_status()


def get_tables_to_backup():
    """Determines which tables to backup. The determination is made based on
    the config options. Return value is a list.

    Order is as follows:

    1. If the TABLE_REGEX environment variable is set, call the `ListTables` API
    for DynamoDB and return tables that match the given TABLE_REGEX;

    2. If the TABLE_FILE environment variable is set, load the TABLE_FILE and
    return the tables list.

    3. If the TABLE_NAME environment variable is set, return the TABLE_NAME.

    It will not combine multiple options. It will return the value(s) from the first
    option with the environment variable present.
    """
    if os.environ.get('TABLE_REGEX'):
        return get_tables_regex(os.environ.get('TABLE_REGEX'))
    elif os.environ.get('TABLE_FILE'):
        return get_tables_from_file(os.environ.get('TABLE_FILE'))
    elif os.environ.get('TABLE_NAME'):
        return [os.environ.get('TABLE_NAME')]

    print("No tables configured. Please use TABLE_REGEX, TABLE_FILE, OR TABLE_NAME environment variables.")

    return []


def get_tables_regex(pattern):
    print("Using regex pattern {} to find tables.".format(pattern))
    tables = []
    paginator = CLIENT.get_paginator('list_tables')
    for page in paginator.paginate():
        for table in page['TableNames']:
            if re.match(pattern, table):
                tables.append(table)

    return tables


def get_tables_from_file(filename):
    print("Using local file {} to find tables.".format(filename))
    with open(filename, 'r') as f:
        tables = json.load(f)

    return tables


if __name__ == "__main__":
    main('', '')
