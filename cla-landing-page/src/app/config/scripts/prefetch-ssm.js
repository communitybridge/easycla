// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const fs = require('fs');
const RetrieveSSMValues = require('./read-ssm');
const configVarArray = ['auth0-clientId', 'auth0-domain', 'proj-console-link', 'corp-console-link'];
const region = 'us-east-1';
const profile = process.env.AWS_PROFILE;
const stageEnv = process.env.STAGE_ENV;
const AWS_SSM_JSON_PATH = './src/app/config/cla-env-config.json';

async function prefetchSSM() {
  console.log(`Start to fetch SSM values at ${stageEnv}...`);
  const result = await RetrieveSSMValues(configVarArray, stageEnv, region, profile);
  console.log('Fetching completed.');

  //test for local
  // result['cla-api-url'] = 'http://localhost:5000';
  console.log(`Saving configuration to file: ${AWS_SSM_JSON_PATH}...`);
  fs.writeFile(AWS_SSM_JSON_PATH, JSON.stringify(result), function (err) {
    if (err) {
      throw new Error(`Couldn't save SSM parameters to disk with error ${err}`);
    }
    console.log('Save complete.');
  });
}

prefetchSSM();
