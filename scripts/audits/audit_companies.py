import csv
import logging
import os
import sys

import boto3
import botocore.session
import click
from botocore.exceptions import ClientError

from audit import CompanyAudit


@click.command()
@click.option("--aws-profile", is_flag=False, default="default", help="The Named AWS profile")
@click.option(
    "--output-file",
    is_flag=False,
    help=" The output file for the companies audited records summarizing invalid records",
)
def main(aws_profile, output_file):
    """
    This script audits invalid records in the companies table and generates a report
    """
    try:
        if os.environ.get("STAGE") is None:
            logging.warning("Please set the 'STAGE' environment variable - typically one of: {dev, staging, prod}")
            return
        stage = os.environ.get("STAGE", "dev")
        companies_table_name = "cla-{}-companies".format(stage)
        users_table_name = "cla-{}-users".format(stage)
        session = boto3.Session(profile_name=aws_profile)
        dynamodb = session.resource("dynamodb")
        companies_table = dynamodb.Table(companies_table_name)
        users_table = dynamodb.Table(users_table_name)
        companies = companies_table.scan()["Items"]

        # set tables used in the audit process
        audit_company = CompanyAudit(dynamodb, batch=companies)
        audit_company.set_companies_table(companies_table)
        audit_company.set_users_table(users_table)
        invalid_fields = audit_company.process_batch()

        columns = ["company_id", "error_type", "column", "data"]
        with open(output_file, "w", newline="") as csv_file:
            writer = csv.DictWriter(csv_file, fieldnames=columns, delimiter=" ")
            writer.writeheader()
            writer.writerows(
                {
                    "company_id": audit["company_id"],
                    "error_type": audit["error_type"],
                    "column": audit["column"],
                    "data": audit["data"],
                }
                for audit in invalid_fields
            )

    except (Exception, ClientError) as err:
        logging.error(err)


if __name__ == "__main__":
    sys.exit(main())
