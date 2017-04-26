import gossip
import core


@gossip.register('local_instance_build_dockerfile', tags=['pmc'])
def local_instance_build_dockerfile(containers, instance_config, dependencies, envs):
    if len(dependencies) >= 1:
        platform_instance = dependencies[0]
        docker_config = platform_instance.artifacts.get('DockerComposeFile')

        for key, port in enumerate(docker_config.config['services']['workspace']['ports']):
            p = port.split(':')
            if p[1] == '5000':
                envs['CINCO_SERVER_URL'] = 'http://' + platform_instance.containers.bridge_ip + ':' + p[0] + '/'
                core.logger.info('Setting CINCO_SERVER_URL to ' + envs['CINCO_SERVER_URL'])

    port = containers['workspace']['ports'][0].split(':')
    envs['CINCO_CONSOLE_URL'] = 'http://' + core.config.local['hostname'] + ':' + port[0] + '/'
    core.logger.info('Setting CINCO_CONSOLE_URL to ' + envs['CINCO_CONSOLE_URL'])


@gossip.register('procedures.local.instance_create.execute', tags=['pmc'])
def local_instance_create_execute(procedure):
    procedure.extra['autorun'] = core.utils. \
        yes_or_no('Would you like to automatically run the console on container startup?')


@gossip.register('preprod_instance_task_build', tags=['pmc'])
def preprod_instance_task_build(containers, instance_config, dependencies, domains, envs):
    if len(dependencies) >= 1:
        platform = dependencies.get('cinco-platform')

        envs.append({
            'name': 'CINCO_SERVER_URL',
            'value': 'https://' + platform.domains['primary'] + '/'
        })
        envs.append({
            'name': 'CINCO_CONSOLE_URL',
            'value': 'https://' + domains['primary'] + '/'
        })

        core.logger.info('Setting CINCO_CONSOLE_URL to ' + 'https://' + platform.domains['primary'] + '/')
        core.logger.info('Setting CINCO_SERVER_URL to ' + 'https://' + domains['primary'] + '/')