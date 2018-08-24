const { dev, prod} = require('@ionic/app-scripts/config/webpack.config');
const webpack = require('webpack');

const ENV = process.env.IONIC_ENV;
const envConfigFile = require(`./config-${ENV}.json`);
const claURLConfig = envConfigFile.CLA_API_URL;

// const customConfig = {
//   plugins: [
//     new webpack.DefinePlugin({
//         webpackGlobalVars: {
//           claApiUrl: JSON.stringify(claURLConfig)
//         }
//     })
//   ]
// }

const claPlugin = new webpack.DefinePlugin({
        webpackGlobalVars: {
          claApiUrl: JSON.stringify(claURLConfig)
        }
      })

dev.plugins.push(claPlugin);
prod.plugins.push(claPlugin);

// devConfig.plugins.push(claPlugin);
// prodConfig.plugins.push(claPlugin);

// module.exports = {
//   dev: devConfig,
//   prod: prodConfig
// };
