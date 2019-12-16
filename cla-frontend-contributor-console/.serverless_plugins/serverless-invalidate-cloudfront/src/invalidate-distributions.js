'use strict';
const randomString = require('./random-string');
const matchDistributions = require('./match-distributions');

/**
 * A distribution configuration object.
 * @typedef {Object} Distribution
 * @property {String[]} paths - An array of paths to invalidate.
 */

/**
 * A distribution/physical id pair.
 * @typedef {Object} DistributionInfo
 * @property {Distribution} distribution
 * @property {String} distributionId
 * @property {String} name
 */

/**
 * Invalidates a cloudfront distribution.
 * @param {DistributionInfo} distributionInfo - Info about the distribution to invalidate.
 * @param {Cloudfront} cloudfront - An AWS cloudfront service.
 * @param {Cli} cli - An serverless cli service.
 */
function invalidateDistribution(distributionInfo, cloudfront, cli) {
  const reference = randomString(16);
  const paths = distributionInfo.distribution.paths;
  const distributionId = distributionInfo.distributionId;
  const name = distributionInfo.name;

  if (paths === undefined) {
    return Promise.reject('No paths defined');
  }

  const params = {
    DistributionId: distributionId,
    InvalidationBatch: {
      CallerReference: reference,
      Paths: {
        Quantity: paths.length,
        Items: paths
      }
    }
  };

  return cloudfront
    .createInvalidation(params, (error, data) => {
      if (!error) {
        cli.log(`InvalidateCloudfront: Invalidating Distribution '${name}(${distributionInfo.distributionId})'`);
      }
    })
    .promise();
}

/**
 * Retrieves information about a stack from AWS, and invalidates all matching
 * cloudfront distributions.
 * @param {Aws} aws - The AWS service.
 * @param {Array<Distribution>} distributions - A list of distribution objects, including invalidation paths.
 * @param {Cloudfront} cloudfront - An AWS cloudfront service.
 * @param {Cli} cli - An serverless cli service.
 */
function invalidateDistributions(aws, distributions, cloudfront, cli) {
  const stackName = aws.naming.getStackName();

  return aws
    .request('CloudFormation', 'describeStackResources', { StackName: stackName })
    .then((resp) => {
      return matchDistributions(distributions, resp.StackResources);
    })
    .then((distributions) => {
      const promises = distributions.map((pair) => invalidateDistribution(pair, cloudfront, cli));
      return Promise.all(promises);
    });
}

module.exports = invalidateDistributions;
