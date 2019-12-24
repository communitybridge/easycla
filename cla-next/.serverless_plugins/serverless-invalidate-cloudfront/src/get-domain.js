function getNameFromDistributionInfo(distributionId, cloudfront) {
  const params = {
    Id: distributionId
  };

  return cloudfront
    .getDistribution(params, (data, err) => {})
    .promise()
    .then(data => {
      if (!data.DomainName) {
        throw Error('GetDomain: No domain name found');
      }
      return data.DomainName;
    });
}

/**
 * Retrieves the domain name from a distribution.
 * @param {string} distributionName The resource name of the distribution
 * @param {AWS} aws The aws sdk
 * @param {Cloudfront} cloudfront An AWS cloudfront service.
 * @returns {Promise<string>} A promise containing the domain name of the resource.
 */
function getDomain(distributionName, aws, cloudfront) {
  const stackName = aws.naming.getStackName();

  return aws.request('CloudFormation', 'describeStackResources', { StackName: stackName }).then(resp => {
    const stackResources = resp.StackResources;
    const resource = stackResources.find(r => r.LogicalResourceId === distributionName);
    if (resource === undefined) {
      throw Error(`Unable to find distribution matching '${name}'.`);
    }
    return getNameFromDistributionInfo(resource.PhysicalResourceId, cloudfront);
  });
}

module.exports = getDomain;
