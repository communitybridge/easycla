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
        del containers['workspace']
        notice('No workspace container provisioned in dev mode')
        notice('Ensure the cla_config.py file has the appropriate configuration for your external dependencies')
        base_url = 'http://' + lfcore.utils.public_ip()
        callback_url = base_url + '/v1/signed'
        github_oauth_url = base_url + '/v1/github/installation'
        mailhog_port = lfcore.utils.host_port(containers, 'mailhog')
        mailhog_web_port = lfcore.utils.host_port(containers, 'mailhog', 1)
        dynamodb_port = lfcore.utils.host_port(containers, 'dynamodb')
        notice('Example: BASE_URL = %s' %base_url)
        notice('Example: SIGNED_CALLBACK_URL = %s' %callback_url)
        notice('Example: GITHUB_OAUTH_CALLBACK_URL = %s' %github_oauth_url)
        notice('Example: DATABASE_HOST = http://localhost:%s' %dynamodb_port)
        notice('Example: KEYVALUE_HOST = http://localhost:%s' %dynamodb_port)
        notice('Example: EMAIL_SERVICE = SMTP')
        notice('Example: SMTP_HOST = localhost')
        notice('Example: SMTP_PORT = %s' %mailhog_port)
        notice('Example: KEYCLOAK_ENDPOINT = http://localhost:<port>/auth/')
        notice('Example: KEYCLOAK_REALM = LinuxFoundation')
        notice('Example: KEYCLOAK_CLIENT_ID = cla')
        notice('Example: KEYCLOAK_CLIENT_SECRET = secret')
        notice('Example: CINCO_ENDPOINT = http://localhost:<port>')
        notice('Example: CLA_CONSOLE_ENDPOINT = http://localhost:<port>/#')
        notice('Example: DOCRAPTOR_API_KEY = <key>')
        notice('Example: DOCUSIGN_INTEGRATOR_KEY = 828f4192-0f3f-4f9f-9cff-56f0b6d61615')
        notice('Example: DOCUSIGN_USERNAME = <username>')
        notice('Example: DOCUSIGN_PASSWORD = <password>')
        notice('If you are behind a NAT/firewall, you will need to add port forwarding from the edge of your network to you local machine')
        notice('The CLA system will work without the port forwading setup, but you will not be able to test GitHub integration (starting flow from GitHub) and DocuSign callbacks (confirmation of completed signatures)')
        success('MailHog: http://localhost:%s' %mailhog_web_port)
        notice('You will need to initialize the DB manually by running the create_database.py script in the helpers/ folder')
    else:
        info('Automatically setting CLA dependency endpoints')
        platform = dependencies.get('cinco')

        if platform is not None and platform.created:
            cinco_endpoint = platform.endpoints.containers.get('workspace', 5000).formatted
            success('CINCO: %s' %cinco_endpoint)
            containers['workspace']['environment']['CLA_CINCO_ENDPOINT'] = cinco_endpoint
            keycloak = platform.instance.dependencies.get('keycloak')

            if keycloak is not None and keycloak.created:
                kc_endpoint = keycloak.endpoints.containers.get('workspace', 8080).formatted
                success('Keycloak: %s' %kc_endpoint)
                containers['workspace']['environment']['CLA_KEYCLOAK_ENDPOINT'] = kc_endpoint + '/auth/'

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
        success('MailHog: http://localhost:%s' %mailhog_port)
        success('Web: http://localhost:%s' %web_port)
