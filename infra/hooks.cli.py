import gossip
import lf
import os
import requests
import json


def public_ip():
    r = requests.get('http://httpbin.org/ip')
    d = json.loads(r.content.decode('utf-8'))
    return d['origin']


def host_port(containers, container_name, idx=0):

    for key, port in enumerate(containers[container_name]['ports']):
        p = port.split(':')
        c_p = containers[container_name]['ports'][idx].split(':')[1]
        if p[1] == c_p:
            return str(p[0])


@gossip.register('local.instance.init.docker-compose', tags=['cla'])
def local_init_docker_compose_file(containers, config, dependencies, envs, mode, path):
    """
    All necessary steps to integrate CINCO and Keycloak into the CLA system.

    If dev mode, we don't provision a workspace container and assume the user has manually
    configured the appropriate dependencies.
    """
    if mode == 'dev':
        lf.logger.info('No workspace container provisioned in dev mode')
        lf.logger.info('Ensure the cla_config.py file has the appropriate configuration for your \
                       external dependencies')
        del containers['workspace']
    else:
        # Should we use CINCO's DynamoDB instance here?
        lf.logger.info('Automatically setting CLA dependency endpoints')
        platform = dependencies.get('cinco')
        if platform:
            cinco_endpoint = platform.endpoints.containers.get('workspace', 5000).formatted
            lf.logger.info('CINCO: %s', cinco_endpoint)
            containers['workspace']['environment']['CLA_CINCO_ENDPOINT'] = cinco_endpoint
            keycloak = platform.instance.dependencies.get('keycloak')
            if keycloak:
                kc_endpoint = keycloak.endpoints.containers.get('workspace', 8080).formatted
                lf.logger.info('Keycloak: %s', kc_endpoint)
                containers['workspace']['environment']['CLA_KEYCLOAK_ENDPOINT'] = kc_endpoint
        # Set the callback_url to the CLA instance.
        ip = public_ip()
        port = host_port(containers, 'workspace')
        base_url = 'http://' + ip
        callback_url = base_url + '/v1/signed'
        github_oauth_url = base_url + '/v1/github/installation'
        containers['workspace']['environment']['CLA_BASE_URL'] = base_url
        containers['workspace']['environment']['CLA_SIGNED_CALLBACK_URL'] = callback_url
        containers['workspace']['environment']['CLA_GITHUB_OAUTH_CALLBACK_URL'] = github_oauth_url

        lf.logger.warning('You public IP address (%s) was used as the base_url for the CLA system.' %ip)
        lf.logger.warning('If you are behind a NAT/firewall, you will need to add port forwarding from the edge of your network to you local machine (80 -> %s)' %port)
        lf.logger.warning('The CLA system will work without the port forwading setup, but you will not be able to test GitHub integration (starting flow from GitHub) and DocuSign callbacks (confirmation of completed signatures)')
        lf.logger.warning('Tip: You can still access the CLA system locally via http://localhost:%s' %port)
