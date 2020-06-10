// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const AWS = require('aws-sdk');
const SSM = require('aws-sdk/clients/ssm');
const program = require('commander');

program
  .version('1.0.0')
  .option('-o, --overwrite', 'overwrite the value?')
  .option('-s, --securestring', 'create parameter as as secure string')
  .requiredOption('-r, --region <region>', 'specifies the AWS region')
  .requiredOption('-p, --parameter <name>', 'the parameter name')
  .requiredOption('-d, --desc <description>', 'the parameter description')
  .requiredOption('-v, --value <value>', 'the parameter value')
  .parse(process.argv);
if (program.debug) console.log(program.opts());
console.log(`program.description ${program.desc}`);
if (program.desc === undefined) {
  console.log('Missing --desc parameter');
  process.exit(1);
}

// Configure AWS
AWS.config.update({ region: program.region });

console.log(`Adding SSM Parameter: ${program.parameter} in region ${program.region}...`);
let query = {}
if (program.overwrite) {
  query = {
    "Name": program.parameter,
    "Value": program.value,
    "Description": program.desc,
    "Type": (program.securestring ? "SecureString" : "String"),
    "Overwrite": (program.overwrite ? true : false),
  };
} else {
  query = {
    "Name": program.parameter,
    "Value": program.value,
    "Description": program.desc,
    "Type": (program.securestring ? "SecureString" : "String"),
    "Overwrite": (program.overwrite ? true : false),
    "Tags": [
      { "Key": "Name", "Value": "EasyCLA" },
      { "Key": "ServiceType", "Value": "Product" },
      { "Key": "Service", "Value": "EasyCLA" },
      { "Key": "ServiceRole", "Value": "Backend" },
      { "Key": "ProgrammingPlatform", "Value": "Go" },
      { "Key": "Owner", "Value": "David Deal" },
    ],
  };
}

const ssm = new SSM();
let param = ssm.putParameter(query, (err, data) => {
  if (err == null) {
    console.log('raw data = %o', data);
    //console.log(`${data.Parameter.Name} = ${data.Parameter.Value}`)
  } else {
    console.log(`Error adding parameter ${program.parameter}. Error code: ${err.code}, Message: ${err.message}`);
    console.log('error = %o', err);
  }
});
