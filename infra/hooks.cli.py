import gossip
import lfcore

def info(msg): print(msg)
def notice(msg): print('\033[93m%s\033[0m' %msg)
def success(msg): print('\033[92m%s\033[0m' %msg)

@gossip.register('lf.instance.compose.generation', tags=['cla'])
def lf_instance_compose_generation(containers, config, dependencies, mode, path):
    """
    All necessary steps to integrate CINCO and Keycloak into the CLA system.

    If dev mode, we don't provision a workspace container and assume the user has manually
    configured the appropriate dependencies.
    """
    if mode == 'dev':
        notice('No workspace container provisioned in dev mode')
        notice('Ensure the cla_config.py file has the appropriate configuration for your \
                       external dependencies')
        del containers['workspace']
    else:
        info('Automatically setting CLA dependency endpoints')
        platform = dependencies.get('cinco')

        if platform.created:
            cinco_endpoint = platform.endpoints.containers.get('workspace', 5000).formatted
            success('CINCO: %s', cinco_endpoint)
            containers['workspace']['environment']['CLA_CINCO_ENDPOINT'] = cinco_endpoint
            keycloak = platform.instance.dependencies.get('keycloak')

            if keycloak.created:
                kc_endpoint = keycloak.endpoints.containers.get('workspace', 8080).formatted
                success('Keycloak: %s', kc_endpoint)
                containers['workspace']['environment']['CLA_KEYCLOAK_ENDPOINT'] = kc_endpoint

        # Set the callback_url to the CLA instance.
        ip = lfcore.utils.public_ip()
        web_port = lfcore.utils.host_port(containers, 'workspace')
        mailhog_port = lfcore.utils.host_port(containers, 'mailhog')
        base_url = 'http://' + ip
        callback_url = base_url + '/v1/signed'
        github_oauth_url = base_url + '/v1/github/installation'
        containers['workspace']['environment']['CLA_BASE_URL'] = base_url
        containers['workspace']['environment']['CLA_SIGNED_CALLBACK_URL'] = callback_url
        containers['workspace']['environment']['CLA_GITHUB_OAUTH_CALLBACK_URL'] = github_oauth_url

        notice('You public IP address (%s) was used as the base_url for the CLA system.' % ip)
        notice('If you are behind a NAT/firewall, you will need to add port forwarding from the edge of your network to you local machine (80 -> %s)' % web_port)
        notice('The CLA system will work without the port forwading setup, but you will not be able to test GitHub integration (starting flow from GitHub) and DocuSign callbacks (confirmation of completed signatures)')
        success('MailHog: http://localhost:%s' % mailhog_port)
        success('Web: http://localhost:%s' % web_port)
