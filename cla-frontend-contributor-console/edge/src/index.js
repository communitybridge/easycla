// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const addHeaders = require('./add-headers');
const setCacheControl = require('./set-cache-control');

exports.handler = (event, context, callback) => {
  const headers = HEADERS;
  const resourcesNotToCache = ['/index.html', '/'];
  const resource = event.Records[0].cf.request.uri;
  const timeToLive = 60 * 60 * 24 * 365;
  const modifiedHeaders = setCacheControl.setCacheControl(headers, resource, resourcesNotToCache, timeToLive);
  const response = addHeaders.addHeaders(event, modifiedHeaders);
  callback(null, response);
};
