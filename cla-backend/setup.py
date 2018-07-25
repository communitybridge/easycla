"""Setup script for the CLA system."""

import re
import codecs
from os import path
from setuptools import setup, find_packages

def find_version(filename):
    """Helper function to find the currently CLA system version."""
    with open(filename) as fhandle:
        version_match = re.search(r"^__version__ = ['\"]([^'\"]*)['\"]",
                                  fhandle.read(), re.M)
        if version_match:
            return version_match.group(1)
        raise RuntimeError("Unable to find version string.")

setup(
    name='cla',
    version=find_version('cla/__init__.py'),
    description='REST endpoint to manage CLAs',
    long_description='See the CLA GitHub repository for more details: ' + \
                     'https://github.com/linuxfoundation/cla',
    url='https://github.com/linuxfoundation/cla',
    author='***REMOVED*** ***REMOVED***',
    author_email='***REMOVED***@linuxfoundation.org',
    #license='BSD',
    classifiers=[
        'Development Status :: 3 - Alpha',
        'Environment :: Web Environment',
        'Framework :: Hug',
        'Intended Audience :: Developers',
        #'License :: OSI Approved :: BSD License',
        'Natural Language :: English',
        'Programming Language :: Python :: 3.6',
    ],
    keywords='cla',
    packages=find_packages(),
    package_data={'cla': ['resources/*.html']},
    install_requires=['boto3>=1.4.4,<2.0',
                      'gossip>=2.2.0,<3.0',
                      'gunicorn>=19.7.1,<20.0',
                      'hug>=2.2.0,<3.0',
                      'pydocusign>=1.2.0,<2.0',
                      'pygithub>=1.34.0,<2.0',
                      'pynamodb>=2.1.6,<3.0',
                      'python-gitlab>=0.21.2,<1.0',
                      'requests-oauthlib>=0.8.0,<1.0',
                      'python-jose>=1.4.0'],
    extras_require={'dev': ['pylint'], 'test': ['nose2', 'coverage']},
    dependency_links=['https://nexus.engineering.tux.rocks/repository/pypi-engineering/packages/lf-keycloak/2.0.4/lf-keycloak-2.0.4.tar.gz'],
    entry_points={'console_scripts': []},
)
