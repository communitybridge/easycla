// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const webpack = require('webpack');
const path = require('path');

const GetSecurityHeaders = require('./security-headers.js');

module.exports = (env) => {
  const securityHeaders = GetSecurityHeaders(env, false);

  return {
    target: 'node',
    context: path.join(__dirname, '/src'),
    entry: {
      index: ['./index.js']
    },
    output: {
      path: path.join(__dirname, './dist'),
      filename: '[name].js',
      libraryTarget: 'umd'
    },
    module: {
      rules: [
        {
          test: /\.js$/,
          loaders: ['babel-loader']
        }
      ]
    },
    plugins: [
      new webpack.DefinePlugin({
        HEADERS: JSON.stringify(securityHeaders)
      })
    ]
  };
};
