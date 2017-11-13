import gossip
import lfcore


@gossip.register('lf.instance.compose.generation', tags=['pmc'])
def local_init_docker_compose_file(containers, config, dependencies, mode, path):
    platform = dependencies.get('cinco')

    if platform.created:
        cinco_endpoint = platform.endpoints.containers.get('workspace', 5000).formatted
        containers['workspace']['environment']['CINCO_SERVER_URL'] = cinco_endpoint
        lfcore.logger.info('Setting CINCO_SERVER_URL to ' + containers['workspace']['environment']['CINCO_SERVER_URL'])

        keycloak = platform.instance.dependencies.get('keycloak')

        if keycloak.created:
            kc_endpoint = keycloak.endpoints.containers.get('workspace', 8080).formatted
            containers['workspace']['environment']['KEYCLOAK_SERVER_URL'] = kc_endpoint
            lfcore.logger.info('Setting KEYCLOAK_SERVER_URL to ' + containers['workspace']['environment']['KEYCLOAK_SERVER_URL'])


@gossip.register('preprod_instance_task_build', tags=['pmc'])
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

        lfcore.logger.info('Setting CINCO_SERVER_URL to ' + 'https://' + domains['primary'])
        lfcore.logger.info('Setting KEYCLOAK_SERVER_URL to ' + kc_endpoint)
