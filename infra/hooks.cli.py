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

    # If dev mode, we don't provision a workspace container
    if mode == 'dev':
        del containers['workspace']

        # We need to write the configuration file
