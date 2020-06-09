// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const AWS = require('aws-sdk');
const SSM = require('aws-sdk/clients/ssm');
const program = require('commander');

program
  .version('1.0.0')
  .option('-d, --debug', 'output extra debugging')
  .requiredOption('-r, --region <region>', 'specifies the AWS region')
  .requiredOption('-s, --stage <stage>', 'the environment stage')
  .parse(process.argv);
if (program.debug) console.log(program.opts());


// Configure AWS
AWS.config.update({ region: program.region });

console.log(`Querying SSM Parameters in region ${program.region} for stage ${program.stage}...`);

const parameters = [
  `cla-gh-app-private-key-${program.stage}`,
  `cla-gh-app-webhook-secret-${program.stage}`,
  `cla-gh-app-id-${program.stage}`,
  `cla-gh-oauth-client-id-${program.stage}`,
  `cla-gh-oauth-secret-${program.stage}`,
  `cla-gh-access-token-${program.stage}`,
  `cla-auth0-domain-${program.stage}`,
  `cla-auth0-clientId-${program.stage}`,
  `cla-auth0-username-claim-${program.stage}`,
  `cla-auth0-algorithm-${program.stage}`,
  `cla-auth0-platform-client-id-${program.stage}`,
  `cla-auth0-platform-client-secret-${program.stage}`,
  `cla-auth0-platform-audience-${program.stage}`,
  `cla-auth0-platform-url-${program.stage}`,
  `cla-auth0-platform-api-gw-${program.stage}`,
  `cla-sf-instance-url-${program.stage}`,
  `cla-sf-consumer-key-${program.stage}`,
  `cla-sf-consumer-secret-${program.stage}`,
  `cla-sf-username-${program.stage}`,
  `cla-sf-password-${program.stage}`,
  `cla-doc-raptor-api-key-${program.stage}`,
  `cla-docusign-root-url-${program.stage}`,
  `cla-docusign-username-${program.stage}`,
  `cla-docusign-password-${program.stage}`,
  `cla-docusign-integrator-key-${program.stage}`,
  `cla-api-base-${program.stage}`,
  `cla-contributor-base-${program.stage}`,
  `cla-contributor-v2-base-${program.stage}`,
  `cla-corporate-base-${program.stage}`,
  `cla-landing-page-${program.stage}`,
  `cla-signature-files-bucket-${program.stage}`,
  `cla-cla-logo-s3-url-${program.stage}`,
  `cla-ses-sender-email-address-${program.stage}`,
  `cla-lf-group-client-id-${program.stage}`,
  `cla-lf-group-client-secret-${program.stage}`,
  `cla-lf-group-refresh-token-${program.stage}`,
  `cla-lf-group-client-url-${program.stage}`,
  `cla-sns-event-topic-arn-${program.stage}`,
  `docraptor-test-mode-${program.stage}`,
  `cla-lfx-portal-url-${program.stage}`
];

const ssm = new SSM();
parameters.forEach((param) => {

  console.log(`Querying SSM Parameter: ${param} in region ${program.region} for stage ${program.stage}...`);
  const query = {
    "Name": param,
    "WithDecryption": false,
  };

  ssm.getParameter(query, (err, data) => {
    if (err == null) {
      console.log(`Parameter ${data.Parameter.Name} = ${data.Parameter.Value}`);
      if (program.debug) {
        console.log('Details: %o', data);
      }
    } else {
      console.log(`Error fetching parameter: ${param}. Error code: ${err.code}, Message: ${err.message}`);
      if (program.debug) {
        console.log('error = %o', err);
      }
    }
  });
});

