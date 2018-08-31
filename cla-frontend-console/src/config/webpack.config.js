const { dev, prod} = require('@ionic/app-scripts/config/webpack.config');
const webpack = require('webpack');
const RetrieveLocalConfigValues = require('./scripts/read-local');
const configVarArray = ['auth0-clientId', 'auth0-domain', 'cla-api-url'];
const stageEnv = process.env.STAGE_ENV;

module.exports = async env => {
  // Here we hard code stage name, it's not perfect since if a new stage created/modified, we also need to change it.
  const shouldReadFromSSM = (stageEnv !== undefined && stageEnv === 'staging') || (stageEnv !== undefined && stageEnv === 'prod') || (stageEnv !== undefined && stageEnv === 'qa');
  let configMap = {};

  // Here in the future, we maybe want to use Enum class to replace hard-code file name as indicator.
  if (shouldReadFromSSM){
    configMap = await RetrieveLocalConfigValues(configVarArray, `config-${stageEnv}.json`);
  } else {
    configMap = await RetrieveLocalConfigValues(configVarArray, 'config-local.json');
  }

  const claConfigPlugin = new webpack.DefinePlugin({
    webpackGlobalVars: {
      CLA_API_URL: JSON.stringify(configMap['cla-api-url']),
      AUTH0_DOMAIN: JSON.stringify(configMap['auth0-domain']),
      AUTH0_CLIENT_ID: JSON.stringify(configMap['auth0-clientId'])
    }
  });

  dev.plugins.push(claConfigPlugin);
  prod.plugins.push(claConfigPlugin);

  return {
    dev: dev,
    prod: prod
  };
};
