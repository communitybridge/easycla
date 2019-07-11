// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

'use strict';

/**
 * A distribution configuration object.
 * @typedef {Object} Distribution
 * @property {String[]} paths - An array of paths to invalidate.
 */

/**
 * A stack resource object.
 * @typedef {Object} StackResource
 * @property {String} LogicalResourceId
 * @property {String} PhysicalResourceId
 * @property {String} ResourceType
 */

/**
 * A distribution/physical id pair.
 * @typedef {Object} DistributionInfo
 * @property {Distribution} distribution
 * @property {String} distributionId
 * @property {String} name
 */

/**
 *
 * @param {Object.<string, Distribution>} distributions - Array of distribution configuration options
 * @param {Array<StackResource>} stackResources
 * @param {String} stackName - The name of the stack.
 * @returns {DistributionInfo[]} - A distribution paired with it's physical id.
 */
function matchDistributions(distributions, stackResources, stackName) {
  const CLOUDFRONT_TYPE = 'AWS::CloudFront::Distribution';

  if (distributions === undefined) {
    return [];
  }

  return Object.keys(distributions)
    .map(distributionName => {
      const distribution = distributions[distributionName];

      const resource = stackResources.find(r => r.LogicalResourceId === distributionName);

      if (!resource) {
        throw new Error(
          `InvalidateCloudfront: Stack '${stackName}'did not have a resource with logical name '${distributionName}'`
        );
      }
      if (resource.ResourceType !== CLOUDFRONT_TYPE) {
        throw new Error(
          `InvalidateCloudfront: Stack '${stackName}' had resource with logical name '${distributionName}', but was of incorrect type '${
            resource.ResourceType
          }'`
        );
      }
      const distributionId = resource.PhysicalResourceId;
      return { distributionId, distribution, name: distributionName };
    })
    .filter(pair => pair !== undefined);
}

module.exports = matchDistributions;
