// @ts-check
const AWS = require('aws-sdk');

/**
 * @param {string[]} variables
 * @param {string} stage
 * @param {string} region
 * @param {string} profile
 * @returns {Promise<{ [key:string]: string}>}
 */
async function retrieveSSMValues(variables, stage, region, profile) {
  const result = await requestSSMParameters(variables, stage, region, profile);
  const parameters = result.Parameters;
  const error = result.$response.error;
  if (error !== null) {
    throw new Error(`Couldn't retrieve SSM paramters from AWS with error ${error}`);
  }
  const params = createParameterMap(parameters, stage);

  variables.forEach(variable => {
    if (params[variable] === undefined) {
      throw new Error(`Missing SSM parameter with name ${variable}-${stage}`);
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
  var credentials = new AWS.SharedIniFileCredentials({ profile });
  AWS.config.credentials = credentials;
  const ssm = new AWS.SSM({ region: region });

  const ps = {
    Names: variables.map(variable => `${variable}-${stage}`),
    WithDecryption: true
  };

  return ssm.getParameters(ps).promise();
}

/**
 * @param {AWS.SSM.Parameter[]} parameters
 * @param {string} stage
 */
function createParameterMap(parameters, stage) {
  const params = parameters
    .filter(param => param.Name.endsWith(`-${stage}`))
    .map(param => {
      const name = nameWithoutStage(param.Name, stage);
      const output = {};
      output[name] = param.Value;
      return output;
    })
    .reduce((prev, current) => {
      return { ...prev, ...current };
    }, {});
  return params;
}

/**
 * @param {string} name
 * @param {string} stage
 */
function nameWithoutStage(name, stage) {
  return name.slice(0, name.length - (stage.length + 1));
}

module.exports = retrieveSSMValues;
