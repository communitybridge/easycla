// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

/**
 * Splits a header property in the format 'key=value', or just 'key', and returns an object in the format {key: value}.
 * @param {string} assignment
 */
function splitAssignmentPair(assignment) {
  const parts = assignment.split('=').map(value => value.trim());

  const obj = {};
  if (parts.length == 1) {
    obj[parts[0]] = true;
  } else {
    obj[parts[0]] = parts[1];
  }
  return obj;
}

/**
 * Splits a comma seperated header into a list of key-value pairs.
 * @param {string} headerValue
 */
function splitHeaderValue(headerValue) {
  return headerValue
    .split(',')
    .map(value => splitAssignmentPair(value))
    .reduce((previous, current) => Object.assign(previous, current), {});
}

/**
 * Modifies the Cache-Control header on a per resource basis.
 * @param {Object.<string, string>} headers - A list of preset headers.
 * @param {String} currentResourceName - The name of the current resource.
 * @param {Array<string>} resourcesNotToCache - A list of resources not to cache.
 * @param {Number} timeToLive - The time to cache objects for, if they are to be cached.
 */
exports.setCacheControl = function(headers, currentResourceName, resourcesNotToCache, timeToLive) {
  const existingCacheControl = headers['Cache-Control'] !== undefined ? headers['Cache-Control'] : '';
  const cacheValues = splitHeaderValue(existingCacheControl);
  const sMaxAge = cacheValues['s-maxage'];

  let newCacheControl = '';
  if (resourcesNotToCache.includes(currentResourceName)) {
    newCacheControl = 'no-cache, no-store, must-revalidate';
  } else {
    newCacheControl = `max-age=${timeToLive}`;
  }
  if (sMaxAge) {
    newCacheControl = `s-maxage=${sMaxAge}, ${newCacheControl}`;
  }

  return Object.assign({}, headers, { 'Cache-Control': newCacheControl });
};
