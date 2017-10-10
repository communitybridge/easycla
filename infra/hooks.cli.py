import gossip
import lf
import os


@gossip.register('local.instance.init.docker-compose', tags=['cla-console'])
def local_init_docker_compose_file(containers, config, dependencies, envs, mode, path):

    cla = dependencies.get('cla')

    if cla:
        cla_endpoint = cla.endpoints.containers.get('workspace', 5000).formatted
        containers['workspace']['environment']['CLA_SERVER_URL'] = cla_endpoint
        lf.logger.info('Setting CLA_SERVER_URL to ' + containers['workspace']['environment']['CLA_SERVER_URL'])


@gossip.register('preprod_instance_task_build', tags=['cla-console'])
def preprod_instance_task_build(containers, instance_config, dependencies, domains, envs):
    if len(dependencies) >= 1:
