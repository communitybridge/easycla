import gossip
import lfcore
import os


@gossip.register('lf.instance.compose.generation', tags=['ccc'])
def lf_instance_compose_generation(containers, config, dependencies, mode, path):

    cla = dependencies.get('cla')

    if cla.created:
        cla_endpoint = cla.endpoints.containers.get('workspace', 5000).formatted
        containers['workspace']['environment']['CLA_SERVER_URL'] = cla_endpoint
        lfcore.logger.info('Setting CLA_SERVER_URL to ' + containers['workspace']['environment']['CLA_SERVER_URL'])

        cinco = cla.instance.dependencies.get('cinco')

        if cinco.created:
            cinco_endpoint = cinco.endpoints.containers.get('workspace', 5000).formatted
            containers['workspace']['environment']['CINCO_SERVER_URL'] = cinco_endpoint
            lfcore.logger.info('Setting CINCO_SERVER_URL to ' + containers['workspace']['environment']['CINCO_SERVER_URL'])

            keycloak = cinco.instance.dependencies.get('keycloak')

            if keycloak.created:
                kc_endpoint = keycloak.endpoints.containers.get('workspace', 8080).formatted
                containers['workspace']['environment']['KEYCLOAK_SERVER_URL'] = kc_endpoint
                lfcore.logger.info('Setting KEYCLOAK_SERVER_URL to ' + containers['workspace']['environment']['KEYCLOAK_SERVER_URL'])
