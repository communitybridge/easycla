# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import json
import os

import requests
from boto3 import client

import cla
from cla import log
from cla.config import THE_LINUX_FOUNDATION

STAGE = os.environ.get('STAGE', '')
REGION = 'us-east-1'


class ProjectService:
    """
    ProjectService Handles external salesforce Project
    """

    def __init__(self):
        self.platform_gateway_url = self.get_ssm_key(REGION, f'cla-auth0-platform-api-gw-{STAGE}')

    def get_ssm_key(self, region, key):
        """
        Fetches the specified SSM key value from the SSM key store
        :param region: aws region
        :type region: string
        :parm key: key
        :type key: string
        """
        ssm_client = client('ssm', region_name=region)

        log.debug(f'Fetching Key: {key}')
        response = ssm_client.get_parameter(Name=key, WithDecryption=True)
        log.debug(f'Fetched Key: {key}, value: {response["Parameter"]["Value"]}')
        return response['Parameter']['Value']

    def is_standalone(self, project_sfid) -> bool:
        """
        Checks if salesforce project is a stand alone project (No subprojects and parent)
        :param project_sfid: salesforce project_id
        :type project_sfid: string
        :return: Check whether sf project is a stand alone
        :rtype: Boolean
        """
        project = self.get_project_by_id(project_sfid)
        if project:
            parent_sf_id = project.get('Parent', None)
            if not self.has_parent(project) and (parent_sf_id is None or parent_sf_id == THE_LINUX_FOUNDATION):
                return True
        return False

    def is_lf_supported(self, project_sfid) -> bool:
        """
        Checks if salesforce project is a LF Supported project
        :param project_sfid: salesforce project_id
        :type project_sfid: string
        :return: Check whether sf project is a stand alone
        :rtype: Boolean
        """
        project = self.get_project_by_id(project_sfid)
        if project:
            return (project.get('Funding', None) == 'Unfunded' or
                    project.get('Funding', None) == 'Supported By Parent') and \
                   project.get('Parent', None) == THE_LINUX_FOUNDATION
        return False

    def has_parent(self, project) -> bool:
        """ checks if project has parent """
        try:
            log.info(f"Checking if {project['Name']} has parent project")
            parent = project['Parent']
            if parent:
                return True
        except KeyError as err:
            log.debug(f"Failed to find parent for {project['Name']} , error: {err}")
            return False
        return False

    def is_parent(self, project) -> bool:
        """ 
        checks whether salesforce project is a parent
        :param project: salesforce project
        :type project: dict
        :return: Whether salesforce project is a parent
        :rtype: Boolean
        """
        try:
            log.info(f"Checking if {project['Name']} is a parent")
            project_type = project['ProjectType']
            if project_type == 'Project Group':
                return True
        except KeyError as err:
            log.debug(f"Failed to get ProjectType for project: {project['Name']}  error: {err}")
            return False
        return False

    def get_project_by_id(self, project_id):
        """
        Gets Salesforce project by ID
        """
        headers = {
            'Authorization': f'bearer {self.get_access_token()}',
            'accept': 'application/json'
        }
        try:
            url = f'{self.platform_gateway_url}/project-service/v1/projects/{project_id}'
            cla.log.debug(f'Sending GET request to {url}')
            r = requests.get(url, headers=headers)
            r.raise_for_status()
            response_model = json.loads(r.text)
            return response_model
        except requests.exceptions.HTTPError as err:
            msg = f'Could not get project: {project_id}, error: {err}'
            cla.log.warning(msg)
            return None

    def get_access_token(self):
        auth0_url = self.get_ssm_key(REGION, f'cla-auth0-platform-url-{STAGE}')
        platform_client_id = self.get_ssm_key(REGION, f'cla-auth0-platform-client-id-{STAGE}')
        platform_client_secret = self.get_ssm_key(REGION, f'cla-auth0-platform-client-secret-{STAGE}')
        platform_audience = self.get_ssm_key(REGION, f'cla-auth0-platform-audience-{STAGE}')

        auth0_payload = {
            'grant_type': 'client_credentials',
            'client_id': platform_client_id,
            'client_secret': platform_client_secret,
            'audience': platform_audience
        }

        headers = {
            'content-type': 'application/x-www-form-urlencoded',
            'accept': 'application/json'
        }

        access_token = ''
        try:
            # logger.debug(f'Sending POST to {auth0_url} with payload: {auth0_payload}')
            log.debug(f'Sending POST to {auth0_url}')
            r = requests.post(auth0_url, data=auth0_payload, headers=headers)
            r.raise_for_status()
            json_data = json.loads(r.text)
            access_token = json_data["access_token"]
            return access_token
        except requests.exceptions.HTTPError as err:
            log.warning(f'Could not get auth token, error: {err}')
            return None
