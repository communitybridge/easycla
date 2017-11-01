import gossip
import lf
import os


@gossip.register('local.instance.init.docker-compose', tags=['ccc'])
def local_init_docker_compose_file(containers, config, dependencies, envs, mode, path):
    platform = dependencies.get('cinco')

    if platform:
        cinco_endpoint = platform.endpoints.containers.get('workspace', 5000).formatted
        containers['workspace']['environment']['CINCO_SERVER_URL'] = cinco_endpoint
        lf.logger.info('Setting CINCO_SERVER_URL to ' + containers['workspace']['environment']['CINCO_SERVER_URL'])

        keycloak = platform.instance.dependencies.get('keycloak')

        if keycloak:
            kc_endpoint = keycloak.endpoints.containers.get('workspace', 8080).formatted
            containers['workspace']['environment']['KEYCLOAK_SERVER_URL'] = kc_endpoint
            lf.logger.info('Setting KEYCLOAK_SERVER_URL to ' + containers['workspace']['environment']['KEYCLOAK_SERVER_URL'])

    cla = dependencies.get('cla')

    if cla:
        cla_endpoint = cla.endpoints.containers.get('workspace', 5000).formatted
        containers['workspace']['environment']['CLA_SERVER_URL'] = cla_endpoint
        lf.logger.info('Setting CLA_SERVER_URL to ' + containers['workspace']['environment']['CLA_SERVER_URL'])


@gossip.register('preprod_instance_task_build', tags=['ccc'])
def preprod_instance_task_build(containers, instance_config, dependencies, domains, envs):
    if len(dependencies) >= 1:
        platform = dependencies.get('cinco')

        task = platform.artifacts.get('ECSPreprodTask')
        workspace = [x for x in task.containers if x['name'] == 'workspace'][0]
        kc_endpoint = [x['value'] for x in workspace['environment'] if x['name'] == 'KEYCLOAK_SERVER_URL'][0]

        envs.append({
            'name': 'CINCO_SERVER_URL',
            'value': 'https://' + platform.domain
        })
        envs.append({
            'name': 'KEYCLOAK_SERVER_URL',
            'value': kc_endpoint
        })

        lf.logger.info('Setting CINCO_SERVER_URL to ' + 'https://' + domains['primary'])
        lf.logger.info('Setting KEYCLOAK_SERVER_URL to ' + kc_endpoint)

        cla = dependencies.get('cla')

        envs.append({
            'name': 'CLA_SERVER_URL',
            'value': 'https://' + cla.domain
        })

        lf.logger.info('Setting CLA_SERVER_URL to ' + 'https://' + cla.domain)
