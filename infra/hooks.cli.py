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
