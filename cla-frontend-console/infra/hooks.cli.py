import gossip
import lfcore
import os


@gossip.register('lf.instance.compose.generation', tags=['cla-console'])
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


@gossip.register('lf.instance.init.started', tags=['cla-console'])
def lf_instance_init_started(instance):
    # Adding CLA-Console endpoint to cla project and rebooting containers
    cla = instance.dependencies.get('cla')

    if cla.created:

        if not cla.instance.mode == 'dev':

            cla.instance.containers.destroy()

            cla_dc = cla.instance.containers.compose_definition
            endpoint = instance.endpoints.containers.get('nginx', 80, True).formatted
            cla_dc['services']['workspace']['environment']['CLA_CONSOLE_ENDPOINT'] = endpoint

            lfcore.utils.writeYaml(cla_dc, os.path.join(cla.instance.path, 'docker-compose.yml'))

            cla.instance.containers.up()
