// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const AWS = require('aws-sdk');
const SSM = require('aws-sdk/clients/ssm');
const program = require('commander');
const Promise = require('bluebird');

program
  .version('1.0.0')
  .option('-d, --debug', 'output extra debugging')
  .requiredOption('--from <region>', 'specifies the AWS region to pull the SSM parameters from')
  .requiredOption('--to <region>', 'specifies the AWS region to push the SSM parameters to')
  .requiredOption('-s, --stage <stage>', 'the environment stage')
  .parse(process.argv);

if (program.debug) console.log(program.opts());

console.log(`From region: ${program.from}`);
console.log(`To region:   ${program.to}`);
console.log(`Migrating SSM Parameters in region ${program.from} to region ${program.to} for stage ${program.stage}...`);

const parameters = [
  `cla-gh-app-private-key-${program.stage}`,
  `cla-gh-app-webhook-secret-${program.stage}`,
  `cla-gh-app-id-${program.stage}`,
  `cla-gh-oauth-client-id-${program.stage}`,
  `cla-gh-oauth-secret-${program.stage}`,
  `cla-auth0-domain-${program.stage}`,
  `cla-auth0-clientId-${program.stage}`,
  `cla-auth0-username-claim-${program.stage}`,
  `cla-auth0-algorithm-${program.stage}`,
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
];

const ssmFrom = new SSM({'region': program.from});
const ssmTo = new SSM({'region': program.to});

// If we do all these at the same time, we end up with the following error:
//   Error adding parameter undefined. Error code: ThrottlingException, Message: Rate exceeded
// So, we use this promise library which allows us to throttle the requests in batches of N...
const concurrencyFactor = 10;
Promise.map(parameters, function(param) {
  console.log(`Querying SSM Parameter: ${param} in region ${program.from} for stage ${program.stage}...`);
  const query = {
    "Name": param,
    "WithDecryption": false,
  };

  ssmFrom.getParameter(query, (err, data) => {
    if (err == null) {
      console.log("Data:");
      console.log(data);
      console.log(`Uploading parameter ${data.Parameter.Name} = ${data.Parameter.Value}, Description: ${data.Parameter.Description}`);
      const query = {
        "Name": data.Parameter.Name,
        "Value": data.Parameter.Value,
        "Description": (data.Parameter.Description !== undefined ? data.Parameter.Description : data.Parameter.Name),
        "Type": "String",
        "Overwrite": true,
        /* Can't use overwrite and Tags at the same time
        "Tags": [
          { "Key": "Name", "Value": "EasyCLA" },
          { "Key": "ServiceType", "Value": "Product" },
          { "Key": "Service", "Value": "EasyCLA" },
          { "Key": "ServiceRole", "Value": "Backend" },
          { "Key": "ProgrammingPlatform", "Value": "Go" },
          { "Key": "Owner", "Value": "David Deal" },
        ],
         */
      };
      console.log('Query:');
      console.log(query);

      let param = ssmTo.putParameter(query, (err, data) => {
        if (err == null) {
          //console.log('raw data = %o', data);
          console.log(`Set ${data.Parameter.Name} = ${data.Parameter.Value}`)
        } else {
          console.log(`Error adding parameter ${program.parameter}. Error code: ${err.code}, Message: ${err.message}`);
          console.log('error = %o', err);
        }
      });
    } else {
      console.log(`Error fetching parameter: ${param}. Error code: ${err.code}, Message: ${err.message}`);
      if (program.debug) {
        console.log('error = %o', err);
      }
    }
  });
}, {concurrency: concurrencyFactor}).then(function() {
  // All items have been processed in sets of concurrencyFactor (10) at most
  console.log("done");
});
