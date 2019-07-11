// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

const getDomain = require('./get-domain');

describe('getDomain', () => {
  let aws, cloudfront, promise;
  beforeEach(() => {
    const stacks = {
      StackResources: [
        {
          ResourceType: 'AWS::CloudFront::Distribution',
          LogicalResourceId: 'dist1',
          PhysicalResourceId: 'physical1'
        }
      ]
    };
    const stackName = 'stack';
    aws = {
      naming: jasmine.createSpyObj('naming', {
        getStackName: jasmine.createSpy('getStackName').and.returnValue(stackName)
      }),
      request: jasmine.createSpy('request').and.returnValue(Promise.resolve(stacks))
    };

    cloudfront = {
      getDistribution: jasmine.createSpy('getDistribution').and.returnValue({
        promise: () =>
          Promise.resolve({
            DomainName: 'www.somedomain.com'
          })
      })
    };
  });

  it('gets a domain name from a distribution name', done => {
    getDomain('dist1', aws, cloudfront).then(name => {
      expect(name).toBe('www.somedomain.com');
      done();
    });
  });
});
