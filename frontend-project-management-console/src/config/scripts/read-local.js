const localConfig = require('../config-local.json');

/**
 * @param {string[]} variables
 * @returns {{ [key:string]: string }}
 */
async function retrieveLocalConfigValues(variables) {
  const parameterMap = {};
  variables.forEach( variable => {
    value = localConfig[variable];
    if (value === undefined) {
      throw new Error(`Couldn't retrieve value from local config for ${variable}`);
    }
    parameterMap[variable] = localConfig[variable];
  });
  return parameterMap;
}

module.exports = retrieveLocalConfigValues;
