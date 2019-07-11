// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

'use strict';
const AWS = require('aws-sdk');
const invalidateDistributions = require('./src/invalidate-distributions');
const getDomain = require('./src/get-domain');

class InvalidateCloudfront {
  constructor(serverless, options) {
    this.serverless = serverless;
    this.options = options || {};

    this.commands = {
      invalidate: {
        usage: 'Invalidate the Cloudfront cache',
        lifecycleEvents: ['invalidate']
      },
      'get-domain': {
        usage: 'Get a cloudfront domain',
        lifecycleEvents: ['get-domain'],
        options: {
          distribution: {
            usage: 'Specify the distribution you want to the domain of (e.g. "--distribution MyDistribution")',
            shortcut: 'd',
            required: true
          }
        }
      }
    };

    this.hooks = {
      'invalidate:invalidate': () => this.invalidateCloudfrontDistributions(),
      'get-domain:get-domain': () => this.getDomain()
    };
  }

  invalidateCloudfrontDistributions() {
    const distributions = this.serverless.service.custom.invalidateCloudfront;
    const resources = this.serverless.service.resources.Resources;
    const cli = this.serverless.cli;
    const aws = this.serverless.getProvider('aws');

    if (distributions === undefined || resources === undefined) {
      return;
    }

    const credentials = aws.getCredentials().credentials;
    const cloudfront = new AWS.CloudFront({
      credentials
    });

    return invalidateDistributions(aws, distributions, cloudfront, cli);
  }

  getDomain() {
    const aws = this.serverless.getProvider('aws');

    const credentials = aws.getCredentials().credentials;
    const cloudfront = new AWS.CloudFront({
      credentials
    });
    const dist = this.options.distribution;

    return getDomain(dist, aws, cloudfront).then(domain => {
      console.log(domain);
    });
  }
}

module.exports = InvalidateCloudfront;
