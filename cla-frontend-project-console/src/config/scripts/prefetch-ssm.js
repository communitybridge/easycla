// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const fs = require('fs');
const RetrieveSSMValues = require('./read-ssm');
const configVarArray = [
  'auth0-clientId',
  'auth0-domain',
  'cla-api-url',
  'gh-app-public-link',
  'cla-logo-s3-url',
  'proj-console-link'
];
const region = 'us-east-1';
const profile = process.env.AWS_PROFILE;
const stageEnv = process.env.STAGE_ENV;

async function prefetchSSM() {
  let result = {};
  console.log(`Start to fetch SSM values at ${stageEnv}...`);
  result = await RetrieveSSMValues(configVarArray, stageEnv, region, profile);

  // result['cla-api-url'] = 'https://api.dev.lfcla.com';
  fs.writeFile('./config/cla-env-config.json', JSON.stringify(result), function(err) {
    if (err) throw new Error(`Couldn't save SSM parameters to disk with error ${err}`);
    console.log('Fetching completed...');
  });
}

prefetchSSM();
