import json
import os
import re
import uuid

from datetime import datetime, timedelta
from pynamodb.models import Model
from cla.models.dynamo_models import User, Signature, Project, Company


def main(event, context):
    with open('sample.json') as json_file:  
        input = json.load(json_file)
        for company in input['companies']: 
            # Create company ID
            company_id = str(uuid.uuid4()) 
            company['company_id'] = company_id
            project = input['project']
            
            # Create users in database from company ACL
            for user in company['company_acl']: 
                try:
                    create_user(user)
                except Exception as e:
                    print("Error creating user data. \nError: {err}".format(err=str(e)))
            
            # Get only the LFIDs of the company managers
            cla_managers = [user['username'] for user in company['company_acl']]
            print(cla_managers)
            try:
                create_company(company, cla_managers)
            except Exception as e:
                print("Error creating company data. \nError: {err}".format(err=str(e)))
                # Iterate to the next company
                continue
            try:
                create_company_signature(company, project, cla_managers)
            except Exception as e: 
                print("Error creating company signature data. \nError: {err}".format(err=str(e)))

    
def create_company(company, cla_managers):
    # Use Constructor to include Company ACL
    company_model = Company(
        company_acl=set(cla_managers)
    )
    company_model.set_company_id(company['company_id'])
    company_model.set_company_name(company['company_name'] )
    company_model.save()

def create_company_signature(company, project, cla_managers):
    # Use Constructor to include Signature ACL
    signature_model = Signature(
        signature_acl = set(cla_managers)
    )
    signature_model.set_signature_id(str(uuid.uuid4()))
    signature_model.set_signature_project_id(project['project_id'])
    signature_model.set_signature_reference_id(company['company_id'])
    signature_model.set_signature_reference_type('company')
    signature_model.set_signature_document_major_version("1")
    signature_model.set_signature_document_minor_version("1")
    signature_model.set_signature_type('ccla')
    signature_model.set_signature_signed(True)
    signature_model.set_signature_approved(True)
    signature_model.save()

def create_user(user):
    user_model = User()
    user_model.set_user_id(str(uuid.uuid4()))
    user_model.set_user_name(user['name'])
    user_model.set_lf_email(user['email'])
    user_model.set_lf_username(user['username'])
    user_model.save()


if __name__ == "__main__":
    main('', '')
