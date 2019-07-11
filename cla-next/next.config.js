// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

require('dotenv').config()

const path = require('path')
const Dotenv = require('dotenv-webpack')
const withSass = require('@zeit/next-sass')

module.exports = withSass(
  {
    webpack: config => {
      config.plugins = config.plugins || []
  
      config.plugins = [
        ...config.plugins,
  
        // Read the .env file
        new Dotenv({
          path: path.join(__dirname, '.env'),
          systemvars: true
        })
      ]
  
      return config
    },
    exportPathMap: function() {
      return {
        '/': { page: '/' }
      }
    }
  }
)