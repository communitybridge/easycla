const fs = require('fs');
const RetrieveSSMValues = require('./read-ssm');
const configVarArray = ['auth0-clientId', 'auth0-domain', /*'cinco-api-url',*/ 'cla-api-url', /*'analytics-api-url',*/ 'gh-app-public-link', 'cla-logo-s3-url'];
const region = 'us-east-1';
const profile = process.env.AWS_PROFILE;
const stageEnv = process.env.STAGE_ENV;

async function prefetchSSM () {
  let result = {};
  console.log(`Start to fetch SSM values at ${stageEnv}...`);
  result = await RetrieveSSMValues(configVarArray, stageEnv, region, profile);

  // result['cla-api-url'] = 'http://localhost:5000';
  fs.writeFile ('./config/cla-env-config.json', JSON.stringify(result), function(err) {
    if (err) throw new Error(`Couldn't save SSM paramters to disk with error ${err}`);
    console.log('Fetching completed...');});
}

prefetchSSM();
