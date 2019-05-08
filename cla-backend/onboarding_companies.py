import json
import os
import re
import uuid

from datetime import datetime, timedelta

from pynamodb.models import Model
from cla.models.dynamo_models import User, Signature, Project, Company

# boto3 requires AWS credentials. Please set either AWS_PROFILE or
# (AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY) environment variables
# prior to running.

dry_run = False
company_input_file = 'tungsten_fabric.json'
stage = os.environ['STAGE']

def main(event, context):
    if dry_run:
        print('Performing dry run')
    else:
        exit_cmd = input('This is NOT a dry run. You are running the script for the {} environment. Press enter to continue ("exit" to exit): '.format(stage))
        if exit_cmd == 'exit':
            return

    with open(company_input_file) as json_file:
        input_file = json.load(json_file)
        project_id = input_file['project']['project_id']
        if project_id == "":
            print('project_id is required')
            return

        for company in input_file['companies']:
            # Generate company ID
            company_id = str(uuid.uuid4())
            company['company_id'] = company_id
            print('\n\nConfiguring company: {} ({})'.format(company, company_id))

            # Create user for each CLA Manager specified
            cla_manager_id = None
            for user in company['company_acl']:
                user_id = str(uuid.uuid4())
                print('Creating user: {} ({})'.format(user, user_id))
                if not dry_run:
                    try:
                        create_user(user, user_id)
                        cla_manager_id = user_id
                    except Exception as e:
                        print("Error creating user data. \nError: {err}".format(err=str(e)))

            if cla_manager_id == None:
                print('No cla manager. Skipping company: ({})'.format(company['company_name']))

            print('Creating company: {}'.format(company))
            # Create Company
            cla_manager_LFIDs = [user['username'] for user in company['company_acl']]
            if not dry_run:
                # Get the LFIDs of the company managers
                try:
                    create_company(company, cla_manager_LFIDs, cla_manager_id)
                except Exception as e:
                    print("Error creating company data. \nError: {err}".format(err=str(e)))
                    # Iterate to the next company
                    continue

            # Create CCLA signature for the company for the project
            signature_id = str(uuid.uuid4())
            print('Creating signature ({}) for company: {} project: {}'.format(signature_id, company, project_id))
            if not dry_run:
                try:
                    create_company_signature(signature_id, company['company_id'], project_id, cla_manager_LFIDs)
                except Exception as e:
                    print("Error creating company signature data. \nError: {err}".format(err=str(e)))
                    continue

def create_user(user, user_id):
    user_model = User()
    user_model.set_user_id(user_id)
    user_model.set_user_name(user['name'])
    user_model.set_lf_email(user['email'])
    user_model.set_lf_username(user['username'])
    user_model.save()

def create_company(company, cla_managers, cla_manager_id):
    company_model = Company(
        company_id=company['company_id'],
        company_name=company['company_name'],
        company_manager_id=cla_manager_id,
        company_acl=set(cla_managers)
    )

    company_model.save()

def create_company_signature(signature_id, company_id, project_id, cla_managers):
    # Use Constructor to include Signature ACL
    signature_model = Signature(
        signature_id=signature_id,
        signature_project_id=project_id,
        signature_reference_type='company',
        signature_reference_id=company_id,
        signature_document_major_version=1,
        signature_document_minor_version=1,
        signature_type='ccla',
        signature_signed=True,
        signature_approved=True,
        signature_acl = set(cla_managers)
    )

    signature_model.save()

if __name__ == "__main__":
    main('', '')
