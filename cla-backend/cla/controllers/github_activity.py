# Copyright The Linux Foundation and each contributor to CommunityBridge.
# SPDX-License-Identifier: MIT

import falcon
import requests

import cla

GITHUB_ACTIVITY_ENDPOINT = "/github/activity"


def v4_easycla_github_activity(base_url: str, request: falcon.Request):
    """
    sends the post request coming from github to v4 api
    so we can start migrating some of the legacy code from
    Python -> Golang
    """
    fn = 'github_activity.v4_easycla_github_activity'
    if not base_url:
        raise ValueError("base url missing, can't find the easyCLA api")

    base_url = base_url.rstrip("/")
    if "v4" not in base_url:  # we need to add the prefix path for v4
        base_url = base_url + "/cla-service/v4"

    url = base_url + GITHUB_ACTIVITY_ENDPOINT
    headers = request.headers
    body = request.bounded_stream.read()

    cla.log.debug(f'{fn} - forwarding github activity to: {url}')
    resp = requests.post(url, data=body, headers=headers)
    cla.log.debug(f'{fn} - forwarding response status is: {resp.status_code}')

    # If the response was successful, no Exception will be raised
    resp.raise_for_status()
