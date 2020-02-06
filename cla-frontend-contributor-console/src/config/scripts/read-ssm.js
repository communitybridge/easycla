// @ts-check

// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT
const AWS = require('aws-sdk');

/**
 * @param {string[]} variables
 * @param {string} stage
 * @param {string} region
 * @param {string} profile
 * @returns {Promise<{ [key:string]: string}>}
 */
async function retrieveSSMValues(variables, stage, region, profile) {
  const scopedVariables = variables.map((param) => {
    return `cla-${param}-${stage}`;
  });

  const result = await requestSSMParameters(scopedVariables, stage, region, profile);
  const parameters = result.Parameters;
  const error = result.$response.error;
  if (error !== null) {
    throw new Error(
      `Couldn't retrieve SSM parameters for stage ${stage} in region ${region} using profile ${profile} - error ${error}`
    );
  }
  const scopedParams = createParameterMap(parameters, stage);
  let params = {};
  Object.keys(scopedParams).forEach((key) => {
    const param = scopedParams[key];
    key = key.replace('cla-', '');
    key = key.replace(`-${stage}`, '');
    params[key] = param;
  });

  variables.forEach((variable) => {
    if (params[variable] === undefined) {
      throw new Error(
        `Missing SSM parameter with name ${variable} for stage ${stage} in region ${region} using profile ${profile}`,
      );
    }
  });
  return params;
}

/**
 * @param {string[]} variables
 * @param {string} stage
 * @param {string} region
 */
function requestSSMParameters(variables, stage, region, profile) {
  AWS.config.credentials = new AWS.SharedIniFileCredentials({profile});
  const ssm = new AWS.SSM({region: region});

  const ps = {
    Names: variables,
    WithDecryption: true
  };

  return ssm.getParameters(ps).promise();
}

/**
 * @param {AWS.SSM.Parameter[]} parameters
 * @param {string} stage
 */
function createParameterMap(parameters, stage) {
  return parameters.filter((param) => param.Name.endsWith(`-${stage}`))
    .map((param) => {
      const output = {};
      output[param.Name] = param.Value;
      return output;
    })
    .reduce((prev, current) => {
      return {...prev, ...current};
    }, {});
}

module.exports = retrieveSSMValues;
