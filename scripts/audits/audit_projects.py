# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import csv
import logging
import os
import sys

import boto3
import botocore.session
import click
from botocore.exceptions import ClientError

from audit import ProjectAudit


@click.command()
@click.option("--aws-profile", is_flag=False, default="default", help="The Named AWS profile")
@click.option(
    "--output-file", is_flag=False, help="The output file showing audit report for invalid records",
)
def main(aws_profile, output_file):
    """
    This script audits invalid records in the projects table - specifically projects template pdfs
    """
    try:
        if os.environ.get("STAGE") is None:
            logging.warning("Please set the 'STAGE' environment variable - typically one of: {dev, staging, prod}")
            return
        stage = os.environ.get("STAGE", "dev")
        projects_table_name = "cla-{}-projects".format(stage)
        session = boto3.Session(profile_name=aws_profile)
        dynamodb = session.resource("dynamodb")
        projects_table = dynamodb.Table(projects_table_name)
        projects = projects_table.scan()["Items"]

        # set the projects table used in the audit process
        audit_project = ProjectAudit(dynamodb, batch=projects)
        invalid_fields = audit_project.process_batch()

        columns = ["project_id", "error_type", "column", "data"]
        with open(output_file, "w", newline="") as csv_file:
            writer = csv.DictWriter(csv_file, fieldnames=columns, delimiter=" ")
            writer.writeheader()
            writer.writerows(
                {
                    "project_id": audit["project_id"],
                    "error_type": audit["error_type"],
                    "column": audit["column"],
                    "data": audit["data"],
                }
                for audit in invalid_fields
            )

    except (Exception, ClientError) as err:
        logging.error(err,exc_info=True)


if __name__ == "__main__":
    sys.exit(main())
